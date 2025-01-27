package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"

// 	"github.com/imanudd/forum-app/pkg/auth"
// )

// type TransactionRepositoryImpl interface {
// 	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
// }

// type TransactionRepository struct {
// 	db *sql.DB
// }

// func NewTransactionRepository(db *sql.DB) TransactionRepositoryImpl {
// 	return &TransactionRepository{
// 		db: db,
// 	}
// }

// func (r *TransactionRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
// 	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	txCtx := auth.SetTrx(ctx, tx)

// 	defer func() {
// 		if recover() != nil || ctx.Done() != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	err := fn(txCtx)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error db %v", err)
// 	}

// 	return tx.Commit().Error
// }
