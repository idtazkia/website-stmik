package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/config"
	"github.com/idtazkia/stmik-admission-api/handler"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/pages"
	"github.com/idtazkia/stmik-admission-api/version"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (development only)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database connection
	ctx := context.Background()
	if err := model.Connect(ctx, cfg.Database.DSN()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer model.Close()
	log.Println("Connected to database")

	// Initialize auth components
	googleOAuth := auth.NewGoogleOAuth(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
		cfg.Google.StaffEmailDomain,
	)
	sessionMgr := auth.NewSessionManager(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.ExpirationHours)*time.Hour,
		false, // secure=false for development
	)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"status":  "ok",
			"version": version.Info(),
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Auth routes
	adminAuthHandler := handler.NewAdminAuthHandler(googleOAuth, sessionMgr)
	adminAuthHandler.RegisterRoutes(mux)

	// Admin routes (protected)
	adminHandler := handler.NewAdminHandler(sessionMgr)
	adminHandler.RegisterRoutes(mux)

	// Portal routes
	portalHandler := handler.NewPortalHandler()
	portalHandler.RegisterRoutes(mux)

	// Test routes for Playwright (only in development)
	mux.HandleFunc("GET /test/portal", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Portal")
		pages.TestPortal(data).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /test/admin", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Admin")
		pages.TestAdmin(data).Render(r.Context(), w)
	})

	mux.HandleFunc("POST /test/submit", func(w http.ResponseWriter, r *http.Request) {
		// Test form submission (CSRF protected)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Form submitted successfully",
		})
	})

	// Apply CSRF protection middleware
	protectedMux := handler.CrossOriginProtection(mux)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      protectedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
