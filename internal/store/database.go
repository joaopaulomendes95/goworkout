package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	dbPassword = os.Getenv("BLUEPRINT_DB_PASSWORD")
	dbUsername = os.Getenv("BLUEPRINT_DB_USERNAME")
	dbPort     = os.Getenv("BLUEPRINT_DB_PORT")
	dbHost     = os.Getenv("BLUEPRINT_DB_HOST")
	dbSchema   = os.Getenv("BLUEPRINT_DB_SCHEMA")
)

func Open() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		database,
		dbSchema,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Connected to the database")

	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
