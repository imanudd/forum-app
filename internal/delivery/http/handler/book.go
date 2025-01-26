package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/internal/delivery/http/helper"
	"github.com/imanudd/forum-app/internal/domain"
)

// AddBook handler
// @Summary add book
// @Description add book
// @Tags book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body domain.CreateBookRequest true "data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/book [POST]
func (h *Handler) AddBook(c *gin.Context) {
	var req domain.CreateBookRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err := h.BookUseCase.AddBook(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)
}

// UpdateBook handler
// @Summary update book
// @Description update book
// @Tags book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param id path string true "book id"
// @Param input body domain.UpdateBookRequest true "data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/book/{id} [POST]
func (h *Handler) UpdateBook(c *gin.Context) {
	var req domain.UpdateBookRequest

	err := c.ShouldBind(&req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	req.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err = h.BookUseCase.UpdateBook(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK)
}

// DeleteBook handler
// @Summary delete book
// @Description delete book
// @Tags book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param id path string true "book id"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/book/{id} [DELETE]
func (h *Handler) DeleteBook(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err = h.BookUseCase.DeleteBook(c, id)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK)
}

// GetDetailBook handler
// @Summary get detail book
// @Description get detail book
// @Tags book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param id path string true "book id"
// @Success 200 {object} helper.JSONResponse{data=domain.DetailBook}
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/book/{id} [GET]
func (h *Handler) GetDetailBook(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	resp, err := h.BookUseCase.GetDetailBook(c, id)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK, resp)
}
