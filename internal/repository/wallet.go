package repository

import (
	"context"
	"fmt"
	"github.com/alob-mtc/wallet-engine/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IWalletRepository interface {
	CreateWallet(ctx context.Context, w *model.Wallet) error
	UpdateWalletState(ctx context.Context, wallet *model.Wallet) error
	GetWallet(ctx context.Context, id, customerId string) (*model.Wallet, error)
	GetWalletById(ctx context.Context, id string) (*model.Wallet, error)
	WithTx(tx *gorm.DB) IWalletRepository
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) IWalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWallet(ctx context.Context, w *model.Wallet) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *walletRepository) UpdateWalletState(ctx context.Context, wallet *model.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *walletRepository) GetWalletById(ctx context.Context, id string) (*model.Wallet, error) {
	var wallet model.Wallet

	if err := r.db.WithContext(ctx).Model(&model.Wallet{}).Where("id = ?", id).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet does not exit")
		}
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) GetWallet(ctx context.Context, id, customerId string) (*model.Wallet, error) {
	var wallet model.Wallet

	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Model(&model.Wallet{}).Where("id = ? AND customer_id = ?", id, customerId).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet does not exit")
		}
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) WithTx(tx *gorm.DB) IWalletRepository {
	return NewWalletRepository(tx)
}

var _ IWalletRepository = &walletRepository{}
