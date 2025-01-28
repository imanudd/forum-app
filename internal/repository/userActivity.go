package repository

import (
	"context"
	"database/sql"

	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/utils"
)

type UserActivityRepositoryImpl interface {
	GetByUserIdAndPostId(ctx context.Context, userId int64, postId int64) (*domain.UserActivity, error)
	UpdateUserActivity(ctx context.Context, userActivity *domain.UserActivity) error
	CreateUserActivity(ctx context.Context, userActivity *domain.UserActivity) error
}

type UserActivityRepository struct {
	db *sql.DB
}

func NewUserActivityRepository(db *sql.DB) UserActivityRepositoryImpl {
	return &UserActivityRepository{
		db: db,
	}
}

func (r *UserActivityRepository) CreateUserActivity(ctx context.Context, userActivity *domain.UserActivity) error {
	insertQuery := `INSERT INTO user_activities (user_id, post_id,is_liked, created_at, created_by) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, insertQuery, userActivity.PostId, userActivity.UserId, userActivity.IsLiked, userActivity.CreatedAt, userActivity.CreatedBy)
	return err
}

func (r *UserActivityRepository) UpdateUserActivity(ctx context.Context, userActivity *domain.UserActivity) error {
	updateQuery := `UPDATE user_activities SET is_liked = ?, updated_at = ?, updated_by = ? WHERE user_id = ? AND post_id = ?`

	_, err := r.db.ExecContext(ctx, updateQuery, userActivity.IsLiked, userActivity.UpdatedAt, userActivity.UpdatedBy, userActivity.UserId, userActivity.PostId)
	return err
}

func (r *UserActivityRepository) GetByUserIdAndPostId(ctx context.Context, userId int64, postId int64) (*domain.UserActivity, error) {
	var userActivity domain.UserActivity

	sqlQuery := `SELECT * FROM user_activities WHERE user_id = ? AND post_id = ?`

	row := r.db.QueryRowContext(ctx, sqlQuery, userId, postId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := utils.ScanToStruct(row, &userActivity); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &userActivity, nil
}
