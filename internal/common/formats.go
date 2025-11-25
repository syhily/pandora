package common

// Image formats
const (
	JPEG = "jpeg"
	JPG  = "jpg"
	PNG  = "png"
	AVIF = "avif"
	WEBP = "webp"
	GIF  = "gif"
	APNG = "apng"
	SVG  = "svg"
	BMP  = "bmp"
)

// SupportExtensions contains all supported image extensions
var SupportExtensions = map[string]struct{}{
	JPEG: {},
	JPG:  {},
	PNG:  {},
	AVIF: {},
	WEBP: {},
	GIF:  {},
	APNG: {},
	SVG:  {},
	BMP:  {},
}
