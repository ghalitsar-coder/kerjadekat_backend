package repository

import (
	"context"
	"errors"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletPostgres struct {
	db *gorm.DB
}

func NewWalletPostgres(db *gorm.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

func (r *WalletPostgres) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	var w domain.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *WalletPostgres) CreateWallet(ctx context.Context, w *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WalletPostgres) CreateTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *WalletPostgres) UpdateBalance(ctx context.Context, walletID uuid.UUID, balance float64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", balance).Error
}

func (r *WalletPostgres) ListTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]domain.WalletTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []domain.WalletTransaction
	if err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

var _ domain.WalletRepository = (*WalletPostgres)(nil)
