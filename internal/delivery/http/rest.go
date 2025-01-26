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

	_ "github.com/imanudd/forum-app/docs"

	"github.com/gin-gonic/gin"
	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/delivery/http/handler"
	"github.com/imanudd/forum-app/internal/delivery/http/middleware"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/internal/usecase"
)

// NewRest
// @title Inventory Service API
// @version 1.0
// @description Inventory Service API
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

func NewRest(cfg *config.MainConfig) *gin.Engine {
	if cfg.Environment != "production" {
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

func Serve(app *gin.Engine, cfg *config.MainConfig) (err error) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServicePort),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %s\n", err)
		}
	}()

	log.Println("-------------------------------------------")
	log.Println("server started")
	log.Printf("running on port %d\n", cfg.ServicePort)
	log.Println("-------------------------------------------")

	return gracefulShutdown(server)
}

func gracefulShutdown(srv *http.Server) error {
	done := make(chan os.Signal)

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
	Config        *config.MainConfig
	App           *gin.Engine
	AuthUseCase   usecase.AuthUseCaseImpl
	BookUseCase   usecase.BookUseCaseImpl
	AuthorUseCase usecase.AuthorUseCaseImpl
	UserRepo      repository.UserRepositoryImpl
}

func (r *Route) RegisterRoutes() {
	r.App.Use(gin.Recovery())

	auth := middleware.NewAuthMiddleware(r.Config, r.UserRepo)

	handler := handler.NewHandler(&handler.Handler{
		AuthUseCase:   r.AuthUseCase,
		BookUseCase:   r.BookUseCase,
		AuthorUseCase: r.AuthorUseCase,
	})

	inventorySvc := r.App.Group("/inventorysvc")
	inventorySvc.POST("/auth/register", handler.Register)
	inventorySvc.POST("/auth/login", handler.Login)

	inventorySvc.POST("/managements/book", auth.JWTAuth(handler.AddBook))
	inventorySvc.PUT("/managements/book/:id", auth.JWTAuth(handler.UpdateBook))
	inventorySvc.DELETE("/managements/book/:id", auth.JWTAuth(handler.DeleteBook))
	inventorySvc.GET("/managements/book/:id", auth.JWTAuth(handler.GetDetailBook))

	inventorySvc.POST("/managements/author/book", auth.JWTAuth(handler.CreateAuthorAndBook))
	inventorySvc.POST("/managements/author", auth.JWTAuth(handler.CreateAuthor))
	inventorySvc.POST("/managements/author/:id", auth.JWTAuth(handler.AddAuthorBook))
	inventorySvc.GET("/managements/author/:id/list", auth.JWTAuth(handler.GetListBookByAuthor))
	inventorySvc.DELETE("managements/author/:id/books/:bookid", auth.JWTAuth(handler.DeleteBookByAuthor))

}
