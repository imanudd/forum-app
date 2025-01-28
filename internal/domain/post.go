package domain

import (
	"time"

	"github.com/imanudd/forum-app/internal/delivery/http/helper"
)

type GetListPostResponse struct {
	ListPost    []*ListPost `json:"listPost"`
	*Pagination `json:"pagination"`
}

type GetListPostRequest struct {
	*helper.Pagination
}

type ListPost struct {
	Id           int64  `json:"id"`
	UserID       int64  `json:"userId"`
	PostTitle    string `json:"postTitle"`
	PostContent  string `json:"postContent"`
	PostHashtags string `json:"postHashtags"`
	Username     string `json:"username"`
}

type CreatePostRequest struct {
	PostTitle    string `json:"postTitle" validate:"required"`
	PostContent  string `json:"postContent" validate:"required"`
	PostHashtags string `json:"postHashtags"`
}

type Post struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	PostTitle    string    `json:"postTitle"`
	PostContent  string    `json:"postContent"`
	PostHashtags string    `json:"postHashtags"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedBy    string    `json:"createdBy"`
	UpdatedBy    *string   `json:"updatedBy"`
}
