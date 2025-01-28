package cmd

import (
	"log"

	"github.com/imanudd/forum-app/config"
	rest "github.com/imanudd/forum-app/internal/delivery/http"
	"github.com/imanudd/forum-app/internal/delivery/http/middleware"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/internal/usecase"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use: "rest",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *config.Config

		if err := config.InitConfig(); err != nil {
			log.Fatalln(" failed to initialize config", err.Error())
		}

		cfg = config.Get()
		mySQL := NewMysql(cfg)

		app := rest.NewRest(cfg)

		//init repo
		userRepo := repository.NewUserRepository(mySQL)
		postRepo := repository.NewPostRepository(mySQL)
		commentRepo := repository.NewCommentRepository(mySQL)
		userActivityRepo := repository.NewUserActivityRepository(mySQL)
		refreshTokenRepo := repository.NewRefreshTokenRepository(mySQL)

		//init usecase
		authUseCase := usecase.NewAuthUseCase(cfg, userRepo, refreshTokenRepo)
		postUseCase := usecase.NewPostUseCase(cfg, postRepo, commentRepo, userActivityRepo)

		//init middleware
		authMiddleware := middleware.NewAuthMiddleware(cfg, userRepo)
		auth := authMiddleware.JWTAuth()

		route := &rest.Route{
			Config:         cfg,
			App:            app,
			AuthMiddleware: auth,
			AuthUseCase:    authUseCase,
			UserRepo:       userRepo,
			PostUseCase:    postUseCase,
		}

		route.RegisterRoutes()

		if err := rest.Serve(app, cfg); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}

	},
}
