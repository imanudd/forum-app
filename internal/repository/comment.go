package repository

import (
	"context"
	"database/sql"

	"github.com/imanudd/forum-app/internal/domain"
)

type CommentRepositoryImpl interface {
	CreateComment(ctx context.Context, req *domain.Comment) error
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepositoryImpl {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, req *domain.Comment) error {
	insertQuery := `INSERT INTO comments (user_id, post_id, comment_content, created_at, created_by) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, insertQuery, req.UserId, req.PostId, req.CommentContent, req.CreatedAt, req.CreatedBy)
	return err
}
