// Package main provides a mockup server without database connection
// This is for UI validation with stakeholders before real implementation
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
	"github.com/idtazkia/stmik-admission-api/handler"
	"github.com/idtazkia/stmik-admission-api/storage"
	"github.com/idtazkia/stmik-admission-api/templates/pages"
	"github.com/idtazkia/stmik-admission-api/version"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Create a mock session manager for mockup mode
	// Uses a fixed secret - this is fine since mockup doesn't need real security
	mockSessionMgr := auth.NewSessionManager(
		"mockup-secret-key-for-development",
		24*time.Hour,
		false, // secure=false for development
	)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"status":  "ok",
			"mode":    "mockup",
			"version": version.Info(),
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Root redirect to admin
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	// Mock login endpoint - sets a mock session cookie
	mux.HandleFunc("GET /mock-login", func(w http.ResponseWriter, r *http.Request) {
		// Create mock session for admin user
		token, err := mockSessionMgr.CreateToken("mock-admin", "admin@mockup.test", "Admin Mockup", "admin")
		if err != nil {
			http.Error(w, "Failed to create mock session", http.StatusInternalServerError)
			return
		}
		mockSessionMgr.SetCookie(w, token)
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Admin routes (mockup) - uses session manager for auth middleware
	adminHandler := handler.NewAdminHandler(mockSessionMgr, nil)
	adminHandler.RegisterRoutes(mux)

	// Initialize local storage for file uploads (mockup mode uses ./uploads)
	mockStorage, err := storage.NewLocalStorage("./uploads", "/uploads")
	if err != nil {
		log.Fatalf("failed to initialize local storage: %v", err)
	}

	// Portal routes (mockup)
	portalHandler := handler.NewPortalHandler(mockSessionMgr, mockStorage)
	portalHandler.RegisterRoutes(mux)

	// Test routes for Playwright
	mux.HandleFunc("GET /test/portal", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Portal")
		pages.TestPortal(data).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /test/admin", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Admin")
		pages.TestAdmin(data).Render(r.Context(), w)
	})

	// Apply CSRF protection middleware
	protectedMux := handler.CrossOriginProtection(mux)

	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{
		Addr:         addr,
		Handler:      protectedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting mockup server on http://%s", addr)
		log.Println("Available routes:")
		log.Println("  GET  /mock-login           - Auto-login as admin (use this first)")
		log.Println("  GET  /admin                - Dashboard")
		log.Println("  GET  /admin/candidates     - Candidates list")
		log.Println("  GET  /admin/candidates/:id - Candidate detail")
		log.Println("  GET  /admin/campaigns      - Campaigns (coming soon)")
		log.Println("  GET  /admin/referrers      - Referrers (coming soon)")
		log.Println("  GET  /health               - Health check")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down mockup server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Mockup server stopped")
}
