package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" 
	"Social/pkg/db/migrations"
)

var DB *sql.DB

// Initialize the SQLite connection and apply migrations
func Initialize() error {
	var err error

	// Open a connection to the SQLite database
	DB, err = sql.Open("sqlite3", "./your_database_name.db")
	if err != nil {
		return err
	}

	// Apply migrations
	err = migrations.ApplyMigrations(DB)
	if err != nil {
		return err
	}

	return nil
}
