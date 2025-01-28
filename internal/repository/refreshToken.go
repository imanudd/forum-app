package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/imanudd/forum-app/internal/domain"
)

type RefreshTokenRepositoryImpl interface {
	CreateRefreshToken(ctx context.Context, req *domain.RefreshToken) error
	GetLatest(ctx context.Context, userId int64) (*domain.RefreshToken, error)
}

type RefreshToken struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepositoryImpl {
	return &RefreshToken{db: db}
}

func (r *RefreshToken) CreateRefreshToken(ctx context.Context, req *domain.RefreshToken) error {
	insertQuery := `
		INSERT INTO refresh_tokens (user_id, refresh_token, expired_at, created_at, created_by, updated_at, updated_by) 
		values (? ,? ,? ,? ,? ,? ,?)`

	_, err := r.db.ExecContext(ctx, insertQuery, req.UserId, req.RefreshToken, req.ExpiredAt, req.CreatedAt, req.CreatedBy, req.UpdatedAt, req.UpdatedBy)
	return err
}

func (r *RefreshToken) GetLatest(ctx context.Context, userId int64) (*domain.RefreshToken, error) {
	var data domain.RefreshToken

	selectQuery := ` 
		SELECT user_id, refresh_token, expired_at, created_at, created_by, updated_at, updated_by 
		FROM refresh_tokens 
		WHERE user_id = ? AND expired_at >= ? `

	row := r.db.QueryRowContext(ctx, selectQuery, userId, time.Now())
	err := row.Scan(&data.UserId, &data.RefreshToken, &data.ExpiredAt, &data.CreatedAt, &data.CreatedBy, &data.UpdatedAt, &data.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}
