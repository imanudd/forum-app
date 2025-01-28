package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/delivery/http/helper"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/pkg/auth"
)

type AuthMiddleware struct {
	cfg  *config.Config
	repo repository.UserRepositoryImpl
}

func NewAuthMiddleware(cfg *config.Config, repo repository.UserRepositoryImpl) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:  cfg,
		repo: repo,
	}
}

func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("authorization")
		if authHeader == "" {
			helper.Error(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		barierToken := strings.Split(authHeader, "Bearer ")
		if len(barierToken) < 2 {
			helper.Error(c, http.StatusUnauthorized, "token not valid")
			return
		}

		token := barierToken[1]

		authJwt := auth.NewAuth(m.cfg)
		userID, err := authJwt.VerifyToken(token)
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := m.repo.GetByID(c, int(userID))
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		auth.SetUserContext(c, user)
		auth.SetTokenContext(c, token)

		c.Next()
	}
}
