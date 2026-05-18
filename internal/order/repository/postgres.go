package repository

import (
	"context"
	"errors"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderPostgres struct {
	db *gorm.DB
}

func NewOrderPostgres(db *gorm.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) WithTx(ctx context.Context, fn func(ctx context.Context, rr domain.OrderRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		child := &OrderPostgres{db: tx}
		return fn(ctx, child)
	})
}

func (r *OrderPostgres) Create(ctx context.Context, o *domain.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *OrderPostgres) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var o domain.Order
	if err := r.db.WithContext(ctx).
		Preload("Skill").
		Preload("Consumer").
		Preload("Worker").
		Preload("Logs", func(db *gorm.DB) *gorm.DB { return db.Order("change_time asc") }).
		First(&o, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *OrderPostgres) ListByConsumer(ctx context.Context, consumerID uuid.UUID, limit, offset int) ([]domain.Order, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []domain.Order
	q := r.db.WithContext(ctx).Preload("Skill").Where("consumer_id = ?", consumerID).Order("created_at desc").Limit(limit).Offset(offset)
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *OrderPostgres) ListByWorker(ctx context.Context, workerUserID uuid.UUID, limit, offset int) ([]domain.Order, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []domain.Order
	q := r.db.WithContext(ctx).Preload("Skill").Where("worker_id = ?", workerUserID).Order("created_at desc").Limit(limit).Offset(offset)
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *OrderPostgres) Update(ctx context.Context, o *domain.Order) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(o).Error
}

func (r *OrderPostgres) AppendStatusLog(ctx context.Context, log *domain.OrderStatusLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *OrderPostgres) CreateIncomeRecord(ctx context.Context, rec *domain.IncomeRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

var _ domain.OrderRepository = (*OrderPostgres)(nil)
