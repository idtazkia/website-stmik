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
	"strings"
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

	// Step 4: Build coverage-instrumented server binary
	log.Println("Building coverage-instrumented server...")
	if err := buildCoverServer(ctx); err != nil {
		return fmt.Errorf("failed to build cover server: %w", err)
	}

	// Step 5: Find available port and start server
	port, err := findAvailablePort()
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}
	log.Printf("Starting server on port %d...", port)

	// Prepare coverage directory for E2E coverage data
	coverDir := "coverage-e2e"
	os.RemoveAll(coverDir)
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		return fmt.Errorf("failed to create coverage dir: %w", err)
	}

	server, err := startServerProcess(pg.host, pg.port, port, coverDir)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForServer(ctx, serverURL, 30*time.Second); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	log.Println("Server is ready")

	// Step 6: Run Go unit tests
	log.Println("Running Go unit tests...")
	if err := runGoTests(ctx); err != nil {
		return fmt.Errorf("go tests failed: %w", err)
	}

	// Step 7: Run E2E tests
	log.Println("Running E2E tests...")
	if err := runE2ETests(ctx, serverURL); err != nil {
		return fmt.Errorf("e2e tests failed: %w", err)
	}

	// Stop server gracefully via SIGTERM so it flushes coverage data
	log.Println("Stopping server...")
	if err := server.stop(); err != nil {
		log.Printf("Warning: server stop: %v", err)
	}

	// Step 8: Generate combined coverage report
	log.Println("Generating coverage report...")
	if err := generateCoverageReport(ctx, coverDir); err != nil {
		log.Printf("Warning: failed to generate coverage report: %v", err)
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
		"JWT_SECRET=test-secret-key-for-testing-only",
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

func buildCoverServer(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "go", "build", "-cover", "-o", "bin/server-cover", "./cmd/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// serverProcess holds the running server process so we can send SIGTERM for graceful shutdown.
// exec.CommandContext sends SIGKILL on cancel, which prevents coverage data from being flushed.
type serverProcess struct {
	cmd *exec.Cmd
}

func startServerProcess(dbHost string, dbPort string, port int, coverDir string) (*serverProcess, error) {
	cmd := exec.Command("./bin/server-cover")
	cmd.Env = append(os.Environ(),
		"GOCOVERDIR="+coverDir,
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
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &serverProcess{cmd: cmd}, nil
}

func (s *serverProcess) stop() error {
	// Send SIGTERM for graceful shutdown so coverage data gets flushed
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %w", err)
	}
	return s.cmd.Wait()
}

func generateCoverageReport(ctx context.Context, coverDir string) error {
	// Convert E2E coverage data to textfmt
	e2eCovFile := "coverage-e2e.out"
	cmd := exec.CommandContext(ctx, "go", "tool", "covdata", "textfmt", "-i="+coverDir, "-o="+e2eCovFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert E2E coverage: %w", err)
	}

	// Read unit test and E2E coverage files, merge into combined report
	unitCov, unitErr := os.ReadFile("coverage.out")
	e2eCov, e2eErr := os.ReadFile(e2eCovFile)

	combined, err := os.Create("coverage-combined.out")
	if err != nil {
		return fmt.Errorf("failed to create combined coverage file: %w", err)
	}
	defer combined.Close()

	// Write the mode line once, then all coverage lines from both files
	combined.WriteString("mode: set\n")
	if unitErr == nil {
		writeProfileLines(combined, unitCov)
	}
	if e2eErr == nil {
		writeProfileLines(combined, e2eCov)
	}

	// Print coverage summary
	log.Println("Coverage report: coverage-combined.out")
	summaryCmd := exec.CommandContext(ctx, "go", "tool", "cover", "-func=coverage-combined.out")
	summaryCmd.Stdout = os.Stdout
	summaryCmd.Stderr = os.Stderr
	summaryCmd.Run()

	return nil
}

// writeProfileLines writes all non-mode lines from a coverage profile to the writer.
func writeProfileLines(w *os.File, data []byte) {
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		w.WriteString(line + "\n")
	}
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
	cmd := exec.CommandContext(ctx, "go", "test", "-v", "-short", "-coverprofile=coverage.out", "./internal/...")
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
