package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/syhily/pandora/internal/config"
)

var configPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pandora",
	Short: "A set of useful tools for writing in weblog",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", config.DefaultConfigRoot(), "The config file directory")
}
