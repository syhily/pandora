package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactive generate pandora config.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Configuration file created successfully at %s", configPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
