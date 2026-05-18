package usecase

import (
	"context"
	"sort"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/geo"

	"github.com/google/uuid"
)

// NearbyWorkerItem is returned by GET /workers/nearby for map and list UIs.
type NearbyWorkerItem struct {
	UserID       uuid.UUID `json:"user_id"`
	FullName     string    `json:"full_name"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	DistanceM    float64   `json:"distance_m"`
	RatingAvg    float64   `json:"rating_avg"`
	RatingCount  int       `json:"rating_count"`
	BaseRate     *float64  `json:"base_rate,omitempty"`
	Availability string    `json:"availability"`
	VerifiedRT   bool      `json:"verified_rt"`
	Skills       []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"skills"`
}

type NearbyInput struct {
	Latitude      float64
	Longitude     float64
	RadiusMeters  float64
	SkillID       *int
}

func (w *Workers) Nearby(ctx context.Context, in NearbyInput) ([]NearbyWorkerItem, error) {
	if in.Latitude < -90 || in.Latitude > 90 || in.Longitude < -180 || in.Longitude > 180 {
		return nil, domain.ErrInvalidInput
	}
	radius := in.RadiusMeters
	if radius <= 0 {
		radius = w.defaultRadiusM
	}
	if radius <= 0 {
		radius = 5000
	}

	profiles, err := w.repo.FindNearbyOnline(ctx, in.Latitude, in.Longitude, radius, in.SkillID)
	if err != nil {
		return nil, err
	}

	if len(profiles) > 0 {
		_ = w.syncProfilesToPresence(ctx, profiles)
	}

	items := make([]NearbyWorkerItem, 0, len(profiles))
	for _, p := range profiles {
		if p.LastLocation == nil || !p.LastLocation.Valid {
			continue
		}
		lat := p.LastLocation.Lat
		lng := p.LastLocation.Lng
		skills := make([]struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}, 0, len(p.Skills))
		for _, s := range p.Skills {
			skills = append(skills, struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: s.SkillID, Name: s.Skill.Name})
		}
		items = append(items, NearbyWorkerItem{
			UserID:       p.UserID,
			FullName:     p.User.FullName,
			Latitude:     lat,
			Longitude:    lng,
			DistanceM:    geo.HaversineM(in.Latitude, in.Longitude, lat, lng),
			RatingAvg:    p.RatingAvg,
			RatingCount:  p.RatingCount,
			BaseRate:     p.BaseRate,
			Availability: p.Availability,
			VerifiedRT:   p.User.VerifiedAt != nil,
			Skills:       skills,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].DistanceM < items[j].DistanceM
	})
	return items, nil
}

// BootstrapPresence loads online workers from PostgreSQL into Redis GEO for matching.
func (w *Workers) BootstrapPresence(ctx context.Context) error {
	profiles, err := w.repo.ListOnlineWithLocation(ctx)
	if err != nil {
		return err
	}
	return w.syncProfilesToPresence(ctx, profiles)
}

func (w *Workers) syncProfilesToPresence(ctx context.Context, profiles []domain.WorkerProfile) error {
	if w.presence == nil {
		return nil
	}
	for _, p := range profiles {
		if p.Availability != domain.WorkerAvailabilityOnline {
			continue
		}
		if p.LastLocation == nil || !p.LastLocation.Valid {
			continue
		}
		skillIDs := make([]int, 0, len(p.Skills))
		for _, s := range p.Skills {
			skillIDs = append(skillIDs, s.SkillID)
		}
		if err := w.presence.SetAvailability(ctx, p.UserID, domain.WorkerAvailabilityOnline); err != nil {
			return err
		}
		if err := w.presence.UpdateWorkerLocation(ctx, p.UserID, p.LastLocation.Lat, p.LastLocation.Lng, skillIDs); err != nil {
			return err
		}
	}
	return nil
}
