package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GatherGithubInfo",
	Short: "GatherGithubInfo is a tool to gather github users' info",
}

func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(initCmd)
}
