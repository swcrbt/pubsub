package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.orayer.com/golang/issue/app"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gitlab.orayer.com/golang/issue",
	Short: "Long connection data Issue service",
}

var App *app.App

// Execute 执行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
