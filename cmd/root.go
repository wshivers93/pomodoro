package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "Simple timer for the pomodoro time management method",
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

