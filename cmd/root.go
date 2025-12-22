package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/upyun"
)

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
	// Loading global config file.
	var configFile string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", config.DefaultConfigFile(), "The config file")
	config.ReadConfig(configFile)
	// Init the UPYUN client.
	upyun.InitUpyunClient()
}
