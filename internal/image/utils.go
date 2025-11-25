package image

import (
	"strings"
)

// IsSupportedImage checks if a file name has a supported image extension
func IsSupportedImage(name string) (bool, string) {
	ext := strings.ToLower(name[strings.LastIndex(name, ".")+1:])
	_, ok := SupportExtensions[ext]
	return ok, ext
}
