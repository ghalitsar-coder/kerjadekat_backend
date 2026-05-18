package ocr

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"kerjadekat/backend/internal/domain"
)

// Mock simulates fast KTP OCR with a valid-looking 16-digit NIK.
type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) ExtractKTP(ctx context.Context, objectKey, hintFullName string) (domain.KTPExtraction, error) {
	_ = objectKey
	select {
	case <-ctx.Done():
		return domain.KTPExtraction{}, ctx.Err()
	case <-time.After(50 * time.Millisecond):
	}

	nik, err := randomNIK()
	if err != nil {
		return domain.KTPExtraction{}, err
	}
	name := hintFullName
	if name == "" {
		name = "Nama Terdeteksi OCR"
	}
	return domain.KTPExtraction{NIK: nik, FullName: name}, nil
}

func randomNIK() (string, error) {
	const prefix = "3174"
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%012d", prefix, n.Int64()), nil
}

var _ domain.OCRService = (*Mock)(nil)
