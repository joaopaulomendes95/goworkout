package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/strangecousinwst/goworkout/cmd/web"
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

	// Static assets
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	// Web routes
	s.registerWebRoutes(r)

	// API routes with version
	s.registerAPIRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(s.Middleware.Authenticate)

		r.Get("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", s.Middleware.RequireUser(s.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", s.Middleware.RequireUser(s.WorkoutHandler.HandleDeleteWorkoutByID))
	})

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Post("/users", s.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", s.TokenHandler.HandleCreateToken)

	r.Get("/login", web.LoginWebHandler)
	r.Post("/login", web.LoginWebHandler)

	r.Get("/web", templ.Handler(web.HelloForm()).ServeHTTP)
	r.Post("/hello", web.HelloWebHandler)

	return r
}

func (s *Server) registerWebRoutes(r chi.Router) {
	// HTML routes
	// r.Get("/login", web.LoginWebHandler)
	// r.Post("/login", web.LoginWebHandler)
	// r.Get("/", web.DashboardHandler) // Main dashboard
	// r.Get("/web", templ.Handler(web.HelloForm()).ServeHTTP)
	// r.Post("/hello", web.HelloWebHandler)

	// // Add workout-related web routes
	// r.Get("/workouts", web.WorkoutListHandler)
	// r.Get("/workouts/{id}", web.WorkoutDetailHandler)
	// r.Get("/workouts/create", web.WorkoutCreateHandler)
	// // etc...
	// Group routs and stuff here
}

func (s *Server) registerAPIRoutes(r chi.Router) {
	// API routes with version prefix
	// r.Route("/api/v1", func(r chi.Router) {
	// 	// Auth
	// 	r.Post("/login", s.APIHandler.HandleLogin)
	// 	r.Post("/register", s.UserHandler.HandleRegisterUser)

	// 	// User management
	// 	r.Route("/users", func(r chi.Router) {
	// 		r.Get("/", s.UserHandler.HandleGetUsers)
	// 		r.Post("/", s.UserHandler.HandleCreateUser)
	// 		r.Get("/{id}", s.UserHandler.HandleGetUser)
	// 		// etc...
	// 	})

	// 	// Workout management
	// 	r.Route("/workouts", func(r chi.Router) {
	// 		r.Get("/", s.WorkoutHandler.HandleGetWorkouts)
	// 		r.Post("/", s.WorkoutHandler.HandleCreateWorkout)
	// 		// etc...
	// 	})
	// })
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
