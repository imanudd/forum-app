package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/internal/delivery/http/helper"
	"github.com/imanudd/forum-app/internal/domain"
)

func (h *Handler) ValidateRefreshToken(c *gin.Context) {
	var req domain.ValidateRefreshTokenRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	resp, err := h.AuthUseCase.ValidateRefreshToken(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK, resp)
}

// Login handler
// @Summary login user
// @Description login user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.LoginRequest true "login data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /forumsvc/auth/login [POST]
func (h *Handler) Login(c *gin.Context) {
	var req *domain.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	resp, err := h.AuthUseCase.Login(c, req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK, resp)
}

// Register handler
// @Summary register user
// @Description register user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.SignUpRequest true "register data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /forumsvc/auth/signup [POST]
func (h *Handler) SignUp(c *gin.Context) {
	var req *domain.SignUpRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err := h.AuthUseCase.SignUp(c, req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK)
}
