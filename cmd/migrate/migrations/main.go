package main

import (
	"database/sql"
	"github.com/ViniciusDSLima/golang01/config"
	"github.com/ViniciusDSLima/golang01/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"

	"log"
)

func main() {
	dbconn, err := db.NewPostgresStorage(config.Env)

	if err != nil {
		log.Fatal(err)
	}

	initStorage(dbconn)

	driver, err := postgres.WithInstance(dbconn, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	log.Println("Successfully connected to database")
}
