package image

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"golang.design/x/clipboard"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/s3"
)

var imageLocalDatePattern = regexp.MustCompile(`^\d{8}$`)

// Processor handles image processing operations
type Processor struct {
	config   *config.PandoraConfig
	s3Client *s3.Client
}

// NewProcessor creates a new image processor
func NewProcessor(cfg *config.PandoraConfig) *Processor {
	return &Processor{
		config:   cfg,
		s3Client: s3.NewClient(&cfg.ImageS3),
	}
}

// Process processes an image file according to the given parameters
func (p *Processor) Process(file *os.File, width, height int, dt time.Time, format string, quality int, upload bool) error {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read the image %s: %w", file.Name(), err)
	}

	// Image conversion.
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
		return fmt.Errorf("image is invalid: %w", err)
	}
	if height == 0 {
		options.Height = width * size.Height / size.Width
		options.Crop = false
	} else {
		options.Crop = true
	}
	bytes, err = image.Process(options)
	if err != nil {
		return fmt.Errorf("failed to convert the images: %w", err)
	}

	// Create directory.
	directory := filepath.Join(p.config.ProjectRoot, "images", dt.Format("2006"), dt.Format("01"))
	err = os.MkdirAll(directory, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("failed to create the image directory: %w", err)
	}

	// Save image file.
	filename := dt.Format("20060102") + time.Now().Format("150405") + fmt.Sprintf("%02d", time.Now().Nanosecond()%100) + "." + format
	file, err = os.OpenFile(filepath.Join(directory, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to generate the target image file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush image: %w", err)
	}

	log.Printf("The image is saved into the [%v]\n", filepath.Join(directory, filename))

	if upload {
		// Upload S3
		err = p.s3Client.UploadObject(context.TODO(), strings.ReplaceAll(filepath.Join(directory, filename)[len(p.config.ProjectRoot)+1:], string(filepath.Separator), "/"), bytes)
		if err != nil {
			return fmt.Errorf("failed to upload the generated images to s3: %w", err)
		}

		link, _ := url.JoinPath(p.config.ImageS3.PublicDomain, "images", dt.Format("2006"), dt.Format("01"), filename)
		log.Printf("You can use link for document [%v]\n", link)
		// Save into clipboard
		clipboard.Write(clipboard.FmtText, []byte(link))
	}

	return nil
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
