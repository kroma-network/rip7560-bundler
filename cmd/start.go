package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stackup-wallet/stackup-bundler/internal/start"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts an instance",
	Run: func(cmd *cobra.Command, args []string) {
		start.Rip7560Mode()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
