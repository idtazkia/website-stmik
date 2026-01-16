package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/idtazkia/stmik-admission-api/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Get direction from command line
	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	if direction != "up" && direction != "down" {
		log.Println("Usage: migrate <up|down>")
		os.Exit(1)
	}

	log.Printf("Running migrations (%s) against %s", direction, cfg.Database.Name)

	m, err := migrate.New(
		"file://migrations",
		cfg.Database.DSN(),
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	if direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
	} else {
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
	}

	log.Println("Migrations completed successfully")
}
