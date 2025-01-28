package repository

import (
	"context"
	"database/sql"

	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/utils"
)

type UserRepositoryImpl interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, req *domain.GetByUsernameOrEmail) (*domain.User, error)
	CreateUser(ctx context.Context, req *domain.User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	sqlQuery := `SELECT * FROM users WHERE id = ?`

	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := utils.ScanToStruct(row, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, req *domain.GetByUsernameOrEmail) (*domain.User, error) {
	var user domain.User

	sqlQuery := `SELECT id,username,email,password,created_at,created_by,updated_at,updated_by FROM users WHERE username = ? OR email = ?`

	row := r.db.QueryRowContext(ctx, sqlQuery, req.Username, req.Email)

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.CreatedBy, &user.UpdatedAt, &user.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, req *domain.User) error {
	insertQuery := `INSERT INTO users (username, email, password, created_at, created_by) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, insertQuery, req.Username, req.Email, req.Password, req.CreatedAt, req.CreatedBy)
	return err
}
