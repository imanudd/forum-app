package domain

import "time"

type UpsertUserActivity struct {
	PostId  int64 `json:"-"`
	IsLiked bool  `json:"isLiked" validate:"required"`
}

type UserActivity struct {
	Id        int64     `json:"id"`
	PostId    int64     `json:"postId"`
	UserId    int64     `json:"userId"`
	IsLiked   bool      `json:"isLiked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy"`
}
