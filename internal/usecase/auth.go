package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/pkg/auth"
	"github.com/imanudd/forum-app/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseImpl interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	SignUp(ctx context.Context, req *domain.SignUpRequest) (err error)
	ValidateRefreshToken(ctx context.Context, req *domain.ValidateRefreshTokenRequest) (*domain.ValidateRefreshTokenResponse, error)
}

type authUseCase struct {
	cfg              *config.Config
	userRepo         repository.UserRepositoryImpl
	refreshTokenRepo repository.RefreshTokenRepositoryImpl
}

func NewAuthUseCase(cfg *config.Config, userRepo repository.UserRepositoryImpl, refreshTokenRepo repository.RefreshTokenRepositoryImpl) AuthUseCaseImpl {
	return &authUseCase{
		cfg:              cfg,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (a *authUseCase) ValidateRefreshToken(ctx context.Context, req *domain.ValidateRefreshTokenRequest) (*domain.ValidateRefreshTokenResponse, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var resp domain.ValidateRefreshTokenResponse

	user := auth.GetUserContext(ctx)

	refreshToken, err := a.refreshTokenRepo.GetLatest(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, errors.New("refresh token is expired")
	}

	if refreshToken.RefreshToken != req.RefreshToken {
		return nil, errors.New("refresh token is invalid")
	}

	auth := auth.NewAuth(a.cfg)
	resp.Token, err = auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}

func (a *authUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	response := domain.LoginResponse{
		Username: req.Username,
	}

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
	response.Token, err = auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.refreshTokenRepo.GetLatest(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	if refreshToken != nil {
		response.RefreshToken = refreshToken.RefreshToken
		return &response, nil
	}

	response.RefreshToken = auth.GenerateRefreshToken()
	err = a.refreshTokenRepo.CreateRefreshToken(ctx, &domain.RefreshToken{
		UserId:       user.Id,
		RefreshToken: response.RefreshToken,
		ExpiredAt:    time.Now().Add(10 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    user.Username,
		UpdatedBy:    user.Username,
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *authUseCase) SignUp(ctx context.Context, req *domain.SignUpRequest) (err error) {
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

	return a.userRepo.CreateUser(ctx, &domain.User{
		Username:  req.Username,
		Password:  string(hash),
		Email:     req.Email,
		CreatedAt: time.Now(),
		CreatedBy: req.Username,
	})

}
