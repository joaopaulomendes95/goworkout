package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joaodias/golang-svelte-fullstack/internal/domain"
	"github.com/joaodias/golang-svelte-fullstack/internal/handler"
	authmw "github.com/joaodias/golang-svelte-fullstack/internal/middleware"
)

func New(userHandler *handler.UserHandler, authMiddleware *authmw.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/users", userHandler.Register)
		r.Post("/tokens/authentication", userHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/me", userHandler.Me)
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				user, _ := r.Context().Value("user").(*domain.User)
				w.Write([]byte("Hello, " + user.Username + "!"))
			})
		})
	})

	return r
}
