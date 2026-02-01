package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) setupRouter() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", a.healthCheck)

	r.Route("/api", func(r chi.Router) {
		r.Post("/posts", a.handlers.Post.CreatePost)
		r.Get("/posts", a.handlers.Post.GetPosts)
		r.Post("/posts/{id}/like", a.handlers.Post.LikePost)
		r.Get("/feed", a.handlers.Post.GetFeed)
	})

	a.router = r
}
