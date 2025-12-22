package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/h2non/bimg"

	"github.com/syhily/pandora/internal/config"
)

// BimgToImage converts a bimg.Image into a standard Go image.Image.
// It supports JPEG, JPG, PNG, AVIF, WEBP, GIF, APNG, SVG, and BMP.
func BimgToImage(b *bimg.Image) (image.Image, error) {
	buf := b.Image()

	meta, err := b.Metadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	format := meta.Type
	reader := bytes.NewReader(buf)

	switch format {
	case config.JPEG, config.JPG:
		return jpeg.Decode(reader)
	case config.PNG, config.APNG:
		return png.Decode(reader)
	case config.GIF:
		return gif.Decode(reader)
	case config.BMP:
		// Convert BMP to PNG first (since Go's std doesn't support BMP decoding)
		pngBuf, err := bimg.NewImage(buf).Convert(bimg.PNG)
		if err != nil {
			return nil, err
		}
		return png.Decode(bytes.NewReader(pngBuf))
	case config.WEBP, config.AVIF, config.SVG:
		// Convert unsupported formats (WebP, AVIF, SVG) to PNG first
		pngBuf, err := bimg.NewImage(buf).Convert(bimg.PNG)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to PNG: %w", format, err)
		}
		return png.Decode(bytes.NewReader(pngBuf))
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ImageType converts a format string to bimg.ImageType
func ImageType(format string) bimg.ImageType {
	switch format {
	case config.JPG:
	case config.JPEG:
		return bimg.JPEG
	case config.PNG:
		return bimg.PNG
	case config.AVIF:
		return bimg.AVIF
	case config.GIF:
		return bimg.GIF
	case config.APNG:
		return bimg.PNG
	case config.BMP:
		return bimg.JPEG
	case config.WEBP:
		return bimg.WEBP
	case config.SVG:
		return bimg.SVG
	}
	return bimg.JPEG
}

// SupportedFormats returns a comma-separated string of supported formats
func SupportedFormats() string {
	extensions := make([]string, 0, len(config.SupportExtensions))
	for k := range config.SupportExtensions {
		extensions = append(extensions, k)
	}
	return strings.Join(extensions, ", ")
}

// IsSupportedImage checks if a file name has a supported image extension
func IsSupportedImage(name string) (bool, string) {
	ext := strings.ToLower(name[strings.LastIndex(name, ".")+1:])
	_, ok := config.SupportExtensions[ext]
	return ok, ext
}
