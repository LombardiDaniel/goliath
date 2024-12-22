package services

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectServiceS3Impl struct {
	client *s3.Client
}

func NewObjectServiceS3Impl(client *s3.Client) ObjectService {
	return &ObjectServiceS3Impl{
		client: client,
	}
}

func (s *ObjectServiceS3Impl) Upload(ctx context.Context, bucket string, path string, size int64, data io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   data,
	})
	return err
}

func (s *ObjectServiceS3Impl) Download(ctx context.Context, bucket string, path string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *ObjectServiceS3Impl) SignedUrl(ctx context.Context, bucket string, path string, exp time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(exp))

	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func (s *ObjectServiceS3Impl) UploadUrl(ctx context.Context, bucket string, path string, exp time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(exp))

	if err != nil {
		return "", err
	}
	return request.URL, nil
}
