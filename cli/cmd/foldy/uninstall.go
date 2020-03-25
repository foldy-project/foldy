package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/foldy-project/foldy/cli/pkg/installer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var force bool

func init() {
	uninstallCmd.PersistentFlags().BoolVar(&skipDependencies, "skip-dependencies", false, "only install the specified components without installing dependencies")
	viper.BindPFlag("skipDependencies", uninstallCmd.PersistentFlags().Lookup("skip-dependencies"))

	uninstallCmd.PersistentFlags().BoolVar(&force, "force", false, "force uninstall without waiting for Argo CD")

	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:  "uninstall",
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Uninstalling...")
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
		install.Force = force
		if force {
			log.Printf("--force was specified. Uninstallation will not use Argo CD")
		}

		if len(args) == 0 {
			log.Printf("Uninstalling everything...")
			if err := install.UninstallAll(); err != nil {
				return err
			}
			log.Printf("all components were uninstalled successfully")
		} else {
			log.Printf("Uninstalling %v", args)
			if err := install.UninstallComponentsByName(args); err != nil {
				return err
			}
			if len(args) == 1 {
				log.Printf("component '%s' uninstalled", args[0])
			} else {
				log.Printf("components %v uninstalled", args)
			}
		}
		return nil
	},
}
