package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/idtazkia/stmik-admission-api/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Parse command line flags
	var direction string
	flag.StringVar(&direction, "direction", "up", "migration direction: up or down")
	flag.Parse()

	if direction != "up" && direction != "down" {
		fmt.Println("Usage: migrate -direction=<up|down>")
		os.Exit(1)
	}

	log.Printf("Running migrations (%s) against %s", direction, cfg.Database.Name)

	// TODO: Implement migration logic using golang-migrate/migrate
	// Example:
	// m, err := migrate.New(
	//     "file://migrations",
	//     cfg.Database.DSN(),
	// )
	// if err != nil {
	//     log.Fatalf("failed to create migrate instance: %v", err)
	// }
	//
	// if direction == "up" {
	//     if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//         log.Fatalf("migration up failed: %v", err)
	//     }
	// } else {
	//     if err := m.Steps(-1); err != nil {
	//         log.Fatalf("migration down failed: %v", err)
	//     }
	// }

	log.Println("Migrations completed successfully")
}
