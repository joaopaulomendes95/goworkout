package server

import (
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

type Server struct {
	port           int
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	db             database.Service
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// NewServer := &Server{
	// 	port: port,

	// 	db: database.New(),
	// }
	dbService := database.New()
	pgDB := dbService.GetDB()

	err = database.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// TODO: Implement stores
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// TODO: Implement handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore)
	userHandler := api.NewUserHandler(userStore)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	server := &Server{
		port:           port,
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		db:             dbService,
	}

	// Declare Server config
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
