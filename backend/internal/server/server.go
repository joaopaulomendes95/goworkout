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

type Server struct {
	port       int
	Logger     *log.Logger
	WorkoutAPI *api.WorkoutAPI
	UserAPI    *api.UserAPI
	TokenAPI   *api.TokenAPI
	Middleware middleware.UserMiddleware
	db         database.Service
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
	workoutAPI := api.NewWorkoutAPI(workoutStore, logger)
	userAPI := api.NewUserAPI(userStore, logger)
	tokenAPI := api.NewTokenAPI(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	server := &Server{
		port:       port,
		Logger:     logger,
		WorkoutAPI: workoutAPI,
		UserAPI:    userAPI,
		TokenAPI:   tokenAPI,
		Middleware: middlewareHandler,
		db:         dbService,
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
