package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client provides HTTP client functionality
type Client struct {
	client *http.Client
}

// NewClient creates a new HTTP client with proper configuration
func NewClient() *Client {
	// Create transport with HTTP/1.1 to avoid HTTP/2 protocol errors
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

	return &Client{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Minute, // Long timeout for large file downloads
		},
	}
}

// DownloadFile downloads a file from the given URL
func (c *Client) DownloadFile(url string) ([]byte, error) {
	// Send request
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get content length if available
	contentLength := resp.ContentLength
	hasLength := contentLength > 0

	var data []byte
	buf := make([]byte, 32*1024) // 32KB buffer
	var downloaded int64

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
			downloaded += int64(n)

			if hasLength {
				percent := float64(downloaded) / float64(contentLength) * 100
				fmt.Printf("\rDownloading... %.2f%% (%d/%d bytes)", percent, downloaded, contentLength)
			} else {
				fmt.Printf("\rDownloading... %d bytes", downloaded)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}
	}

	fmt.Println("\nDownload complete!")
	return data, nil
}

// PostJSON sends a POST request with JSON payload
func (c *Client) PostJSON(url string, payload []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	// Set custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}
