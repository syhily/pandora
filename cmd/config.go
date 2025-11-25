package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/syhily/pandora/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactive configure pandora tool.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.PromptConfig()
		if err := config.WriteConfig(configPath, cfg); err != nil {
			log.Fatalf("Failed to write config: %v", err)
		}
		log.Printf("Configuration file created successfully at %s", configPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
