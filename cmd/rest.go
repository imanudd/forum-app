package cmd

import (
	"log"

	"github.com/imanudd/forum-app/config"
	rest "github.com/imanudd/forum-app/internal/delivery/http"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/internal/usecase"
	"github.com/imanudd/forum-app/pkg/elasticsearch"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use: "rest",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()

		pgDB := InitPostgreSQL(cfg)

		if cfg.LogMode {
			pgDB = pgDB.Debug()
		}

		client := InitElastic(cfg)

		//init elasticsearch
		es := elasticsearch.New(client)

		app := rest.NewRest(cfg)

		//init repo
		userRepo := repository.NewUserRepository(pgDB)
		bookRepo := repository.NewBookRepository(pgDB)
		authorRepo := repository.NewAuthorRepository(pgDB)
		trx := repository.NewTransactionRepository(pgDB)

		//init usecase
		authUseCase := usecase.NewAuthUseCase(cfg, trx, userRepo)
		bookUseCase := usecase.NewBookUseCase(cfg, es, trx, bookRepo, authorRepo)
		authorUseCase := usecase.NewAuthorUseCase(cfg, trx, authorRepo, bookRepo)

		route := &rest.Route{
			Config:        cfg,
			App:           app,
			AuthUseCase:   authUseCase,
			BookUseCase:   bookUseCase,
			AuthorUseCase: authorUseCase,
			UserRepo:      userRepo,
		}

		route.RegisterRoutes()

		if err := rest.Serve(app, cfg); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}

	},
}
