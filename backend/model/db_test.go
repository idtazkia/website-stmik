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

		// Verify columns (matching actual schema from 001_create_users.up.sql)
		columns := []string{"id", "email", "name", "google_id", "role", "id_supervisor", "is_active", "last_login_at", "created_at", "updated_at"}
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

	t.Run("candidates table exists with correct columns", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'candidates'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check candidates table: %v", err)
		}
		if !exists {
			t.Error("candidates table does not exist")
		}

		// Verify key columns
		columns := []string{"id", "email", "phone", "name", "status", "prodi_id", "assigned_consultant_id", "created_at", "updated_at"}
		for _, col := range columns {
			var colExists bool
			err := pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM information_schema.columns
					WHERE table_name = 'candidates' AND column_name = $1
				)
			`, col).Scan(&colExists)
			if err != nil {
				t.Fatalf("failed to check column %s: %v", col, err)
			}
			if !colExists {
				t.Errorf("column %s does not exist in candidates table", col)
			}
		}
	})

	t.Run("interactions table exists with correct columns", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'interactions'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check interactions table: %v", err)
		}
		if !exists {
			t.Error("interactions table does not exist")
		}

		// Verify columns
		columns := []string{"id", "candidate_id", "consultant_id", "channel", "category_id", "obstacle_id", "remarks", "next_followup_date", "created_at"}
		for _, col := range columns {
			var colExists bool
			err := pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM information_schema.columns
					WHERE table_name = 'interactions' AND column_name = $1
				)
			`, col).Scan(&colExists)
			if err != nil {
				t.Fatalf("failed to check column %s: %v", col, err)
			}
			if !colExists {
				t.Errorf("column %s does not exist in interactions table", col)
			}
		}
	})

	t.Run("can insert and query user", func(t *testing.T) {
		var userID string
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, name, role)
			VALUES ('test@example.com', 'Test User', 'consultant')
			RETURNING id
		`).Scan(&userID)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		var email, name, role string
		err = pool.QueryRow(ctx, `SELECT email, name, role FROM users WHERE id = $1`, userID).Scan(&email, &name, &role)
		if err != nil {
			t.Fatalf("failed to query user: %v", err)
		}
		if email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", email)
		}
		if name != "Test User" {
			t.Errorf("expected name Test User, got %s", name)
		}
		if role != "consultant" {
			t.Errorf("expected role consultant, got %s", role)
		}
	})

	t.Run("foreign key constraint works for candidates", func(t *testing.T) {
		// Insert user first (consultant)
		var userID string
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, name, role)
			VALUES ('consultant@example.com', 'Consultant', 'consultant')
			RETURNING id
		`).Scan(&userID)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		// Insert candidate with valid assigned_consultant_id
		var candidateID string
		err = pool.QueryRow(ctx, `
			INSERT INTO candidates (email, name, password_hash, status, assigned_consultant_id)
			VALUES ('candidate@example.com', 'Test Candidate', '$2a$10$dummyhashfortesting', 'registered', $1)
			RETURNING id
		`, userID).Scan(&candidateID)
		if err != nil {
			t.Fatalf("failed to insert candidate: %v", err)
		}

		// Verify candidate exists
		var count int
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM candidates WHERE id = $1`, candidateID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to count candidates: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 candidate, got %d", count)
		}
	})
}
