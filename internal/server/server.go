package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/strangecousinwst/goworkout/internal/api"
	"github.com/strangecousinwst/goworkout/internal/database"
	"github.com/strangecousinwst/goworkout/internal/middleware"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/migrations"
)

// Server holds all dependencies for the HTTP server,
// including API handlers, database service, logger, and middleware.
type Server struct {
	port       int
	Logger     *log.Logger
	WorkoutAPI *api.WorkoutAPI
	UserAPI    *api.UserAPI
	TokenAPI   *api.TokenAPI
	Middleware middleware.UserMiddleware
	db         database.Service
}

// NewServer initializes a new HTTP server with configured routes,
// middleware, and all necessary dependencies.
func NewServer() *http.Server {
	// Determine port from environment variable, fallback is 8080.
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	// Initialize a standard logger.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize database service.
	dbService := database.New()
	pgDB := dbService.GetDB()

	// Apply database migrations using the embedded filesystem.
	err = database.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Initialize data stores (data access layer).
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// Initialize API handlers with their respective stores and logger.
	workoutAPI := api.NewWorkoutAPI(workoutStore, logger)
	userAPI := api.NewUserAPI(userStore, logger)
	tokenAPI := api.NewTokenAPI(tokenStore, userStore, logger)

	// Initialize middleware.
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	// Create the server instance with all dependencies.
	server := &Server{
		port:       port,
		Logger:     logger,
		WorkoutAPI: workoutAPI,
		UserAPI:    userAPI,
		TokenAPI:   tokenAPI,
		Middleware: middlewareHandler,
		db:         dbService,
	}

	// Configure the HTTP server.
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port), // Read from .env or fallback to 8080
		Handler:      server.RegisterRoutes(),         // server.RegisterRoutes() is defined in routes.go.
		IdleTimeout:  time.Minute,                     // Max time to wait for idle connections.
		ReadTimeout:  10 * time.Second,                // Max time to read request.
		WriteTimeout: 30 * time.Second,                // Max time to write response.
	}

	return httpServer
}

// Basic health check endpoint for the service.
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
