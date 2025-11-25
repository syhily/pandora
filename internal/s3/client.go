package s3

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/qingstor/go-mime"

	"github.com/syhily/pandora/internal/config"
)

// Client encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the sync command.
// It contains client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type Client struct {
	client *s3.Client
	bucket string
}

// newHTTPClient creates an HTTP client with HTTP/2 disabled to avoid protocol errors
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		// Force HTTP/1.1 to avoid HTTP/2 protocol errors
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Minute,
	}
}

// NewClient creates a new S3 client from configuration
func NewClient(cfg *config.S3Configuration) *Client {
	// Create HTTP client with HTTP/2 disabled
	httpClient := newHTTPClient()

	// Create AWS config with custom HTTP client
	awsCfg := aws.Config{
		Region:      cfg.Region,
		Credentials: cfg,
		HTTPClient:  httpClient,
	}

	// Override region if endpoint is provided (for custom S3-compatible services)
	if cfg.Endpoint != "" {
		awsCfg.Region = "auto"
	}

	var s3Client *s3.Client
	if cfg.Endpoint == "" {
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
				return smithyhttp.AddContentChecksumMiddleware(stack)
			})
		})
	} else {
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
				return smithyhttp.AddContentChecksumMiddleware(stack)
			})
		})
	}
	return &Client{
		client: s3Client,
		bucket: cfg.Bucket,
	}
}

// UploadObject reads from a file and puts the data into an object in a bucket.
func (c *Client) UploadObject(ctx context.Context, objectKey string, content []byte) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(content),
		ContentType:   aws.String(mime.DetectFileExt(objectKey[strings.LastIndex(objectKey, ".")+1:])),
		ContentLength: aws.Int64(int64(len(content))),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", c.bucket)
		} else {
			log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n", c.bucket, objectKey, err)
		}
	} else {
		err = s3.NewObjectExistsWaiter(c.client).
			Wait(ctx, &s3.HeadObjectInput{Bucket: aws.String(c.bucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
		}
	}
	return err
}

// ListObjects lists the objects in a bucket.
func (c *Client) ListObjects(ctx context.Context, objectKey string) ([]types.Object, error) {
	var err error
	var output *s3.ListObjectsV2Output
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(objectKey),
	}
	var objects []types.Object
	objectPaginator := s3.NewListObjectsV2Paginator(c.client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", c.bucket)
				err = noBucket
			}
			break
		} else {
			objects = append(objects, output.Contents...)
		}
	}
	return objects, err
}
