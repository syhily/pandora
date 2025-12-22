package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/music"
)

var (
	musicCmd = &cobra.Command{
		Use:   "music",
		Short: "Download and upload Netease music to S3 with metadata file.",
		Run: func(cmd *cobra.Command, args []string) {
			processor := music.NewProcessor(config.GetConfig())
			if err := processor.Process(musicId, vip); err != nil {
				log.Fatalf("Failed to process music: %v", err)
			}
		},
	}
	musicId int
	vip     bool
)

func init() {
	musicCmd.Flags().IntVarP(&musicId, "id", "i", 0, "The music id")
	musicCmd.Flags().BoolVarP(&vip, "vip", "v", false, "Download through VIP API")
	if err := musicCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("%v", err)
	}
	rootCmd.AddCommand(musicCmd)
}
