package services

import (
	"context"
	"io"
	"time"
)

// ObjectService defines the interface for object storage operations.
// It provides methods for uploading and downloading objects, as well as
// generating signed URLs for secure access and upload URLs for pre-signed uploads.
type ObjectService interface {
	// Upload uploads an object to the specified bucket and path.
	Upload(ctx context.Context, bucket string, path string, size int64, data io.Reader) error

	// Download retrieves an object from the specified bucket and path.
	Download(ctx context.Context, bucket string, path string) ([]byte, error)

	// SignedUrl generates a signed URL for accessing an object in the specified bucket and path.
	SignedUrl(ctx context.Context, bucket string, path string, exp time.Duration) (string, error)

	// UploadUrl generates a pre-signed URL for uploading an object to the specified bucket and path.
	UploadUrl(ctx context.Context, bucket string, path string, exp time.Duration) (string, error)
}
