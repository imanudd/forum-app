package cmd

import (
	"log"

	"github.com/imanudd/forum-app/config"
	rest "github.com/imanudd/forum-app/internal/delivery/http"
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

		//init usecase
		authUseCase := usecase.NewAuthUseCase(cfg, userRepo)

		route := &rest.Route{
			Config:      cfg,
			App:         app,
			AuthUseCase: authUseCase,
			UserRepo:    userRepo,
		}

		route.RegisterRoutes()

		if err := rest.Serve(app, cfg); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}

	},
}
