package cmd

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/image"
	"github.com/syhily/pandora/internal/upyun"
)

var (
	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Convert images, resize, rename and upload to S3.",
		Run: func(cmd *cobra.Command, args []string) {
			// Check the image source path is valid.
			info, err := os.Stat(imageSource)
			if err != nil {
				log.Fatalf("Couldn't read the given file from the path %s, err: %v", imageSource, err)
			}

			if info.IsDir() {
				log.Fatalf("The given path %s is a directory. Only image is accepted", imageSource)
			}

			if ok, ext := image.IsSupportedImage(info.Name()); !ok {
				log.Fatalf("Unsupported file extension %s. Allowed extensions: %s", ext, image.SupportedFormats())
			}

			// Get the file operand
			file, err := os.Open(imageSource)
			if err != nil {
				log.Fatalf("Failed to read image %v", err)
			}
			defer file.Close()

			// File convert format check.
			if ok, _ := image.IsSupportedImage(imageFormat); !ok {
				log.Fatalf("Invalid convert format, only supports %s", image.SupportedFormats())
			}

			// Check the time pattern is valid.
			if imageLocalDate == "" {
				imageLocalDate = time.Now().Format("20060102")
			}
			dt, err := image.ValidateDate(imageLocalDate)
			if err != nil {
				log.Fatalf("Invalid date: %v", err)
			}

			if imageQuality == 0 {
				imageQuality = 85
			}
			if imageFormat == "" {
				imageFormat = "jpg"
			}

			key, bytes, err := image.Process(file, width, height, dt, imageFormat, imageQuality)
			if err != nil {
				log.Fatalf("Failed to process image: %v", err)
			}

			// Upload the image metadata to UPYUN.
			meta, err := image.GenMetadata(key, bytes)
			if err != nil {
				log.Fatalf("Failed to generate metadata: %v", err)
			}
			metaBytes, err := json.Marshal(meta)
			if err != nil {
				log.Fatalf("Failed to serialize metadata: %v", err)
			}
			upyun.Upload(key[:strings.LastIndex(key, ".")]+".json", metaBytes)

			log.Printf("Image processed successfully: %s (%d bytes)", key, len(bytes))
		},
	}

	width          = 1280
	height         = 0
	imageSource    = ""
	imageLocalDate = ""
	imageFormat    = ""
	imageQuality   = 0
	uploadImage    = true
)

func init() {
	imageCmd.Flags().StringVarP(&imageSource, "source", "s", "", "The image file path (absolute of relative)")
	imageCmd.Flags().IntVarP(&width, "width", "w", 1280, "The resized image width")
	imageCmd.Flags().IntVarP(&height, "height", "", 0, "The optional image height, 0 for keep ratio")
	imageCmd.Flags().StringVarP(&imageLocalDate, "time", "t", "", "The date time, in 20060102 format")
	imageCmd.Flags().StringVarP(&imageFormat, "format", "f", config.JPG, "The image format")
	imageCmd.Flags().IntVarP(&imageQuality, "quality", "q", 0, "The image quality")
	imageCmd.Flags().BoolVarP(&uploadImage, "upload", "u", true, "Whether to upload image")

	if err := imageCmd.MarkFlagRequired("source"); err != nil {
		log.Fatalf("%v", err)
	}

	rootCmd.AddCommand(imageCmd)
}
