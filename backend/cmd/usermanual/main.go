// Command usermanual orchestrates the complete user manual generation:
// 1. Start PostgreSQL via testcontainers
// 2. Run migrations
// 3. Seed comprehensive test data (cmd/seedmanual)
// 4. Start HTTP server
// 5. Capture screenshots via Playwright
// 6. Generate HTML manual from markdown + screenshots
// 7. Shutdown everything
//
// Usage:
//
//	go run ./cmd/usermanual              # full pipeline
//	go run ./cmd/usermanual --generate   # generate HTML only (skip infra)
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

const (
	testDBUser     = "test"
	testDBPassword = "test"
	testDBName     = "manual_db"
	testJWTSecret  = "test-secret-key-for-manual-only"
	testEncKey     = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[usermanual] ")

	// --generate flag: skip infra, just generate HTML from existing markdown + screenshots
	if len(os.Args) > 1 && os.Args[1] == "--generate" {
		if err := GenerateManual("docs/user-manual", "docs/user-manual/output"); err != nil {
			log.Fatalf("Generation failed: %v", err)
		}
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received interrupt, cleaning up...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("Manual generation failed: %v", err)
	}
}

func run(ctx context.Context) error {
	startTime := time.Now()

	// Step 1: Start PostgreSQL container
	log.Println("Step 1/6: Starting PostgreSQL container...")
	pg, err := startPostgresContainer(ctx)
	if err != nil {
		return fmt.Errorf("failed to start postgres: %w", err)
	}
	defer func() {
		log.Println("Stopping PostgreSQL container...")
		pg.container.Terminate(context.Background())
	}()

	// Step 2: Run migrations
	log.Println("Step 2/6: Running migrations...")
	if err := runDBMigrations(pg.connStr); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Step 3: Seed comprehensive test data
	log.Println("Step 3/6: Seeding comprehensive test data...")
	if err := seedManualData(ctx, pg.host, pg.port); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	// Step 4: Start HTTP server
	port, err := findFreePort()
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}
	log.Printf("Step 4/6: Starting server on port %d...", port)

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- startHTTPServer(serverCtx, pg.host, pg.port, port)
	}()

	serverURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForHealthy(ctx, serverURL, 30*time.Second); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	log.Println("Server is ready")

	// Step 5: Capture screenshots via Playwright
	log.Println("Step 5/6: Capturing screenshots...")
	if err := captureScreenshots(ctx, serverURL); err != nil {
		return fmt.Errorf("screenshot capture failed: %w", err)
	}

	// Step 6: Generate HTML manual
	log.Println("Step 6/6: Generating HTML manual...")
	if err := GenerateManual("docs/user-manual", "docs/user-manual/output"); err != nil {
		return fmt.Errorf("manual generation failed: %w", err)
	}

	// Shutdown server
	serverCancel()
	select {
	case <-serverDone:
	case <-time.After(5 * time.Second):
		log.Println("Warning: server did not stop gracefully")
	}

	duration := time.Since(startTime)
	log.Printf("User manual generated in %s", duration.Round(time.Second))
	log.Println("Output: docs/user-manual/output/index.html")
	return nil
}

type pgInstance struct {
	container *postgres.PostgresContainer
	connStr   string
	host      string
	port      string
}

func startPostgresContainer(ctx context.Context) (*pgInstance, error) {
	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPassword),
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

	return &pgInstance{
		container: container,
		connStr:   connStr,
		host:      host,
		port:      port.Port(),
	}, nil
}

func runDBMigrations(connStr string) error {
	migrationsPath := getLocalMigrationsPath()
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

func getLocalMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "migrations")
}

func dbEnv(dbHost, dbPort string) []string {
	return []string{
		"DATABASE_HOST=" + dbHost,
		"DATABASE_PORT=" + dbPort,
		"DATABASE_USER=" + testDBUser,
		"DATABASE_PASSWORD=" + testDBPassword,
		"DATABASE_NAME=" + testDBName,
		"DATABASE_SSL_MODE=disable",
		"JWT_SECRET=" + testJWTSecret,
		"ENCRYPTION_KEY=" + testEncKey,
	}
}

func seedManualData(ctx context.Context, dbHost, dbPort string) error {
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/seedmanual")
	cmd.Env = append(os.Environ(), dbEnv(dbHost, dbPort)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func startHTTPServer(ctx context.Context, dbHost, dbPort string, port int) error {
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/server")
	env := append(dbEnv(dbHost, dbPort),
		fmt.Sprintf("SERVER_PORT=%d", port),
		"SERVER_HOST=localhost",
		"SECURE_COOKIE=false",
	)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func waitForHealthy(ctx context.Context, url string, timeout time.Duration) error {
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

func captureScreenshots(ctx context.Context, serverURL string) error {
	cmd := exec.CommandContext(ctx, "npx", "playwright", "test",
		"--config=playwright.testrunner.config.ts",
		"screenshot-capture",
	)
	cmd.Env = append(os.Environ(),
		"BASE_URL="+serverURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
