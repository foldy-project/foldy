package installer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Installer struct {
	client               client.Client
	Verbose              bool          // If true, print every command to stdout
	IgnoreAlreadyExists  bool          // If true, ignore AlreadyExists errors upon deleted resources
	IgnoreDeleteNotFound bool          // If true, ignore NotFound errors upon deleted resources
	SkipDependencies     bool          //
	Password             string        //
	RepoURL              string        //
	ShowSecrets          bool          // If false, inject secrets into commands as environment variables
	Force                bool          // If true, don't use argocd-server to manage resource deletion
	argocdL              sync.Mutex    //
	RestartArgoCD        bool          // If true, reapply Argo CD manifest, causing restart
	StatusUpdateInterval time.Duration // frequency to print periodic updates for asynchronous tasks
	argoCDPodName        string        //
}

func NewInstaller(cl client.Client) *Installer {
	s := &Installer{
		client:               cl,
		Verbose:              false,
		IgnoreAlreadyExists:  true,
		IgnoreDeleteNotFound: true,
		ShowSecrets:          false,
		Force:                false,
		Password:             "password",
		RepoURL:              "https://github.com/foldy-project/foldy.git",
		StatusUpdateInterval: 5 * time.Second,
	}
	s.ConfigureEnv()
	return s
}

func (s *Installer) ConfigureEnv() {
	s.SkipDependencies, _ = viper.Get("skipDependencies").(bool)
	s.Password, _ = viper.Get("password").(string)
	s.ShowSecrets, _ = viper.Get("showSecrets").(bool)
	s.Verbose, _ = viper.Get("verbose").(bool)
}

func (s *Installer) Reuse() {
	for _, comp := range components {
		comp.reuse()
	}
}

var ErrArgoCDNotInstalled = fmt.Errorf("Argo CD is not installed")

// ArgoCDSession retrieves and pins a port forward to Argo CD
func (s *Installer) ArgoCDSession(requireArgoCDExistImmediately bool) error {
	s.argocdL.Lock()
	defer s.argocdL.Unlock()
	if requireArgoCDExistImmediately {
		if exists, err := NamespaceExists(s.client, "argocd"); err != nil {
			return err
		} else if !exists {
			return ErrArgoCDNotInstalled
		}
	}
	if err := s.WaitForArgoCD(); err != nil {
		return err
	}
	// Login into the server
	pods := &corev1.PodList{}
	if err := s.client.List(
		context.TODO(),
		pods,
		client.InNamespace("argocd"),
	); err != nil {
		return err
	}
	var podName string
	for _, pod := range pods.Items {
		appName, ok := pod.ObjectMeta.Labels["app.kubernetes.io/name"]
		if !ok {
			continue
		} else if appName == "argocd-server" && pod.Status.Phase == "Running" {
			podName = pod.ObjectMeta.Name
			break
		}
	}
	if podName == "" {
		return fmt.Errorf("did not find ready argocd-server pod")
	} else if podName != s.argoCDPodName {
		// Pod changed
		if s.ShowSecrets {
			if err := s.exec("kubectl exec -n argocd %s -- argocd login localhost:8080 --plaintext --username admin --password %s", podName, s.Password); err != nil {
				return err
			}
		} else {
			if err := os.Setenv("ARGO_USERNAME", "admin"); err != nil {
				return nil
			}
			if err := os.Setenv("ARGO_PASSWORD", s.Password); err != nil {
				return nil
			}
			if err := s.exec("kubectl exec -n argocd %s -- argocd login localhost:8080 --plaintext --username $ARGO_USERNAME --password $ARGO_PASSWORD", podName); err != nil {
				return err
			}
			if err := os.Unsetenv("ARGO_USERNAME"); err != nil {
				return nil
			}
			if err := os.Unsetenv("ARGO_PASSWORD"); err != nil {
				return nil
			}
		}
		s.argoCDPodName = podName
	}
	return nil
}

func (s *Installer) exec(command string, args ...interface{}) error {
	var interpolated string
	if len(args) > 0 {
		interpolated = fmt.Sprintf(command, args...)
	} else {
		interpolated = command
	}
	cmd := exec.Command("bash", "-c", interpolated)
	if s.Verbose {
		log.Printf("> %s", interpolated)
	}
	r, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	resultStdout := make(chan interface{}, 1)
	resultStderr := make(chan interface{}, 1)
	go func() {
		defer close(resultStderr)
		stderr, err := ioutil.ReadAll(r)
		if err != nil {
			resultStderr <- err
			return
		}
		resultStderr <- string(stderr)
	}()
	go func() {
		defer close(resultStdout)
		stdout, err := ioutil.ReadAll(stdout)
		if err != nil {
			resultStdout <- err
			return
		}
		resultStdout <- string(stdout)
	}()
	if s.Verbose {
		defer NewWaitingMessage(interpolated, s.StatusUpdateInterval).Stop()
	}
	if err := cmd.Run(); err != nil {
		stdoutStr, _ := (<-resultStdout).(string)
		vStderr := <-resultStderr
		if stderr, ok := vStderr.(string); ok {
			return fmt.Errorf("%v: %s\n%s", err, stdoutStr, stderr)
		} else if err2, ok := vStderr.(error); ok {
			return fmt.Errorf("%v: %s\n%s", err, err2, stdoutStr)
		} else {
			panic("unreachable branch detected")
		}
	}
	return nil
}

func (s *Installer) installComponent(comp Component) error {
	comp.Lock()
	defer comp.Unlock()
	if comp.IsHandled() {
		return nil
	}
	if len(comp.GetDependencies()) > 0 && !s.SkipDependencies {
		if err := s.InstallComponentsByName(comp.GetDependencies()); err != nil {
			return err
		}
	}
	log.Printf("Installing %s", comp.GetName())
	if err := comp.RunInstall(s); err != nil {
		return err
	}
	log.Printf("Installed %s", comp.GetName())
	return nil
}

func (s *Installer) installComponents(components []Component) error {
	dones := make([]<-chan int, len(components), len(components))
	var errL sync.Mutex
	var multi error
	for i, comp := range components {
		done := make(chan int, 1)
		dones[i] = done
		go func(comp Component, done chan<- int) {
			defer close(done) // I much prefer this syntax
			err := s.installComponent(comp)
			if err != nil {
				// disambiguate the error just a little
				err = fmt.Errorf("%s: %v", comp.GetName(), err)
				errL.Lock()
				multi = multierror.Append(multi, err)
				errL.Unlock()
			}
			done <- 0
		}(comp, done)
	}
	for _, done := range dones {
		<-done
	}
	return multi
}

func (s *Installer) uninstallComponent(comp Component) error {
	comp.Lock()
	defer comp.Unlock()
	if comp.IsHandled() {
		return nil
	}
	if !s.SkipDependencies {
		dependees := GetDirectDependees(comp.GetName())
		if numDependees := len(dependees); numDependees > 0 {
			dones := make([]chan int, numDependees, numDependees)
			var multi error
			var multiL sync.Mutex
			for i, dependee := range dependees {
				done := make(chan int, 1)
				dones[i] = done
				go func(dependee string, done chan<- int) {
					defer close(done)
					if err := s.uninstallComponent(GetComponentByName(dependee)); err != nil {
						wrapped := fmt.Errorf("%v: %v", dependee, err)
						multiL.Lock()
						log.Printf("%v", wrapped)
						multi = multierror.Append(multi, wrapped)
						multiL.Unlock()
					}
					done <- 0
				}(dependee, done)
			}
			for _, done := range dones {
				<-done
			}
			if multi != nil {
				return multi
			}
		}
	}
	if err := comp.RunUninstall(s); err != nil {
		return err
	}
	if err := s.AsyncDelete("crd", comp.GetCRDs()); err != nil {
		return err
	}
	return nil
}

func (s *Installer) AsyncDelete(resource string, names []string) error {
	dones := make([]chan error, len(names), len(names))
	for i, name := range names {
		done := make(chan error, 1)
		dones[i] = done
		go func(name string, done chan<- error) {
			defer close(done)
			done <- func() error {
				if err := s.exec("kubectl delete %s %s", resource, name); err != nil {
					if !s.IgnoreDeleteNotFound {
						return err
					}
					if err := func() error {
						// Weird code structure so we get coverage
						msg := strings.TrimSpace(err.Error())
						if strings.Contains(msg, `Error from server (NotFound):`) {
							return nil
						}
						if strings.Contains(msg, `Error from server (Conflict):`) {
							return nil
						}
						return err
					}(); err != nil {
						return err
					}
				}
				return nil
			}()
		}(name, done)
	}
	var multi error
	for _, done := range dones {
		if err := <-done; err != nil {
			multi = multierror.Append(multi, err)
		}
	}
	return multi
}

func (s *Installer) uninstallComponents(components []Component) error {
	dones := make([]<-chan int, len(components), len(components))
	var errL sync.Mutex
	var multi error
	for i, comp := range components {
		done := make(chan int, 1)
		dones[i] = done
		go func(comp Component, done chan<- int) {
			defer func() {
				done <- 0
				close(done)
			}()
			if err := s.uninstallComponent(comp); err != nil {
				errL.Lock()
				log.Printf("Error uninstalling %s: %v", comp.GetName(), err)
				multi = multierror.Append(multi, err)
				errL.Unlock()
			}
		}(comp, done)
	}
	for _, done := range dones {
		<-done
	}
	return multi
}

func (s *Installer) InstallComponentsByName(names []string) error {
	components, err := GetComponentsByName(names)
	if err != nil {
		return err
	}
	return s.installComponents(components)
}

func (s *Installer) UninstallComponentsByName(names []string) error {
	components, err := GetComponentsByName(names)
	if err != nil {
		return err
	}
	return s.uninstallComponents(components)
}

func (s *Installer) InstallAll() error {
	return s.installComponents(components)
}

func (s *Installer) UninstallAll() error {
	return s.uninstallComponents(components)
}

func (s *Installer) createNamespace(namespace string) error {
	if err := s.exec("kubectl create namespace %s", namespace); err != nil {
		if !s.IgnoreAlreadyExists {
			return err
		}
		msg := fmt.Sprintf(`Error from server (AlreadyExists): namespaces "%s" already exists`, namespace)
		if !strings.Contains(strings.TrimSpace(err.Error()), msg) {
			return err
		}
	}
	return nil
}

func (s *Installer) deleteNamespace(namespace string) error {
	if err := s.exec("kubectl delete namespace %s", namespace); err != nil {
		if !s.IgnoreDeleteNotFound {
			return err
		}
		if err := func() error {
			// Weird code structure so we get coverage
			msg := strings.TrimSpace(err.Error())
			if strings.Contains(msg, fmt.Sprintf(`Error from server (NotFound): namespaces "%s" not found`, namespace)) {
				return nil
			}
			if strings.Contains(msg, fmt.Sprintf(`Error from server (Conflict): Operation cannot be fulfilled on namespaces "%s": The system is ensuring all content is removed from this namespace.  Upon completion, this namespace will automatically be purged by the system.`, namespace)) {
				return nil
			}
			return err
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Installer) CreateApplication(
	name string,
	repoURL string,
	path string,
) error {
	if err := s.createNamespace(name); err != nil {
		return err
	}
	if err := s.RunCommandInArgoCDServer("argocd app create %s --repo %s --path %s --dest-namespace %s --dest-server https://kubernetes.default.svc", name, repoURL, path, name); err != nil {
		return err
	}
	if err := s.RunCommandInArgoCDServer("argocd app sync %s", name); err != nil {
		return err
	}
	return nil
}

func (s *Installer) DeleteApplication(name string) error {
	exists, err := NamespaceExists(s.client, "argocd")
	if err != nil {
		return err
	}
	if !s.Force {
		// Do gentle uninstallation with Argo CD --cascade delete
		// This causes sub-applications to also be deleted
		if exists {
			if err := s.RunCommandInArgoCDServer("argocd app delete %s --cascade", name); err != nil {
				return err
			}
		} else if s.Verbose {
			log.Printf("Argo CD not installed. Proceeding with cleanup...")
		}
	}

	if err := s.deleteNamespace(name); err != nil {
		return err
	}

	if exists {
		deleted := make(chan error, 1)
		go func() {
			defer close(deleted)
			deleted <- func() error {
				if err := s.exec("kubectl delete application -n argocd %s", name); err != nil {
					msg := strings.TrimSpace(err.Error())
					if strings.Contains(msg, `error: the server doesn't have a resource type "application"`) {
						return nil
					}
					if strings.Contains(msg, `Error from server (NotFound):`) {
						return nil
					}
					return err
				}
				return nil
			}()
		}()
		delay := 10 * time.Second
		select {
		case err := <-deleted:
			if err != nil {
				return err
			}
		case <-time.After(delay):
			if s.Force {
				if s.Verbose {
					log.Printf("Removing finalizers for application argocd/%s", name)
				}
				if err := s.RemoveFinalizers("application", name, "argocd"); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("application argocd/%s was not deleted after %v. Maybe try --force?", name, delay)
			}
		}
	}

	return nil
}

func (s *Installer) RemoveFinalizers(
	resource string,
	name string,
	namespace string,
) error {
	return s.exec(`kubectl patch %s %s -n %s -p '{"metadata":{"finalizers": []}}' --type=merge`, resource, name, namespace)
}
