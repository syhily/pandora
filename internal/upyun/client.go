package upyun

import (
	"bytes"

	"github.com/qingstor/go-mime"
	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/syhily/pandora/internal/config"
)

var (
	client *upyun.UpYun
)

func InitUpyunClient() {
	if client == nil {
		client = upyun.NewUpYun(&upyun.UpYunConfig{
			Bucket:   config.GetConfig().Upyun.Bucket,
			Operator: config.GetConfig().Upyun.Operator,
			Password: config.GetConfig().Upyun.Password,
		})
	}
}

func Upload(path string, data []byte) error {
	InitUpyunClient()
	return client.Put(&upyun.PutObjectConfig{
		Path:    path,
		Reader:  bytes.NewReader(data),
		Headers: map[string]string{"Content-Length": mime.DetectFilePath(path)},
	})
}
