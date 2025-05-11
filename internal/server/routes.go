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

	// Public endpoints
	// Health check
	r.Get("/health", s.healthHandler)
	// Create new user
	r.Post("/users", s.UserAPI.HandleRegisterUser)
	// Login user / create token
	r.Post("/tokens/authentication", s.TokenAPI.HandleCreateToken)

	// Grouping this enpdoints to be protected by the authentication middleware
	r.Group(func(r chi.Router) {
		r.Use(s.Middleware.Authenticate)

		// Get user info
		r.Get("/users/me", s.Middleware.RequireUser(s.UserAPI.HandleGetCurrentUser))
		// Update user info
		r.Put("/users/me", s.Middleware.RequireUser(s.UserAPI.HandleUpdateUser))
		// Get all workouts
		r.Get("/workouts/", s.Middleware.RequireUser(s.WorkoutAPI.HandleGetUserWorkouts))
		// Get a specific workout
		r.Get("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleGetWorkoutByID))
		// Create new workout
		r.Post("/workouts/", s.Middleware.RequireUser(s.WorkoutAPI.HandleCreateWorkout))
		// Update a specific workout
		r.Put("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleUpdateWorkoutByID))
		// Delete a specific workout
		r.Delete("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutAPI.HandleDeleteWorkoutByID))
	})

	return r
}
