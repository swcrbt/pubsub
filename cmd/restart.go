package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	/*var (
		configFile string
		isDaemon   bool
	)*/

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart service",
		Run: func(cmd *cobra.Command, args []string) {
			/*if err := App.Reload(); err != nil {
				cmd.Printf("reload services failed: %v\n", err)
			}*/
		},
	}

	//restartCmd.Flags().StringVarP(&configFile, "config", "c", "./config.toml", "config file path")
	//restartCmd.Flags().BoolVarP(&isDaemon, "daemon", "d", false, "running as a daemon")

	rootCmd.AddCommand(restartCmd)
}
