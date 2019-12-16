package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.orayer.com/golang/pubsub/app"
	"gitlab.orayer.com/golang/server"
)

func init() {
	var (
		configFile string
		isDaemon   bool
	)

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start service",
		Run: func(cmd *cobra.Command, args []string) {
			if (isDaemon) {
				server.Daemon()
			}

			App = app.New(configFile)
			App.Run()
		},
	}

	startCmd.Flags().StringVarP(&configFile, "config", "c", "./config.toml", "config file path")
	startCmd.Flags().BoolVarP(&isDaemon, "daemon", "d", false, "running as a daemon")

	rootCmd.AddCommand(startCmd)
}
