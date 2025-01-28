package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/utils"
)

type PostRepositoryImpl interface {
	CountPost(ctx context.Context) (int64, error)
	GetListPost(ctx context.Context, req *domain.GetListPostRequest) ([]*domain.ListPost, error)
	GetById(ctx context.Context, id int64) (*domain.Post, error)
	CreatePost(ctx context.Context, req *domain.Post) error
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepositoryImpl {
	return &PostRepository{db: db}
}

func (r *PostRepository) CountPost(ctx context.Context) (int64, error) {
	var count int64
	countQuery := `SELECT COUNT(*) FROM posts `

	row := r.db.QueryRowContext(ctx, countQuery)

	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (r *PostRepository) GetListPost(ctx context.Context, req *domain.GetListPostRequest) ([]*domain.ListPost, error) {
	var posts []*domain.ListPost

	selectQuery := ` 
	SELECT 
		p.id, p.user_id, p.post_title, p.post_content, p.post_hastags, 
		u.username
	FROM posts p
	JOIN users u ON p.user_id = u.id 
	ORDER BY p.created_at DESC 
	LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, selectQuery, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data domain.ListPost
		if err := rows.Scan(&data.Id, &data.UserID, &data.PostTitle, &data.PostContent, &data.PostHashtags, &data.Username); err != nil {
			return nil, err
		}

		posts = append(posts, &data)
	}

	return posts, nil
}

func (r *PostRepository) GetById(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post

	sqlQuery := `SELECT * FROM posts WHERE id = ?`

	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := utils.ScanToStruct(row, &post); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post is not exist")
		}
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, req *domain.Post) error {
	insertQuery := `INSERT INTO posts (user_id, post_title, post_content, post_hastags, created_at, created_by) VALUES (?, ?, ?, ?, ?,?)`

	_, err := r.db.ExecContext(ctx, insertQuery, req.UserID, req.PostTitle, req.PostContent, req.PostHashtags, req.CreatedAt, req.CreatedBy)
	return err
}
