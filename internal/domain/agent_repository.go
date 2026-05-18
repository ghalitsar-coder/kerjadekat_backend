package domain

import (
	"context"

	"github.com/google/uuid"
)

// RegisterWorkerParams holds rows to insert atomically for agent onboarding.
type RegisterWorkerParams struct {
	User    *User
	Profile *WorkerProfile
	Skills  []WorkerSkill
}

// AgentRepository persists agent-driven worker registration flows.
type AgentRepository interface {
	RegisterWorker(ctx context.Context, p RegisterWorkerParams) error
	FindWorkerByUserID(ctx context.Context, userID uuid.UUID) (*WorkerProfile, error)
	KelurahanExists(ctx context.Context, id int) (bool, error)
	AgentHasTerritory(ctx context.Context, agentID uuid.UUID, kelurahanID int) (bool, error)
	ListAgentTerritories(ctx context.Context, agentID uuid.UUID) ([]Kelurahan, error)
}
