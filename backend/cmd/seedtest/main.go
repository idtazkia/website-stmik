// Command seedtest seeds test users for E2E testing.
// This should only be run in test/development environments.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/idtazkia/stmik-admission-api/internal/config"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/internal/pkg/crypto"
	"github.com/joho/godotenv"
)

// Test users to seed
var testUsers = []struct {
	Email string
	Name  string
	Role  string
}{
	{"test-admin@tazkia.ac.id", "Test Admin User", "admin"},
	{"test-supervisor@tazkia.ac.id", "Test Supervisor User", "supervisor"},
	{"test-consultant@tazkia.ac.id", "Test Consultant User", "consultant"},
	{"test-finance@tazkia.ac.id", "Test Finance User", "finance"},
	{"test-academic@tazkia.ac.id", "Test Academic User", "academic"},
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize encryption
	if err := crypto.Init(cfg.Encryption.Key); err != nil {
		log.Fatalf("failed to initialize encryption: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := model.Connect(ctx, cfg.Database.DSN()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer model.Close()

	log.Println("Seeding test users...")

	var supervisorID string

	for _, tu := range testUsers {
		// Check if user already exists
		existing, err := model.FindUserByEmail(ctx, tu.Email)
		if err != nil {
			log.Printf("Error checking user %s: %v", tu.Email, err)
			os.Exit(1)
		}

		if existing != nil {
			log.Printf("  [skip] %s (%s) - already exists", tu.Email, tu.Role)
			if tu.Role == "supervisor" {
				supervisorID = existing.ID
			}
			continue
		}

		// Create user
		user, err := model.CreateUser(ctx, tu.Email, tu.Name, "", tu.Role)
		if err != nil {
			log.Printf("Error creating user %s: %v", tu.Email, err)
			os.Exit(1)
		}

		log.Printf("  [created] %s (%s)", tu.Email, tu.Role)

		if tu.Role == "supervisor" {
			supervisorID = user.ID
		}
	}

	// Assign supervisor to consultant
	if supervisorID != "" {
		consultant, err := model.FindUserByEmail(ctx, "test-consultant@tazkia.ac.id")
		if err == nil && consultant != nil && consultant.IDSupervisor == nil {
			if err := model.UpdateUserSupervisor(ctx, consultant.ID, &supervisorID); err != nil {
				log.Printf("Warning: failed to assign supervisor to consultant: %v", err)
			} else {
				log.Println("  [updated] Assigned supervisor to consultant")
			}
		}
	}

	fmt.Println("Test users seeded successfully")
}
