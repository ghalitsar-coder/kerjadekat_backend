package domain

import (
	"context"
	"io"
	"time"
)

type StoredObject struct {
	Key         string
	ContentType string
	Size        int64
}

type FileStorage interface {
	Store(ctx context.Context, bucket, filename string, r io.Reader, contentType string, size int64) (StoredObject, error)

	PresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)

	Delete(ctx context.Context, bucket, key string) error
}

const (
	BucketKTP      = "ktp"
	BucketProfiles = "profiles"

	PresignedURLExpiry = 15 * time.Minute
)
