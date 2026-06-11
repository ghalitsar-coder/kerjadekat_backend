package domain

import (
	"time"

	"github.com/google/uuid"
)

// Wallet maps wallets — one per user.
type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_wallets_user;not null"`
	Balance   float64   `gorm:"type:numeric(12,2);not null;default:0"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// WalletType constants.
const (
	WalletTxTypeCredit = "credit"
	WalletTxTypeDebit  = "debit"
)

// WalletTransaction maps wallet_transactions.
type WalletTransaction struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_wallet_tx_wallet"`
	Type          string     `gorm:"type:varchar(20);not null"`
	Amount        float64    `gorm:"type:numeric(12,2);not null"`
	BalanceBefore float64    `gorm:"type:numeric(12,2);not null"`
	BalanceAfter  float64    `gorm:"type:numeric(12,2);not null"`
	ReferenceType string     `gorm:"type:varchar(30);index:idx_wallet_tx_ref"`
	ReferenceID   *string    `gorm:"type:varchar(100)"`
	Description   *string    `gorm:"type:varchar(255)"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null;autoCreateTime;index:idx_wallet_tx_created"`

	Wallet Wallet `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
