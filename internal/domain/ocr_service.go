package domain

import "context"

// KTPExtraction holds OCR output from an identity document.
type KTPExtraction struct {
	NIK      string
	FullName string
}

// OCRService extracts structured fields from KTP images.
type OCRService interface {
	ExtractKTP(ctx context.Context, objectKey, hintFullName string) (KTPExtraction, error)
}
