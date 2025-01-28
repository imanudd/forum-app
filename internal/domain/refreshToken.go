package domain

import "time"

type ValidateRefreshTokenResponse struct {
	Token string `json:"token"`
}

type ValidateRefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshToken struct {
	Id           int64
	UserId       int64
	RefreshToken string
	ExpiredAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
}
