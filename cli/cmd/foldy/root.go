package main

import (
	"fmt"
	"log"
	"os"

	"github.com/foldy-project/foldy/cli/pkg/installer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "foldy",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("verbose=%b", viper.Get("verbose"))
		log.Printf("password=%v", viper.Get("password"))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode (print mutating shell commands)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("show-secrets", "x", false, "print secrets to stdout instead of injecting them as env variables")
	viper.BindPFlag("showSecrets", rootCmd.PersistentFlags().Lookup("show-secrets"))

	installer.ConfigureViper()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
