package installer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ComparePasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func NamespaceExists(cl client.Client, namespace string) (bool, error) {
	if err := cl.Get(
		context.TODO(),
		types.NamespacedName{Name: "argocd"},
		&corev1.Namespace{},
	); err == nil {
		return true, nil
	} else if errors.IsNotFound(err) {
		return false, nil
	} else {
		return false, err
	}
}

func WaitForDeployment(
	cl client.Client,
	name string,
	namespace string,
	retryInterval time.Duration,
	timeout time.Duration,
) error {
	for deadline := time.Now().Add(timeout); time.Now().Before(deadline); {
		// Wait for argocd-server deployment to good
		deployment := &appsv1.Deployment{}
		if err := cl.Get(
			context.TODO(),
			types.NamespacedName{Name: name, Namespace: namespace},
			deployment,
		); err == nil {
			var replicas int32 = 1
			if deployment.Spec.Replicas != nil {
				replicas = *deployment.Spec.Replicas
			}
			if deployment.Status.AvailableReplicas >= replicas {
				// All replicas are available
				return nil
			}
		} else if err != nil && !errors.IsNotFound(err) {
			return err
		}
		<-time.After(retryInterval)
	}
	return nil
}

var ErrDeploymentNotReady = fmt.Errorf("deployment is not ready")

func DeploymentIsHealthy(
	cl client.Client,
	name string,
	namespace string,
) error {
	// Wait for argocd-server deployment to good
	deployment := &appsv1.Deployment{}
	if err := cl.Get(
		context.TODO(),
		types.NamespacedName{Name: name, Namespace: namespace},
		deployment,
	); err != nil {
		return err
	}
	var replicas int32 = 1
	if deployment.Spec.Replicas != nil {
		replicas = *deployment.Spec.Replicas
	}
	if deployment.Status.AvailableReplicas >= replicas {
		// All replicas are available
		return nil
	}
	return ErrDeploymentNotReady
}

func ConfigureViper() {
	// Search in ~/.foldy and /foldy (for when inside docker)
	viper.AddConfigPath("/foldy")
	viper.AddConfigPath("$HOME/.foldy")
	viper.SetConfigType("yaml")

	// Load in config.yaml
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil && !strings.Contains(strings.TrimSpace(err.Error()), `Config File "config" Not Found in`) {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}

	// Load in credentials.yaml
	viper.SetConfigName("credentials")
	err = viper.MergeInConfig()
	if err != nil && !strings.Contains(strings.TrimSpace(err.Error()), `Config File "credentials" Not Found in`) {
		panic(fmt.Errorf("fatal error credentials file: %v", err))
	}

	// Allow environment override of username
	if username, ok := os.LookupEnv("FOLDY_USERNAME"); ok {
		viper.Set("username", username)
	}

	// Allow environment override of password
	if password, ok := os.LookupEnv("FOLDY_PASSWORD"); ok {
		viper.Set("password", password)
	}
}
