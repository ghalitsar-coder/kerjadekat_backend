package domain

import "context"

// KelurahanRepository provides access to kelurahan (administrative area) data.
type KelurahanRepository interface {
	// ListAll returns all kelurahans with their geographic centroids.
	ListAll(ctx context.Context) ([]Kelurahan, error)
}
