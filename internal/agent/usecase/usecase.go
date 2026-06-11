package usecase

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ktpBucket     = "ktp"
	profileBucket = "profiles"
)

type Agents struct {
	agents domain.AgentRepository
	users  domain.UserRepository
	skills domain.SkillRepository
	store  domain.FileStorage
	ocr    domain.OCRService
}

func NewAgents(
	agents domain.AgentRepository,
	users domain.UserRepository,
	skills domain.SkillRepository,
	store domain.FileStorage,
	ocr domain.OCRService,
) *Agents {
	return &Agents{
		agents: agents,
		users:  users,
		skills: skills,
		store:  store,
		ocr:    ocr,
	}
}

type RegisterWorkerInput struct {
	AgentID         uuid.UUID
	AgentRole       string
	PhoneNumber     string
	FullName        string
	RtRw            string
	KelurahanID     int
	SkillIDs        []int
	KTPPhoto        io.Reader
	KTPFilename     string
	KTPContentType  string
	KTPSize         int64
	ProfilePhoto    io.Reader
	ProfileFilename string
	ProfileContentType string
	ProfileSize     int64
	Latitude        float64
	Longitude       float64
}

type OCRPreview struct {
	NIK      string `json:"nik"`
	FullName string `json:"full_name"`
}

type RegisterWorkerResult struct {
	User          domain.User          `json:"user"`
	WorkerProfile domain.WorkerProfile `json:"worker_profile"`
	OCRPreview    OCRPreview           `json:"ocr_preview"`
}

type AgentWorkerSummary struct {
	UserID        uuid.UUID  `json:"user_id"`
	FullName      string     `json:"full_name"`
	ProfilePhoto  *string    `json:"profile_photo,omitempty"`
	PhoneNumber   *string    `json:"phone_number"`
	Status        string     `json:"status"`
	Availability  string     `json:"availability"`
	KelurahanName string     `json:"kelurahan"`
	RtRw          *string    `json:"rt_rw"`
	RatingAvg     float64    `json:"rating_avg"`
	RatingCount   int        `json:"rating_count"`
	CreatedAt     time.Time  `json:"created_at"`
	VerifiedAt    *time.Time `json:"verified_at"`
	Skills        []string   `json:"skills"`
}

type ListAgentWorkersResult struct {
	Items []AgentWorkerSummary `json:"items"`
}

func (a *Agents) ListWorkers(ctx context.Context, agentID uuid.UUID) (*ListAgentWorkersResult, error) {
	rows, err := a.agents.ListWorkersByAgentTerritory(ctx, agentID)
	if err != nil {
		return nil, err
	}
	items := make([]AgentWorkerSummary, 0, len(rows))
	for _, row := range rows {
		user := row.User
		kelName := ""
		if user.Kelurahan != nil {
			kelName = user.Kelurahan.Name
		}
		skills := make([]string, 0, len(row.Skills))
		for _, s := range row.Skills {
			if s.Skill.Name != "" {
				skills = append(skills, s.Skill.Name)
			}
		}
		items = append(items, AgentWorkerSummary{
			UserID:        user.ID,
			FullName:      user.FullName,
			ProfilePhoto:  user.ProfilePhoto,
			PhoneNumber:   user.PhoneNumber,
			Status:        user.Status,
			Availability:  row.Availability,
			KelurahanName: kelName,
			RtRw:          user.RtRw,
			RatingAvg:     row.RatingAvg,
			RatingCount:   row.RatingCount,
			CreatedAt:     user.CreatedAt,
			VerifiedAt:    user.VerifiedAt,
			Skills:        skills,
		})
	}
	return &ListAgentWorkersResult{Items: items}, nil
}

func (a *Agents) RegisterWorker(ctx context.Context, in RegisterWorkerInput) (*RegisterWorkerResult, error) {
	phone := strings.TrimSpace(in.PhoneNumber)
	name := strings.TrimSpace(in.FullName)
	if phone == "" || len(phone) < 8 || name == "" {
		return nil, domain.ErrInvalidInput
	}
	if in.KelurahanID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	if len(in.SkillIDs) == 0 {
		return nil, domain.ErrInvalidInput
	}
	if in.KTPPhoto == nil || in.ProfilePhoto == nil {
		return nil, domain.ErrInvalidInput
	}

	if _, err := a.users.FindByPhone(ctx, phone); err == nil {
		return nil, domain.ErrConflict
	} else if err != domain.ErrNotFound {
		return nil, err
	}

	exists, err := a.agents.KelurahanExists(ctx, in.KelurahanID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrInvalidInput
	}

	if in.AgentRole == domain.RoleAgent {
		ok, err := a.agents.AgentHasTerritory(ctx, in.AgentID, in.KelurahanID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrForbidden
		}
	}

	for _, skillID := range in.SkillIDs {
		if _, err := a.skills.FindCategoryByID(ctx, skillID); err != nil {
			return nil, err
		}
	}

	ktpObj, err := a.store.Store(ctx, ktpBucket, in.KTPFilename, in.KTPPhoto, in.KTPContentType, in.KTPSize)
	if err != nil {
		return nil, fmt.Errorf("store ktp: %w", err)
	}
	profileObj, err := a.store.Store(ctx, profileBucket, in.ProfileFilename, in.ProfilePhoto, in.ProfileContentType, in.ProfileSize)
	if err != nil {
		return nil, fmt.Errorf("store profile: %w", err)
	}

	extracted, err := a.ocr.ExtractKTP(ctx, ktpObj.Key, name)
	if err != nil {
		return nil, fmt.Errorf("ocr ktp: %w", err)
	}

	nikHash, err := bcrypt.GenerateFromPassword([]byte(extracted.NIK), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(nikHash)

	now := time.Now()
	userID := uuid.New()
	profileID := uuid.New()
	rtRw := strings.TrimSpace(in.RtRw)
	var rtRwPtr *string
	if rtRw != "" {
		rtRwPtr = &rtRw
	}
	kelurahanID := in.KelurahanID

	user := &domain.User{
		ID:           userID,
		PhoneNumber:  &phone,
		FullName:     name,
		Role:         domain.RoleWorker,
		NikHash:      &hashStr,
		ProfilePhoto: &profileObj.Key,
		KtpPhotoRef:  &ktpObj.Key,
		Status:       domain.UserStatusActive,
		VerifiedBy:   &in.AgentID,
		VerifiedAt:   &now,
		RtRw:         rtRwPtr,
		KelurahanID:  &kelurahanID,
	}

	profile := &domain.WorkerProfile{
		ID:           profileID,
		UserID:       userID,
		Availability: domain.WorkerAvailabilityOffline,
	}
	if in.Latitude != 0 && in.Longitude != 0 {
		profile.LastLocation = &domain.NullPoint{
			Lat:   in.Latitude,
			Lng:   in.Longitude,
			Valid: true,
		}
	}

	skills := make([]domain.WorkerSkill, 0, len(in.SkillIDs))
	seen := make(map[int]struct{}, len(in.SkillIDs))
	for _, skillID := range in.SkillIDs {
		if _, dup := seen[skillID]; dup {
			continue
		}
		seen[skillID] = struct{}{}
		skills = append(skills, domain.WorkerSkill{
			ID:       uuid.New(),
			WorkerID: profileID,
			SkillID:  skillID,
			Level:    "beginner",
		})
	}

	if err := a.agents.RegisterWorker(ctx, domain.RegisterWorkerParams{
		User:    user,
		Profile: profile,
		Skills:  skills,
	}); err != nil {
		return nil, err
	}

	loaded, err := a.agents.FindWorkerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	loaded.User.NikHash = nil

	return &RegisterWorkerResult{
		User:          loaded.User,
		WorkerProfile: *loaded,
		OCRPreview:    OCRPreview{NIK: extracted.NIK, FullName: extracted.FullName},
	}, nil
}

func (a *Agents) ListTerritories(ctx context.Context, agentID uuid.UUID) ([]domain.Kelurahan, error) {
	return a.agents.ListAgentTerritories(ctx, agentID)
}
