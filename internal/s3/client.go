package s3

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/syhily/pandora/internal/config"
)

// Client encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the sync command.
// It contains client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type Client struct {
	bucket string
}

// NewClient creates a new S3 client from configuration
func NewClient(cfg *config.S3Configuration) *Client {
	return &Client{
		bucket: cfg.Bucket,
	}
}

// UploadObject reads from a file and puts the data into an object in a bucket.
func (c *Client) UploadObject(ctx context.Context, objectKey string, content []byte) error {
	return errors.New("not implemented")
}

// ListObjects lists the objects in a bucket.
func (c *Client) ListObjects(ctx context.Context, objectKey string) ([]types.Object, error) {
	return nil, errors.New("not implemented")
}
