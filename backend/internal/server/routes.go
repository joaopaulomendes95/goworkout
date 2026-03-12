package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Common middleware stack
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(s.Middleware.Authenticate)

		r.Get("/workouts/", s.Middleware.RequireUser(s.WorkoutAPI.HandleGetUserWorkouts))
		r.Get("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleGetWorkoutByID))
		r.Post("/workouts/", s.Middleware.RequireUser(s.WorkoutAPI.HandleCreateWorkout))
		r.Put("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleDeleteWorkoutByID))
	})

	r.Get("/health", s.healthHandler)
	r.Post("/users", s.UserAPI.HandleRegisterUser)
	r.Post("/tokens/authentication", s.TokenAPI.HandleCreateToken)

	return r
}
