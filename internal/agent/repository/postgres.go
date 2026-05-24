package repository

import (
	"context"
	"errors"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentPostgres struct {
	db *gorm.DB
}

func NewAgentPostgres(db *gorm.DB) *AgentPostgres {
	return &AgentPostgres{db: db}
}

func (r *AgentPostgres) RegisterWorker(ctx context.Context, p domain.RegisterWorkerParams) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p.User).Error; err != nil {
			return err
		}
		p.Profile.UserID = p.User.ID
		if err := tx.Create(p.Profile).Error; err != nil {
			return err
		}
		if len(p.Skills) == 0 {
			return nil
		}
		for i := range p.Skills {
			p.Skills[i].WorkerID = p.Profile.ID
		}
		return tx.Create(&p.Skills).Error
	})
}

func (r *AgentPostgres) KelurahanExists(ctx context.Context, id int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Kelurahan{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AgentPostgres) AgentHasTerritory(ctx context.Context, agentID uuid.UUID, kelurahanID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.AgentTerritory{}).
		Where("agent_id = ? AND kelurahan_id = ?", agentID, kelurahanID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AgentPostgres) ListAgentTerritories(ctx context.Context, agentID uuid.UUID) ([]domain.Kelurahan, error) {
	var rows []domain.Kelurahan
	err := r.db.WithContext(ctx).
		Table("kelurahans").
		Joins("INNER JOIN agent_territories at ON at.kelurahan_id = kelurahans.id").
		Where("at.agent_id = ?", agentID).
		Order("kelurahans.name asc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AgentPostgres) FindWorkerByUserID(ctx context.Context, userID uuid.UUID) (*domain.WorkerProfile, error) {
	var p domain.WorkerProfile
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Kelurahan").
		Preload("Skills.Skill").
		First(&p, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *AgentPostgres) ListWorkersByAgentTerritory(ctx context.Context, agentID uuid.UUID) ([]domain.WorkerProfile, error) {
	var rows []domain.WorkerProfile
	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = worker_profiles.user_id").
		Joins("JOIN agent_territories at ON at.kelurahan_id = users.kelurahan_id").
		Where("at.agent_id = ?", agentID).
		Preload("User").
		Preload("User.Kelurahan").
		Preload("Skills.Skill").
		Order("users.created_at desc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

var _ domain.AgentRepository = (*AgentPostgres)(nil)
