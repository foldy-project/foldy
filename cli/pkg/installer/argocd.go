package installer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
)

func init() {
	AddComponent(&CustomComponent{
		Name: "argocd",
		CRDs: []string{
			"applications.argoproj.io",
			"appprojects.argoproj.io",
		},
		Install: func(s *Installer) error {
			return s.installArgoCD()
		},
		Uninstall: func(s *Installer) error {
			// Delete argocd namespace (CRD removal happens elsewhere)
			return s.deleteNamespace("argocd")
		},
	})
}

// WaitForArgoCD waits for all the Argo CD deployments to come online
func (s *Installer) WaitForArgoCD() error {
	if s.Verbose {
		defer NewWaitingMessage("Argo CD", s.StatusUpdateInterval).Stop()
	}
	waitForDeployments := []string{
		"argocd-application-controller",
		"argocd-dex-server",
		"argocd-redis",
		"argocd-repo-server",
		"argocd-server",
	}
	dones := make([]chan error, len(waitForDeployments), len(waitForDeployments))
	for i, deploymentName := range waitForDeployments {
		done := make(chan error, 1)
		dones[i] = done
		go func(deploymentName string, done chan<- error) {
			defer close(done)
			done <- WaitForDeployment(
				s.client,
				deploymentName,
				"argocd",
				5*time.Second,
				30*time.Second)
		}(deploymentName, done)
	}
	var multi error
	for _, done := range dones {
		if err := <-done; err != nil {
			multi = multierror.Append(multi, err)
		}
	}
	return multi
}

func (s *Installer) installArgoCD() error {
	if err := s.createNamespace("argocd"); err != nil {
		return err
	}

	apply := s.RestartArgoCD

	if !apply {
		// We're not *trying* to restart, but we may have to
		deployment := &appsv1.Deployment{}
		if err := s.client.Get(
			context.TODO(),
			types.NamespacedName{
				Name:      "argocd-server",
				Namespace: "argocd"},
			deployment,
		); err != nil || deployment.Status.AvailableReplicas == 0 {
			// It's not known that argocd-server is reachable
			// at this point, so let's go ahead and reapply
			// the official manifests + our patch.
			apply = true
		}
	}

	if apply {
		// Install the yaml
		if err := s.exec("kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml"); err != nil {
			return err
		}
	}

	// Use a newer argocd image that has Helm v3+
	newImage, ok := viper.Get("argocd.image").(string)
	if !ok {
		newImage = "argoproj/argocd:latest"
	}

	checkArgoCDImageTag := func(deploymentName string) error {
		deployment := &appsv1.Deployment{}
		if err := s.client.Get(
			context.TODO(),
			types.NamespacedName{
				Name:      deploymentName,
				Namespace: "argocd",
			},
			deployment,
		); err != nil {
			return err
		}
		// TODO: configure argocd image in config
		image := deployment.Spec.Template.Spec.Containers[0].Image
		if image == "argoproj/argocd:v1.4.2" {
			command := fmt.Sprintf(`kubectl patch deployment %s -n argocd --type=json -p='[{"op": "add", "path": "/spec/template/spec/containers/0/image", "value": "%s"}]'`, deploymentName, newImage)
			if err := s.exec(command); err != nil {
				return err
			}
			// Verify the result
			if err := s.client.Get(
				context.TODO(),
				types.NamespacedName{
					Name:      deploymentName,
					Namespace: "argocd",
				},
				deployment,
			); err != nil {
				return err
			}
			if deployment.Spec.Template.Spec.Containers[0].Image != newImage {
				return fmt.Errorf("deployments/%s image change did not take effect", deploymentName)
			}
		} else if s.Verbose {
			log.Printf("deployments/%s image patch already applied", deploymentName)
		}
		return nil
	}

	isRunningInsecurely := func() (bool, error) {
		deployment := &appsv1.Deployment{}
		if err := s.client.Get(
			context.TODO(),
			types.NamespacedName{
				Name:      "argocd-server",
				Namespace: "argocd",
			},
			deployment,
		); err != nil {
			return false, err
		}
		command := deployment.Spec.Template.Spec.Containers[0].Command
		for _, component := range command {
			if component == "--insecure" {
				return true, nil
			}
		}
		return false, nil
	}

	configOK := make(chan error, 1)
	go func() {
		configOK <- s.patchArgoCDConfigMap()
		close(configOK)
	}()

	appControllerOK := make(chan error, 1)
	go func() {
		defer close(appControllerOK)
		appControllerOK <- checkArgoCDImageTag("argocd-application-controller")
	}()

	serverOK := make(chan error, 1)
	go func() {
		defer close(serverOK)
		serverOK <- checkArgoCDImageTag("argocd-server")
	}()

	repoServerOK := make(chan error, 1)
	go func() {
		defer close(repoServerOK)
		repoServerOK <- checkArgoCDImageTag("argocd-repo-server")
	}()

	secretOK := make(chan error, 1)
	go func() {
		defer close(secretOK)
		secret := &corev1.Secret{}
		if err := s.client.Get(
			context.TODO(),
			types.NamespacedName{
				Name:      "argocd-secret",
				Namespace: "argocd",
			},
			secret,
		); err != nil {
			secretOK <- err
			return
		}
		adminPassword, ok := secret.Data["admin.password"]
		if ok {
			if ComparePasswordHash(s.Password, adminPassword) {
				if s.Verbose {
					log.Printf("Argo CD admin password already matches local config")
				}
				secretOK <- nil
				return
			}
			if s.Verbose {
				log.Printf("Argo CD admin password differs from config. Synchronizing...")
			}
		}
		hash := HashPassword(s.Password)
		if s.ShowSecrets {
			secretPatchCommand := fmt.Sprintf(`kubectl -n argocd patch secret argocd-secret -p '{"stringData":{"admin.password": "%s","admin.passwordMtime": "'%s'"}}'`,
				hash,
				"$(date +%FT%T%Z)")
			secretOK <- s.exec(secretPatchCommand)
			return
		}
		if err := os.Setenv("PASSWORD_HASH", hash); err != nil {
			secretOK <- err
			return
		}
		secretPatchCommand := fmt.Sprintf(`kubectl -n argocd patch secret argocd-secret -p '{"stringData":{"admin.password": "'${PASSWORD_HASH}'","admin.passwordMtime": "'%s'"}}'`,
			"$(date +%FT%T%Z)")
		execErr := s.exec(secretPatchCommand)
		if err := os.Unsetenv("PASSWORD_HASH"); err != nil {
			secretOK <- err
			return
		}
		secretOK <- execErr
	}()
	if insecure, err := isRunningInsecurely(); err != nil {
		return err
	} else if insecure {
		if s.Verbose {
			log.Printf("deployments/argocd-server already running with insecure flag")
		}
	} else {
		// Patch the deployment so it's running in insecure mode.
		// We'll be using traefik and cert-manager to handle TLS.
		if err := s.exec(`kubectl patch deployment argocd-server -n argocd --type=json -p='[{"op": "add", "path": "/spec/template/spec/containers/0/command", "value": ["argocd-server", "--staticassets", "/shared/app", "--insecure"]}]'`); err != nil {
			return err
		}
		if insecure, err := isRunningInsecurely(); err != nil {
			return err
		} else if !insecure {
			return fmt.Errorf("--insecure was not added to the command")
		}
	}
	if err := <-configOK; err != nil {
		return err
	}
	if err := <-appControllerOK; err != nil {
		return err
	}
	if err := <-appControllerOK; err != nil {
		return err
	}
	if err := <-repoServerOK; err != nil {
		return err
	}
	if err := <-secretOK; err != nil {
		return err
	}

	return s.WaitForArgoCD()
}

func (s *Installer) RunCommandInArgoCDServer(command string, args ...interface{}) error {
	if s.argoCDPodName == "" {
		if err := s.ArgoCDSession(false); err != nil {
			return err
		}
	}
	interpolated := fmt.Sprintf(command, args...)
	// TODO: catch errors that could be fixed by calling ArgoCDSession() and retrying
	return s.exec("kubectl exec -n argocd %s -- %s", s.argoCDPodName, interpolated)
}

// IsArgoCDHealthy checks if argocd-server is up and running
// without looping.
func (s *Installer) IsArgoCDHealthy() error {
	return DeploymentIsHealthy(s.client, "argocd-server", "argocd")
}

func (s *Installer) patchArgoCDConfigMap() error {
	return nil
	config := &corev1.ConfigMap{}
	if err := s.client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      "argocd-cm",
			Namespace: "argocd",
		},
		config,
	); err != nil {
		return err
	}
	customizations, _ := config.Data["resource.customizations"]
	if strings.Contains(customizations, "extensions/Ingress") {
		if s.Verbose {
			log.Printf("argocd-cm already has ingress health check patch")
		}
		return nil
	}
	yaml := `---
apiVersion: v1
data:
  resource.customizations: |
      extensions/Ingress:
        health.lua: |
          hs = {}
          hs.status = "Healthy"
          return hs
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
  name: argocd-cm
  namespace: argocd
`
	if strings.Contains(yaml, "\t") {
		panic("const string should use spaces instead of tabs")
	}
	return s.exec("cat <<EOF | kubectl apply -f -\n%s", yaml)
	/*

			customizations += `|
		extensions/Ingress:
		  health.lua: |
		    hs = {}
			hs.status = "Healthy"
		    return hs
		`
			// I wish I had the source to cite for this...
			if s.Verbose {
				// I don't think it's currently possible to jsonpatch
				// this multiline string, so just faux-print the command
				// to the console and let's sneakily do it in code
				log.Printf(`> kubectl patch configmap argocd-cm -n argocd --type=merge -p '{"data":{"resource.customizations":"%s"}}'`, customizations)
			}
			newConfig := config.DeepCopy()
			newConfig.Data["resource.customizations"] = customizations
			if err := s.client.Update(
				context.TODO(),
				newConfig,
			); err != nil {
				return err
			}
			return nil
			//return s.exec(`kubectl patch configmap argocd-cm -n argocd --type=merge -p '{"data":{"resource.customizations":"%s"}}'`, customiziations)
	*/
}
