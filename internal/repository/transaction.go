package repository

import (
	"context"
	"fmt"

	"github.com/imanudd/forum-app/pkg/auth"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepositoryImpl {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := r.db.Begin()

	txCtx := auth.SetTrx(ctx, tx)

	defer func() {
		if recover() != nil || ctx.Done() != nil {
			tx.Rollback()
		}
	}()

	err := fn(txCtx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error db %v", err)
	}

	return tx.Commit().Error
}
