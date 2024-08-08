package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stackup-bundler",
	Short: "RIP-7560 Bundler",
	Long:  "A modular Go implementation of an RIP-7560 Bundler.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
