package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "ping-and-report",
	Short: "A tool for alive monitoring",
	Long:  "A tool for alive monitoring. Look at: https://github.com/blackbracken/ping-and-report",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to start this command using cobra")
	}
}
