package main

import (
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/foldy-project/foldy/cli/pkg/installer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var skipDependencies bool

func init() {
	installCmd.PersistentFlags().BoolVar(&skipDependencies, "skip-dependencies", false, "only install the specified components without installing dependencies")
	viper.BindPFlag("skipDependencies", installCmd.PersistentFlags().Lookup("skip-dependencies"))

	installCmd.PersistentFlags().StringP("password", "p", "", "installation password")
	viper.BindPFlag("password", installCmd.PersistentFlags().Lookup("password"))

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [components...]",
	Short: "Ensures components of foldy are healthy",
	Long: `Ensures components of foldy are healthy, reinstalling them where necessary. This operation is idempotent and can be applied multiple times to converge the installation to a healthy state.

  # Install everything according to config.yaml
  foldy install

  # Install core application with all of its dependencies
  foldy install foldy

  # Install only specified components
  foldy install argocd cert-manager argo-events`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing...")

		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}
		cl, err := client.New(config, client.Options{})
		if err != nil {
			return err
		}
		install := installer.NewInstaller(cl)
		if len(args) == 0 {
			log.Printf("Installing everything...")
			if err := install.InstallAll(); err != nil {
				return err
			}
			log.Printf("all components appear to be healthy")
		} else {
			log.Printf("Installing %v", args)
			if err := install.InstallComponentsByName(args); err != nil {
				return err
			}
			if len(args) == 1 {
				log.Printf("component '%s' appears to be healthy", args[0])
			} else {
				log.Printf("components %v appear to be healthy", args)
			}
		}

		return nil
	},
}
