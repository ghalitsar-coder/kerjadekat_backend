package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
)

// Mock stores objects in-memory metadata only; returns deterministic bucket keys.
type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Store(ctx context.Context, bucket, filename string, r io.Reader, contentType string, size int64) (domain.StoredObject, error) {
	_ = ctx
	if bucket == "" {
		return domain.StoredObject{}, fmt.Errorf("bucket required")
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExt(contentType)
	}
	key := fmt.Sprintf("buckets/%s/%s%s", bucket, uuid.New().String(), ext)
	if _, err := io.Copy(io.Discard, r); err != nil {
		return domain.StoredObject{}, err
	}
	return domain.StoredObject{Key: key, ContentType: contentType, Size: size}, nil
}

func guessExt(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}

var _ domain.FileStorage = (*Mock)(nil)
