package installer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestInstall(t *testing.T) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	require.NoError(t, err)
	cl, err := client.New(config, client.Options{})
	require.NoError(t, err)
	ConfigureViper()
	install := NewInstaller(cl)
	t.Run("install", func(t *testing.T) {
		defer install.Reuse()
		assert.NoError(t, install.InstallAll())
		assert.NoError(t, install.CleanUp())
	})
	t.Run("uninstall", func(t *testing.T) {
		defer install.Reuse()
		assert.NoError(t, install.UninstallAll())
		assert.NoError(t, install.CleanUp())
	})
}
