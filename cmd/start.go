package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if err := startCmd.MarkFlagRequired("mode"); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("mode", startCmd.Flags().Lookup("mode")); err != nil {
		panic(err)
	}
}
