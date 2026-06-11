package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client *minio.Client
}

func NewS3(endpoint, accessKey, secretKey, region string, useSSL bool) (*S3, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}
	return &S3{client: client}, nil
}

func (s *S3) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("bucket exists: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("make bucket: %w", err)
		}
	}
	return nil
}

func (s *S3) Store(ctx context.Context, bucket, filename string, r io.Reader, contentType string, size int64) (domain.StoredObject, error) {
	if bucket == "" {
		return domain.StoredObject{}, fmt.Errorf("bucket required")
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExt(contentType)
	}
	key := fmt.Sprintf("buckets/%s/%s%s", bucket, uuid.New().String(), ext)

	if err := s.ensureBucket(ctx, bucket); err != nil {
		return domain.StoredObject{}, err
	}

	_, err := s.client.PutObject(ctx, bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return domain.StoredObject{}, fmt.Errorf("put object: %w", err)
	}

	return domain.StoredObject{Key: key, ContentType: contentType, Size: size}, nil
}

func (s *S3) PresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presigned url: %w", err)
	}
	return url.String(), nil
}

func (s *S3) Delete(ctx context.Context, bucket, key string) error {
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

var _ domain.FileStorage = (*S3)(nil)
