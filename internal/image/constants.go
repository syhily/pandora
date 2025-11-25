package image

import "github.com/syhily/pandora/internal/common"

// Re-export constants for convenience
const (
	JPEG = common.JPEG
	JPG  = common.JPG
	PNG  = common.PNG
	AVIF = common.AVIF
	WEBP = common.WEBP
	GIF  = common.GIF
	APNG = common.APNG
	SVG  = common.SVG
	BMP  = common.BMP
)

// SupportExtensions re-exports the support extensions map
var SupportExtensions = common.SupportExtensions
