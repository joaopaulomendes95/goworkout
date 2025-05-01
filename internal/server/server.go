package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/strangecousinwst/goworkout/internal/api"
	"github.com/strangecousinwst/goworkout/internal/database"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/migrations"
)

type Server struct {
	port        int
	UserHandler *api.UserHandler

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	// NewServer := &Server{
	// 	port: port,

	// 	db: database.New(),
	// }
	dbService := database.New()
	pgDB := dbService.GetDB()

	err := database.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// TODO: Implement stores
	userStore := store.NewPostgresUserStore(pgDB)

	// TODO: Implement handlers
	userHandler := api.NewUserHandler(userStore)

	server := &Server{
		port:        port,
		UserHandler: userHandler,
		db:          dbService,
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
