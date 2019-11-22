package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop service",
		Run: func(cmd *cobra.Command, args []string) {
			if App != nil {
				App.Stop()
			}
		},
	}

	rootCmd.AddCommand(stopCmd)
}
