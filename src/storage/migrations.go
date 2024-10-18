package storage

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(db *sql.DB) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		log.Fatalf("Error creating migration: %v", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations are required - database is up to date.")
		} else {
			log.Fatalf("Error applying migrations\n: %v", err)
		}
		return
	}

	// Get a list of all applied migrations
	appliedMigrations, _, err := m.Version()
	if err != nil {
		log.Fatalf("Error getting migration version: %v", err)
	}

	log.Printf("Migrations completed successfully, current database version: %d", appliedMigrations)
}
