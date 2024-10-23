package db

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

// Initialize the SQLite connection and apply migrations
func Initialize() error {
	// Correct absolute migration directory path
	migrationsDir := "../pkg/db/migrations"

	// Open a connection to the SQLite database
	var err error
	DB, err = sql.Open("sqlite3", "./socialNetwork1.db")
	if err != nil {
		return err
	}

	// Create a new migration instance with the absolute path
	m, err := migrate.New(
		"file://"+filepath.ToSlash(migrationsDir),
		"sqlite3://./socialNetwork1.db",
	)
	if err != nil {
		return err
	}

	// Apply all available migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
