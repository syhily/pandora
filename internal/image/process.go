package image

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"golang.design/x/clipboard"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/upyun"
)

var imageLocalDatePattern = regexp.MustCompile(`^\d{8}$`)

// Process processes an image file according to the given parameters
func Process(file *os.File, width, height int, dt time.Time, format string, quality int) (string, []byte, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read the image %s: %w", file.Name(), err)
	}

	// Image conversion with optimization for web page.
	image := bimg.NewImage(bytes)
	it := ImageType(format)
	options := bimg.Options{
		Width:   width,
		Height:  height,
		Crop:    false,
		Quality: quality,
		Rotate:  0,
		Type:    it,
	}
	size, err := image.Size()
	if err != nil {
		return "", nil, fmt.Errorf("image is invalid: %w", err)
	}
	if height == 0 {
		options.Height = width * size.Height / size.Width
		options.Crop = false
	} else {
		options.Crop = true
	}
	bytes, err = image.Process(options)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert the images: %w", err)
	}

	// Generate remote image path.
	filename := dt.Format("20060102") + time.Now().Format("150405") + fmt.Sprintf("%02d", time.Now().Nanosecond()%100) + "." + format
	targetPath := strings.Join([]string{config.GetConfig().Asset.Path.Image, dt.Format("2006"), dt.Format("01"), filename}, "/")
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	// Upload to UPYUN.
	err = upyun.Upload(targetPath, bytes)
	if err != nil {
		return "", nil, fmt.Errorf("failed to upload the image to UPYUN. Cause: %w", err)
	}

	// Save into clipboard.
	link := fmt.Sprintf("%s://%s%s", config.GetConfig().Asset.Scheme, config.GetConfig().Asset.Domain, targetPath)
	clipboard.Write(clipboard.FmtText, []byte(link))

	return targetPath, bytes, nil
}

// ValidateDate validates the date format
func ValidateDate(date string) (time.Time, error) {
	if !imageLocalDatePattern.Match([]byte(date)) {
		return time.Time{}, fmt.Errorf("invalid local date format %s", date)
	}
	t, err := time.Parse("20060102", date)
	if err != nil {
		return time.Time{}, fmt.Errorf(`invalid time str %v. It should be "yyyyMMdd" like %v`, date, time.Now().Format("20060102"))
	}
	return t, nil
}
