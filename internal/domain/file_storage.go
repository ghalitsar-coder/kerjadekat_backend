package domain

import (
	"context"
	"io"
)

// StoredObject describes a file persisted in object storage.
type StoredObject struct {
	Key         string
	ContentType string
	Size        int64
}

// FileStorage abstracts cloud object storage (S3/GCS); MVP uses a mock implementation.
type FileStorage interface {
	Store(ctx context.Context, bucket, filename string, r io.Reader, contentType string, size int64) (StoredObject, error)
}
