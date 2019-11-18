package cmd

import (
	"github.com/spf13/cobra"
	"runtime"
)

var (
	appVersion string
	buildDate  string
	gitCommit  string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf(`
Version: %s
GO Version: %s
Commit: %s
BuildTime: %s
`, appVersion, runtime.Version(), gitCommit, buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func SetVersion(version, date, commitId string)  {
	appVersion = version
	buildDate = date
	gitCommit = commitId
}