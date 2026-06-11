package domain

import (
	"context"

	"github.com/google/uuid"
)

// WalletRepository persists wallets and transactions.
type WalletRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	CreateWallet(ctx context.Context, w *Wallet) error
	CreateTransaction(ctx context.Context, tx *WalletTransaction) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, balance float64) error
	ListTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]WalletTransaction, error)
}
