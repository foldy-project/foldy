package main

import (
	"path/filepath"

	"github.com/foldy-project/foldy/cli/pkg/portfwd"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(portfwdCmd)
}

var portfwdCmd = &cobra.Command{
	Use:   "portfwd",
	Short: "Forwards local ports to the cluster",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}
		p, err := portfwd.NewFoldyPortForwarder(config)
		if err != nil {
			return err
		}
		p.Verbose = true
		p.AddAllPorts()
		for {
		}
		return nil
	},
}
