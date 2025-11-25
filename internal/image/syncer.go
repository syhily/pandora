package image

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/syhily/pandora/internal/s3"
)

// Syncer handles directory synchronization operations
type Syncer struct {
	client      *s3.Client
	forceUpload bool
}

// NewSyncer creates a new syncer
func NewSyncer(client *s3.Client, forceUpload bool) *Syncer {
	return &Syncer{
		client:      client,
		forceUpload: forceUpload,
	}
}

// SyncDirectory synchronizes a directory to S3 and returns image metadata
func (s *Syncer) SyncDirectory(root, path string) []ImageMetadata {
	var metas []ImageMetadata
	var wg sync.WaitGroup

	if stat, err := os.Stat(path); err != nil {
		log.Printf("Failed to read current directory %v", path)
		return metas
	} else if stat.IsDir() && !strings.HasPrefix(stat.Name(), ".") {
		// Load the files/directories from the current directory.
		files, e := os.ReadDir(path)
		if e != nil {
			log.Printf("Failed to read directory %v", path)
			return metas
		}

		// Load the path prefix from AWS S3.
		objs, e := s.client.ListObjects(context.TODO(), path[len(root)+1:])
		if e != nil {
			log.Printf("Failed to read directory from S3: %v\nError: %v", path[len(root):], e)
		}
		awsMetas := map[string]int64{}
		for _, obj := range objs {
			awsMetas[*obj.Key] = *obj.Size
		}

		// Range the files in the current directory.
		resultChan := make(chan []ImageMetadata, len(files))
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			} else if file.IsDir() {
				// Process directories concurrently.
				wg.Add(1)
				go func(subDir string) {
					defer wg.Done()
					syncer := NewSyncer(s.client, s.forceUpload)
					m := syncer.SyncDirectory(root, filepath.Join(path, subDir))
					if m != nil {
						resultChan <- m
					}
				}(file.Name())
			} else {
				// Process files concurrently.
				wg.Add(1)
				go func(filename string) {
					defer wg.Done()
					info, e1 := os.Stat(filename)
					if e1 != nil {
						log.Printf("Failed to read the file %v info", filename)
						return
					}
					key := strings.ReplaceAll(filename[len(root)+1:], string(filepath.Separator), "/")
					content, e2 := os.ReadFile(filename)
					if e2 != nil {
						log.Printf("Failed to read the file %v content", filename)
						return
					}
					if ok, _ := IsSupportedImage(filepath.Base(filename)); ok {
						meta := ReadImageMetadata(filename, filename[len(root):], content)
						if meta != nil {
							resultChan <- []ImageMetadata{*meta}
						}
					}
					if info.Size() != awsMetas[key] || s.forceUpload {
						log.Printf("Try to upload the file [%v] to the aws s3", filename)
						e2 = s.client.UploadObject(context.TODO(), key, content)
						if e2 != nil {
							log.Printf("Failed to upload the file %v to s3", filename)
							return
						}
					} else {
						log.Printf("Skip the existing file [%v] in aws s3", filename)
					}
				}(filepath.Join(path, file.Name()))
			}
		}

		// Wait for all goroutines to finish processing
		wg.Wait()
		close(resultChan)

		// Collect all metadata results from the channel
		for result := range resultChan {
			metas = append(metas, result...)
		}
	}

	return metas
}
