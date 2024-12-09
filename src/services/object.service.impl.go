package services

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectServiceMinioImpl struct {
	client *minio.Client
}

func NewObjectServiceMinioImpl(client *minio.Client) ObjectService {
	return &ObjectServiceMinioImpl{
		client: client,
	}
}

func (s *ObjectServiceMinioImpl) Upload(ctx context.Context, bucket string, path string, data io.Reader) error {
	_, err := s.client.PutObject(ctx, bucket, path, data, -1, minio.PutObjectOptions{})
	return err
}

func (s *ObjectServiceMinioImpl) Download(ctx context.Context, bucket string, path string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ObjectServiceMinioImpl) SignedUrl(ctx context.Context, bucket string, path string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, bucket, path, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *ObjectServiceMinioImpl) UploadUrl(ctx context.Context, bucket string, path string) (string, error) {
	url, err := s.client.PresignedPutObject(ctx, bucket, path, time.Minute*5)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
