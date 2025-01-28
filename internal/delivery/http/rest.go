package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/imanudd/forum-app/docs"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/delivery/http/handler"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/internal/usecase"
)

// NewRest
// @title Forum Service API
// @version 1.0
// @description Forum Service API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name authorization

func NewRest(cfg *config.Config) *gin.Engine {
	if cfg.Service.Environment != "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	return app
}

func Serve(app *gin.Engine, cfg *config.Config) (err error) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Service.Port),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %s\n", err)
		}
	}()

	log.Println("-------------------------------------------")
	log.Println("server started")
	log.Printf("running on port %s\n", cfg.Service.Port)
	log.Println("-------------------------------------------")

	return gracefulShutdown(server)
}

func gracefulShutdown(srv *http.Server) error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error while shutting down Server. Initiating force shutdown...")
		return err
	}

	log.Println("Server exiting...")

	return nil
}

type Route struct {
	Config         *config.Config
	App            *gin.Engine
	AuthMiddleware gin.HandlerFunc
	AuthUseCase    usecase.AuthUseCaseImpl
	UserRepo       repository.UserRepositoryImpl
	PostUseCase    usecase.PostUseCaseImpl
}

func (r *Route) RegisterRoutes() {
	r.App.Use(gin.Recovery())

	handler := handler.NewHandler(&handler.Handler{
		AuthUseCase: r.AuthUseCase,
		PostUseCase: r.PostUseCase,
	})

	forumsvc := r.App.Group("/forumsvc")
	forumsvc.POST("/auth/signup", handler.SignUp)
	forumsvc.POST("/auth/login", handler.Login)

	forumsvc.Use(r.AuthMiddleware)
	forumsvc.POST("/refresh-token/validate", handler.ValidateRefreshToken)

	postGroup := forumsvc.Group("posts")
	postGroup.GET("", handler.GetListPost)
	postGroup.POST("", handler.CreatePost)
	postGroup.POST("/comment/:postId", handler.CreateCommentOnPost)
	postGroup.PUT("/user-activity/:postId", handler.UpsertUserActivity)
}
