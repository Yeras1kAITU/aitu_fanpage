package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
)

func (a *App) setupRouter(authMid *middleware.AuthMiddleware) {
	r := chi.NewRouter()

	allowedOrigins := []string{"*"}
	if a.cfg.Server.Env == "production" {
		allowedOrigins = []string{"https://aitufanpage-production.up.railway.app/"}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.CleanPath)
	r.Use(chimiddleware.Heartbeat("/ping"))

	r.Get("/health", a.healthCheck)

	r.Route("/api", func(r chi.Router) {

		r.Post("/auth/register", a.handlers.Auth.Register)
		r.Post("/auth/login", a.handlers.Auth.Login)

		r.Get("/posts/pinned", a.handlers.Post.GetPinnedPosts)
		r.Get("/posts/featured", a.handlers.Post.GetFeaturedPosts)
		r.Get("/posts/popular", a.handlers.Post.GetPopularPosts)
		r.Get("/posts/categories/stats", a.handlers.Post.GetCategoriesStats)
		r.Get("/posts/search", a.handlers.Post.SearchPosts)
		r.Get("/posts/feed", a.handlers.Post.GetFeed)

		r.Get("/posts", a.handlers.Post.GetPosts)

		r.Route("/posts/{id}", func(r chi.Router) {
			r.Get("/", a.handlers.Post.GetPost)
			r.Group(func(r chi.Router) {
				r.Use(authMid.Authenticator)
				r.Put("/", a.handlers.Post.UpdatePost)
				r.Delete("/", a.handlers.Post.DeletePost)
				r.Post("/like", a.handlers.Post.LikePost)
				r.Delete("/like", a.handlers.Post.UnlikePost)
				r.Get("/likes", a.handlers.Post.GetPostLikes)
				r.Post("/pin", a.handlers.Post.PinPost)
				r.Delete("/pin", a.handlers.Post.UnpinPost)
				r.Post("/feature", a.handlers.Post.FeaturePost)
				r.Delete("/feature", a.handlers.Post.UnfeaturePost)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", a.handlers.Comment.CreateComment)
					r.Get("/", a.handlers.Comment.GetComments)
					r.Get("/count", a.handlers.Comment.GetCommentCount)
				})
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(authMid.Authenticator)

			r.Post("/posts", a.handlers.Post.CreatePost)

			r.Route("/users", func(r chi.Router) {
				r.Get("/me", a.handlers.Auth.GetProfile)
				r.Put("/me", a.handlers.Auth.UpdateProfile)
				r.Put("/me/password", a.handlers.Auth.ChangePassword)
				r.Get("/{id}", a.handlers.User.GetUserProfile)
				r.Get("/{id}/stats", a.handlers.User.GetUserStats)
			})

			r.Route("/comments/{id}", func(r chi.Router) {
				r.Put("/", a.handlers.Comment.UpdateComment)
				r.Delete("/", a.handlers.Comment.DeleteComment)
			})

			r.Post("/media/upload", a.handlers.Media.UploadMedia)
			r.Get("/media/info/{url}", a.handlers.Media.GetMediaInfo)

			r.Route("/admin", func(r chi.Router) {
				r.Use(authMid.RequireRole(models.RoleAdmin))
				r.Get("/stats", a.handlers.Admin.GetSystemStats)
				r.Get("/users", a.handlers.Admin.GetAllUsers)
				r.Get("/users/search", a.handlers.Admin.SearchUsers)
				r.Route("/users/{id}", func(r chi.Router) {
					r.Put("/role", a.handlers.Admin.UpdateUserRole)
					r.Put("/status/{action}", a.handlers.Admin.ToggleUserStatus)
					r.Delete("/", a.handlers.Admin.DeleteUser)
				})
			})
		})
	})

	r.Get("/uploads/*", a.handlers.Media.ServeMedia)

	fs := http.FileServer(http.Dir("./frontend"))

	r.Handle("/*", fs)

	a.router = r
}
