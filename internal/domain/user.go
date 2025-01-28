package domain

import "time"

type LoginResponse struct {
	Username     string `json:"username"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetByUsernameOrEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type User struct {
	Id        int64      `db:"id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	CreatedBy string     `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
}

func (User) TableName() string {
	return "users"
}
