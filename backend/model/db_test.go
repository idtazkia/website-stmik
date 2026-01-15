package model_test

import (
	"context"
	"testing"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/testutil"
)

func TestDatabaseSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)

	ctx := context.Background()
	if err := model.Connect(ctx, testDB.ConnectionStr); err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer model.Close()

	pool := model.Pool()

	t.Run("users table exists with correct columns", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'users'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check users table: %v", err)
		}
		if !exists {
			t.Error("users table does not exist")
		}

		// Verify columns
		columns := []string{"id", "email", "password_hash", "full_name", "phone", "role", "google_id", "email_verified", "created_at", "updated_at"}
		for _, col := range columns {
			var colExists bool
			err := pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM information_schema.columns
					WHERE table_name = 'users' AND column_name = $1
				)
			`, col).Scan(&colExists)
			if err != nil {
				t.Fatalf("failed to check column %s: %v", col, err)
			}
			if !colExists {
				t.Errorf("column %s does not exist in users table", col)
			}
		}
	})

	t.Run("applications table exists with correct columns", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'applications'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check applications table: %v", err)
		}
		if !exists {
			t.Error("applications table does not exist")
		}

		// Verify key columns
		columns := []string{"id", "user_id", "application_number", "program", "academic_year", "status", "created_at", "updated_at"}
		for _, col := range columns {
			var colExists bool
			err := pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM information_schema.columns
					WHERE table_name = 'applications' AND column_name = $1
				)
			`, col).Scan(&colExists)
			if err != nil {
				t.Fatalf("failed to check column %s: %v", col, err)
			}
			if !colExists {
				t.Errorf("column %s does not exist in applications table", col)
			}
		}
	})

	t.Run("sessions table exists with correct columns", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'sessions'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check sessions table: %v", err)
		}
		if !exists {
			t.Error("sessions table does not exist")
		}

		// Verify columns
		columns := []string{"id", "user_id", "refresh_token", "user_agent", "ip_address", "expires_at", "created_at"}
		for _, col := range columns {
			var colExists bool
			err := pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM information_schema.columns
					WHERE table_name = 'sessions' AND column_name = $1
				)
			`, col).Scan(&colExists)
			if err != nil {
				t.Fatalf("failed to check column %s: %v", col, err)
			}
			if !colExists {
				t.Errorf("column %s does not exist in sessions table", col)
			}
		}
	})

	t.Run("can insert and query user", func(t *testing.T) {
		var userID string
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, full_name, role)
			VALUES ('test@example.com', 'Test User', 'registrant')
			RETURNING id
		`).Scan(&userID)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		var email, fullName, role string
		err = pool.QueryRow(ctx, `SELECT email, full_name, role FROM users WHERE id = $1`, userID).Scan(&email, &fullName, &role)
		if err != nil {
			t.Fatalf("failed to query user: %v", err)
		}
		if email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", email)
		}
		if fullName != "Test User" {
			t.Errorf("expected full_name Test User, got %s", fullName)
		}
		if role != "registrant" {
			t.Errorf("expected role registrant, got %s", role)
		}
	})

	t.Run("foreign key constraint works", func(t *testing.T) {
		// Insert user first
		var userID string
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, full_name, role)
			VALUES ('applicant@example.com', 'Applicant', 'registrant')
			RETURNING id
		`).Scan(&userID)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		// Insert application with valid user_id
		var appID string
		err = pool.QueryRow(ctx, `
			INSERT INTO applications (user_id, application_number, program, academic_year, status)
			VALUES ($1, 'APP-2025-001', 'SI', '2025/2026', 'draft')
			RETURNING id
		`, userID).Scan(&appID)
		if err != nil {
			t.Fatalf("failed to insert application: %v", err)
		}

		// Verify application exists
		var count int
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM applications WHERE id = $1`, appID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to count applications: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 application, got %d", count)
		}
	})
}
