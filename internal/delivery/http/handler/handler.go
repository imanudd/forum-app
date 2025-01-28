package handler

import (
	"github.com/imanudd/forum-app/internal/usecase"
)

type Handler struct {
	AuthUseCase usecase.AuthUseCaseImpl
	PostUseCase usecase.PostUseCaseImpl
}

func NewHandler(usecase *Handler) *Handler {
	return &Handler{
		AuthUseCase: usecase.AuthUseCase,
		PostUseCase: usecase.PostUseCase,
	}
}
