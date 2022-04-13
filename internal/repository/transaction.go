package repository

import (
	"context"
	"github.com/alob-mtc/wallet-engine/internal/model"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	Update(ctx context.Context, tnx *model.Transaction) error
	WithTx(tx *gorm.DB) ITransactionRepository
}

type transactionRepo struct {
	db *gorm.DB
}

// NewTransactionRepo will instantiate Transaction Repository
func NewTransactionRepo(db *gorm.DB) ITransactionRepository {
	return &transactionRepo{db: db}
}

func (t *transactionRepo) Create(ctx context.Context, transaction *model.Transaction) error {

	return t.db.WithContext(ctx).Create(transaction).Error
}

func (t *transactionRepo) WithTx(tx *gorm.DB) ITransactionRepository {
	return NewTransactionRepo(tx)
}

func (t *transactionRepo) Update(ctx context.Context, tnx *model.Transaction) error {
	return t.db.WithContext(ctx).Save(tnx).Error
}
