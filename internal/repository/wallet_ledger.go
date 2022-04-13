package repository

import (
	"context"
	"github.com/alob-mtc/wallet-engine/internal/model"
	"gorm.io/gorm"
)

type IWalletLedgerRepository interface {
	Create(ctx context.Context, w *model.WalletLedger) error
	WithTx(tx *gorm.DB) IWalletLedgerRepository
}

type walletLedgerRepository struct {
	db *gorm.DB
}

func NewWalletLedgerRepository(db *gorm.DB) IWalletLedgerRepository {
	return &walletLedgerRepository{db: db}
}

func (r *walletLedgerRepository) Create(ctx context.Context, w *model.WalletLedger) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *walletLedgerRepository) WithTx(tx *gorm.DB) IWalletLedgerRepository {
	return NewWalletLedgerRepository(tx)
}

var _ IWalletLedgerRepository = &walletLedgerRepository{}
