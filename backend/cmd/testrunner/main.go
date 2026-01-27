// cmd/testrunner/main.go
// Test runner that orchestrates the entire test suite using testcontainers.
// Similar to mvn clean test - starts fresh database, runs migrations, seeds data, runs all tests.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[testrunner] ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received interrupt, cleaning up...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("Test run failed: %v", err)
	}
}

func run(ctx context.Context) error {
	startTime := time.Now()

	// Step 1: Start PostgreSQL container
	log.Println("Starting PostgreSQL container...")
	pg, err := startPostgres(ctx)
	if err != nil {
		return fmt.Errorf("failed to start postgres: %w", err)
	}
	defer func() {
		log.Println("Stopping PostgreSQL container...")
		if err := pg.container.Terminate(context.Background()); err != nil {
			log.Printf("Warning: failed to terminate container: %v", err)
		}
	}()

	// Step 2: Run migrations
	log.Println("Running migrations...")
	if err := runMigrations(pg.connStr); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Step 3: Seed test data
	log.Println("Seeding test data...")
	if err := seedTestData(ctx, pg.host, pg.port); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	// Step 4: Find available port and start server
	port, err := findAvailablePort()
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}
	log.Printf("Starting server on port %d...", port)

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- startServer(serverCtx, pg.host, pg.port, port)
	}()

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForServer(ctx, serverURL, 30*time.Second); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	log.Println("Server is ready")

	// Step 5: Run Go unit tests
	log.Println("Running Go unit tests...")
	if err := runGoTests(ctx); err != nil {
		return fmt.Errorf("go tests failed: %w", err)
	}

	// Step 6: Run E2E tests
	log.Println("Running E2E tests...")
	if err := runE2ETests(ctx, serverURL); err != nil {
		return fmt.Errorf("e2e tests failed: %w", err)
	}

	// Stop server
	serverCancel()
	select {
	case <-serverDone:
	case <-time.After(5 * time.Second):
		log.Println("Warning: server did not stop gracefully")
	}

	duration := time.Since(startTime)
	log.Printf("All tests passed in %s", duration.Round(time.Second))
	return nil
}

type pgInfo struct {
	container *postgres.PostgresContainer
	connStr   string
	host      string
	port      string
}

func startPostgres(ctx context.Context) (*pgInfo, error) {
	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &pgInfo{
		container: container,
		connStr:   connStr,
		host:      host,
		port:      port.Port(),
	}, nil
}

func runMigrations(connStr string) error {
	migrationsPath := getMigrationsPath()
	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "migrations")
}

func seedTestData(ctx context.Context, dbHost, dbPort string) error {
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/seedtest")
	cmd.Env = append(os.Environ(),
		"DATABASE_HOST="+dbHost,
		"DATABASE_PORT="+dbPort,
		"DATABASE_USER=test",
		"DATABASE_PASSWORD=test",
		"DATABASE_NAME=test_db",
		"DATABASE_SSL_MODE=disable",
		"ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func startServer(ctx context.Context, dbHost string, dbPort string, port int) error {
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/server")
	cmd.Env = append(os.Environ(),
		"DATABASE_HOST="+dbHost,
		"DATABASE_PORT="+dbPort,
		"DATABASE_USER=test",
		"DATABASE_PASSWORD=test",
		"DATABASE_NAME=test_db",
		"DATABASE_SSL_MODE=disable",
		fmt.Sprintf("SERVER_PORT=%d", port),
		"SERVER_HOST=localhost",
		"JWT_SECRET=test-secret-key-for-testing-only",
		"ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"SECURE_COOKIE=false",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func waitForServer(ctx context.Context, url string, timeout time.Duration) error {
	healthURL := url + "/health"
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := http.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("server did not become ready within %s", timeout)
}

func runGoTests(ctx context.Context) error {
	// Run unit tests - they use their own testcontainers
	cmd := exec.CommandContext(ctx, "go", "test", "-v", "-short", "./internal/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runE2ETests(ctx context.Context, serverURL string) error {
	cmd := exec.CommandContext(ctx, "npx", "playwright", "test",
		"--config=playwright.testrunner.config.ts",
	)
	cmd.Env = append(os.Environ(),
		"BASE_URL="+serverURL,
		"CI=true", // Use CI mode for retries
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
