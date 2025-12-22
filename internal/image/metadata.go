package image

import (
	"encoding/base64"
	"fmt"

	"github.com/h2non/bimg"
	"go.n16f.net/thumbhash"
)

const (
	BlurSize = 8
)

// ImageMetadata represents image metadata
type ImageMetadata struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Blurhash string `json:"blurhash"`
}

// GenMetadata generates metadata for an converted image
func GenMetadata(key string, content []byte) (*ImageMetadata, error) {
	if ok, _ := IsSupportedImage(key); ok {
		img := bimg.NewImage(content)
		img.Image()
		size, err := img.Size()
		if err != nil {
			return nil, fmt.Errorf("Failed to read the image size for %v. Cause: %w", key, err)
		}
		coded, err := BimgToImage(img)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode the image for %v. Cause: %w", key, err)
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
			Width:    size.Width,
			Height:   size.Height,
			Blurhash: base64.StdEncoding.EncodeToString(hash),
		}, nil
	}
	return nil, fmt.Errorf("Unsupported image format: %s", key)
}
