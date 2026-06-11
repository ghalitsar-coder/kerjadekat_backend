package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary(cloudURL string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(cloudURL)
	if err != nil {
		return nil, fmt.Errorf("cloudinary new: %w", err)
	}
	return &Cloudinary{cld: cld}, nil
}

func (c *Cloudinary) Store(ctx context.Context, bucket, filename string, r io.Reader, contentType string, size int64) (domain.StoredObject, error) {
	if bucket == "" {
		return domain.StoredObject{}, fmt.Errorf("bucket required")
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExt(contentType)
	}
	publicID := fmt.Sprintf("buckets/%s/%s%s", bucket, uuid.New().String(), ext)

	uploadResult, err := c.cld.Upload.Upload(ctx, r, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	if err != nil {
		return domain.StoredObject{}, fmt.Errorf("cloudinary upload: %w", err)
	}

	return domain.StoredObject{
		Key:         uploadResult.PublicID,
		ContentType: contentType,
		Size:        size,
	}, nil
}

func (c *Cloudinary) PresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	asset, err := c.cld.Admin.Asset(ctx, admin.AssetParams{
		PublicID: key,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary asset: %w", err)
	}
	return asset.SecureURL, nil
}

func (c *Cloudinary) Delete(ctx context.Context, bucket, key string) error {
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: key,
	})
	if err != nil {
		return fmt.Errorf("cloudinary destroy: %w", err)
	}
	return nil
}

var _ domain.FileStorage = (*Cloudinary)(nil)
