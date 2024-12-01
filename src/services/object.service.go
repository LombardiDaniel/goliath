package services

import (
	"context"
	"io"
)

type ObjectService interface {
	UploadObject(ctx context.Context, bucket string, path string, data io.Reader) error
	DownloadObject(ctx context.Context, bucket string, path string) ([]byte, error)
	SignedUrl(ctx context.Context, bucket string, path string) (string, error)
}
