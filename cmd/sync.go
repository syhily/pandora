package cmd

import (
	"bufio"
	"cmp"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/image"
	"github.com/syhily/pandora/internal/s3"
)

const (
	ImageMetadataFile = "images.yml"
)

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync files to S3. Generate image metadata file with thumbhash.",
		Run: func(cmd *cobra.Command, args []string) {
			// Create S3 client.
			cfg := config.ReadConfig(configPath)
			client := s3.NewClient(&cfg.ImageS3)

			// Upload the files into the S3.
			var metas []image.ImageMetadata
			syncer := image.NewSyncer(client, forceUpload)
			for _, directory := range []string{"images", "uploads", "fonts"} {
				r := syncer.SyncDirectory(cfg.ProjectRoot, filepath.Join(cfg.ProjectRoot, directory))
				if r != nil {
					metas = append(metas, r...)
				}
			}
			log.Println("Successfully sync the directories")

			// Save metadata file into blog project.
			filename := filepath.Join(cfg.BlogRoot, "src", "content", "metas", ImageMetadataFile)

			log.Println("Start to generate image metadata")
			slices.SortStableFunc(metas, func(a, b image.ImageMetadata) int {
				return cmp.Compare(a.Slug, b.Slug)
			})
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
			if err != nil {
				log.Fatalf("Failed to generate the metadata file: %v", filename)
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			defer writer.Flush()

			encoder := yaml.NewEncoder(writer)
			encoder.SetIndent(2)
			err = encoder.Encode(metas)
			if err != nil {
				log.Fatalf("Failed to save image metadata: %v", err)
			}

			log.Printf("The image metadata is saved into the [%v]\n", filename)
		},
	}

	forceUpload = false
)

func init() {
	syncCmd.Flags().BoolVarP(&forceUpload, "force", "", false, "Force upload the files to S3")
	rootCmd.AddCommand(syncCmd)
}
