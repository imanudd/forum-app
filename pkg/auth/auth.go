package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
)

const (
	userKey  = "user-ctx"
	tokenKey = "token-ctx"
)

type AuthMiddleware interface {
	VerifyToken(tokenStr string) (userID int64, err error)
	GenerateToken(user *domain.User) (string, error)
}

type AuthJwt struct {
	config *config.Config
}

func NewAuth(cfg *config.Config) AuthMiddleware {
	return &AuthJwt{
		config: cfg,
	}
}

type MyClaims struct {
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (a AuthJwt) GenerateToken(user *domain.User) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    a.config.Service.Name,
			ExpiresAt: time.Now().Add(time.Duration(1) * time.Hour).Unix(),
		},
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.config.Service.SignatureKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthJwt) VerifyToken(tokenStr string) (userID int64, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.config.Service.SignatureKey), nil
	})
	if err != nil {
		return userID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, fmt.Errorf("invalid token")
	}

	switch v := claims["user_id"].(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		userID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return
}

func SetUserContext(c *gin.Context, users *domain.User) {
	c.Set(userKey, users)
}
func SetTokenContext(c *gin.Context, token string) {
	c.Set(tokenKey, token)
}

func SetTrx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, "tx", tx)
}

func GetTokenContext(ctx context.Context) string {
	raw, ok := ctx.Value(tokenKey).(string)
	if ok {
		return raw
	}

	return ""
}

func GetUserContext(ctx context.Context) *domain.User {
	raw, ok := ctx.Value(userKey).(*domain.User)
	if ok {
		return raw
	}
	return nil
}

func GetTxContext(ctx context.Context) *sql.Tx {
	raw, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		return raw
	}
	return nil
}
