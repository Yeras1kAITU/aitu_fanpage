package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/config"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/handlers"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
	mongorepo "github.com/Yeras1kAITU/aitu_fanpage/internal/repository/mongo"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type App struct {
	cfg      *config.Config
	router   *chi.Mux
	db       *mongo.Database
	handlers *handlers.HandlerContainer
}

func New(cfg *config.Config) (*App, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database.Name)

	var postRepo repository.PostRepository = mongorepo.NewPostRepository(db)
	var userRepo repository.UserRepository = mongorepo.NewUserRepository(db)
	var commentRepo repository.CommentRepository = mongorepo.NewCommentRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	postService := service.NewPostService(postRepo, userRepo, commentRepo)
	commentService := service.NewCommentService(commentRepo, userRepo, postRepo)
	fileService := service.NewFileService(cfg.Upload)
	userService := service.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService, fileService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(userService)
	adminHandler := handlers.NewAdminHandler(postService, userService, commentService)
	mediaHandler := handlers.NewMediaHandler(fileService, postService)

	app := &App{
		cfg: cfg,
		db:  db,
		handlers: &handlers.HandlerContainer{
			Auth:    authHandler,
			Post:    postHandler,
			Comment: commentHandler,
			User:    userHandler,
			Admin:   adminHandler,
			Media:   mediaHandler,
		},
	}

	app.setupRouter(authService.GetTokenAuth())

	return app, nil
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (a *App) Router() *chi.Mux {
	return a.router
}
