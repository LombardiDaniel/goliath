package services

import (
	"context"
	"io"
)

type ObjectService interface {
	Upload(ctx context.Context, bucket string, path string, data io.Reader) error
	Download(ctx context.Context, bucket string, path string) ([]byte, error)
	SignedUrl(ctx context.Context, bucket string, path string) (string, error)
	UploadUrl(ctx context.Context, bucket string, path string) (string, error)
}
