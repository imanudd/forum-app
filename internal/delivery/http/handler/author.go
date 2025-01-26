package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/internal/delivery/http/helper"
	"github.com/imanudd/forum-app/internal/domain"
)

func (h *Handler) CreateAuthorAndBook(c *gin.Context) {
	var req domain.CreateAuthorAndBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err := h.AuthorUseCase.CreateAuthorAndBook(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)

}

// CreateAuthor handler
// @Summary create new author
// @Description create new author
// @Tags author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body domain.CreateAuthorRequest true "data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/author [POST]
func (h *Handler) CreateAuthor(c *gin.Context) {
	var req domain.CreateAuthorRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(req)
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err := h.AuthorUseCase.CreateAuthor(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)

}

// AddAuthorBook handler
// @Summary add author book
// @Description add author book
// @Tags author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "author id"
// @Param input body domain.AddAuthorBookRequest true "data"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/author/{id} [POST]
func (h *Handler) AddAuthorBook(c *gin.Context) {

	var req domain.AddAuthorBookRequest
	err := c.ShouldBind(&req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	req.AuthorID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err = h.AuthorUseCase.AddAuthorBook(c, &req)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusCreated)
}

// GetListBookByAuthor handler
// @Summary get list book by author
// @Description get list book by author
// @Tags author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "author id"
// @Success 200 {object} helper.JSONResponse{data=[]domain.Book}
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/author/{id}/list [GET]
func (h *Handler) GetListBookByAuthor(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	resp, err := h.AuthorUseCase.GetListBookByAuthor(c, id)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK, resp)
}

// GetListBookByAuthor handler
// @Summary get list book by author
// @Description get list book by author
// @Tags author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "author id"
// @Param bookid path string true "book id"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /inventorysvc/managements/author/{id}/books/{bookid} [DELETE]
func (h *Handler) DeleteBookByAuthor(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	bookId, err := strconv.Atoi(c.Param("bookid"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "error bad request")
		return
	}

	err = h.AuthorUseCase.DeleteBookByAuthor(c, id, bookId)
	if err != nil {
		helper.InternalError(c, err)
		return
	}

	helper.Success(c, http.StatusOK)
}
