package handler

import (
	"github.com/imanudd/forum-app/internal/usecase"
)

type Handler struct {
	AuthUseCase usecase.AuthUseCaseImpl
}

func NewHandler(useCase *Handler) *Handler {
	return &Handler{
		AuthUseCase: useCase.AuthUseCase,
	}
}
