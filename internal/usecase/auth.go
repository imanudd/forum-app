package usecase

import (
	"context"
	"errors"

	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/pkg/auth"
	"github.com/imanudd/forum-app/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseImpl interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (err error)
}

type authUseCase struct {
	cfg      *config.Config
	userRepo repository.UserRepositoryImpl
}

func NewAuthUseCase(cfg *config.Config, userRepo repository.UserRepositoryImpl) AuthUseCaseImpl {
	return &authUseCase{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (a *authUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetByUsernameOrEmail(ctx, &domain.GetByUsernameOrEmail{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user is not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	auth := auth.NewAuth(a.cfg)
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Username: req.Username,
		Token:    token,
	}, nil
}

func (a *authUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (err error) {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}
	user, err := a.userRepo.GetByUsernameOrEmail(ctx, &domain.GetByUsernameOrEmail{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return
	}

	if user != nil {
		return errors.New("user is already exist")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error when hashing password")
	}

	return a.userRepo.RegisterUser(ctx, &domain.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
	})

}
