package image

import (
	"encoding/base64"
	"log"

	"github.com/h2non/bimg"
	"go.n16f.net/thumbhash"
)

const (
	BlurSize = 8
)

// ImageMetadata represents image metadata
type ImageMetadata struct {
	Slug     string `json:"slug"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Blurhash string `json:"blurhash"`
}

// ReadImageMetadata reads and generates metadata for an image
func ReadImageMetadata(file string, key string, content []byte) *ImageMetadata {
	if ok, _ := IsSupportedImage(file); ok {
		img := bimg.NewImage(content)
		img.Image()
		size, err := img.Size()
		if err != nil {
			log.Printf("Failed to read the image size for %v", file)
			return nil
		}
		coded, err := BimgToImage(img)
		if err != nil {
			log.Printf("Failed to decode the image for %v", file)
			return nil
		}
		var blurWidth = BlurSize
		var blurHeight = BlurSize
		if size.Width > size.Height {
			blurHeight = size.Height * blurWidth / size.Width
			if blurHeight == 0 {
				blurHeight = 1
			}
		} else if size.Width < size.Height {
			blurWidth = size.Width * blurHeight / size.Height
			if blurWidth == 0 {
				blurWidth = 1
			}
		}

		hash := thumbhash.EncodeImage(coded)
		return &ImageMetadata{
			Slug:     key,
			Width:    size.Width,
			Height:   size.Height,
			Blurhash: base64.StdEncoding.EncodeToString(hash),
		}
	}
	return nil
}
