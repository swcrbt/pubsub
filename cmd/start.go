package cmd

import (
	"github.com/spf13/cobra"
	"go-issued-service/app"
)

func init() {
	var configFile string

	var startCmd = &cobra.Command{
		Use: "start",
		Short: "Start service",
		Run: func(cmd *cobra.Command, args []string) {
			app.New(configFile).Run()
		},
	}

	startCmd.Flags().StringVarP(&configFile, "config", "c", "./config.toml", "config file path")

	rootCmd.AddCommand(startCmd)
}