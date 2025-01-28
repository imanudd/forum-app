package domain

import "time"

type CreateCommentRequest struct {
	PostId         int64  `json:"-"`
	CommentContent string `json:"commentContent" validate:"required"`
}

type Comment struct {
	Id             int64     `json:"id"`
	PostId         int64     `json:"postId"`
	UserId         int64     `json:"userId"`
	CommentContent string    `json:"commentContent"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedBy      string    `json:"createdBy"`
	UpdatedBy      *string   `json:"updatedBy"`
}
