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
	case "jpeg", "jpg":
		return jpeg.Decode(reader)
	case "png", "apng":
		return png.Decode(reader)
	case "gif":
		return gif.Decode(reader)
	case "bmp":
		// Convert BMP to PNG first (since Go's stdlib doesn't support BMP decoding)
		pngBuf, err := bimg.NewImage(buf).Convert(bimg.PNG)
		if err != nil {
			return nil, err
		}
		return png.Decode(bytes.NewReader(pngBuf))
	case "webp", "avif", "svg":
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
	case JPG:
	case JPEG:
		return bimg.JPEG
	case PNG:
		return bimg.PNG
	case AVIF:
		return bimg.AVIF
	case GIF:
		return bimg.GIF
	case APNG:
		return bimg.PNG
	case BMP:
		return bimg.JPEG
	case WEBP:
		return bimg.WEBP
	case SVG:
		return bimg.SVG
	}
	return bimg.JPEG
}

// SupportedFormats returns a comma-separated string of supported formats
func SupportedFormats() string {
	extensions := make([]string, 0, len(SupportExtensions))
	for k := range SupportExtensions {
		extensions = append(extensions, k)
	}
	return strings.Join(extensions, ", ")
}
