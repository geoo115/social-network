package migrations

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
)

// ApplyMigrations applies all migrations in the migrations folder
func ApplyMigrations(db *sql.DB) error {
	migrationsDir := "./pkg/db/migrations/"

	files, err := filepath.Glob(migrationsDir + "*.up.sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Applying migration: %s", file)
		migration, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(migration))
		if err != nil {
			return err
		}
	}

	return nil
}
