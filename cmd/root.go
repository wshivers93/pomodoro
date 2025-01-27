package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pom",
	Short: "Simple timer for the pomodoro time management method",
}

func Execute() error {
	return rootCmd.Execute()
}

