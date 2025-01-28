package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/internal/delivery/http/helper"
	"github.com/imanudd/forum-app/internal/domain"
)

func (h *Handler) GetListPost(c *gin.Context) {
	pagination, err := helper.SetPagination(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var req domain.GetListPostRequest
	req.Pagination = pagination

	datas, err := h.PostUseCase.GetListPost(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK, datas)
}

// CreatePost handler
// @Summary create post for user
// @Description create post for user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.CreatePost true "create post data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /forumsvc/posts [POST]
func (h *Handler) CreatePost(c *gin.Context) {
	var req domain.CreatePostRequest

	if err := c.BindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.PostUseCase.CreatePost(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)
}

func (h *Handler) CreateCommentOnPost(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var req domain.CreateCommentRequest

	if err := c.BindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.PostId = postId

	err = h.PostUseCase.CreateCommentOnPost(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)
}

func (h *Handler) UpsertUserActivity(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var req domain.UpsertUserActivity

	if err := c.BindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.PostId = postId

	err = h.PostUseCase.UpsertUserActivity(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK)
}
