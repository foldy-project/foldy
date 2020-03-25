package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:  "version",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
