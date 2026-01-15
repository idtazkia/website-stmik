package testutil

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB holds the test database container and connection info
type TestDB struct {
	Container     *postgres.PostgresContainer
	ConnectionStr string
}

// SetupTestDB creates a PostgreSQL container and runs migrations
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Run migrations
	migrationsPath := getMigrationsPath()
	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		container.Terminate(ctx)
		t.Fatalf("failed to run migrations: %v", err)
	}

	return &TestDB{
		Container:     container,
		ConnectionStr: connStr,
	}
}

// Teardown terminates the container
func (tdb *TestDB) Teardown(t *testing.T) {
	t.Helper()
	if err := tdb.Container.Terminate(context.Background()); err != nil {
		t.Errorf("failed to terminate container: %v", err)
	}
}

// getMigrationsPath returns the absolute path to migrations directory
func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "migrations")
}
