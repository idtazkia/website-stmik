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

	"github.com/idtazkia/stmik-admission-api/internal/auth"
	"github.com/idtazkia/stmik-admission-api/internal/config"
	"github.com/idtazkia/stmik-admission-api/internal/handler"
	"github.com/idtazkia/stmik-admission-api/internal/integration"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/internal/pkg/crypto"
	"github.com/idtazkia/stmik-admission-api/internal/storage"
	"github.com/idtazkia/stmik-admission-api/web/templates/pages"
	"github.com/idtazkia/stmik-admission-api/internal/version"
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

	// Initialize encryption
	if err := crypto.Init(cfg.Encryption.Key); err != nil {
		log.Fatalf("failed to initialize encryption: %v", err)
	}
	log.Println("Encryption initialized")

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
		cfg.Server.SecureCookie,
	)

	// Initialize integration clients (optional - nil if not configured)
	resendClient := integration.NewResendClient(cfg.Resend.APIKey, cfg.Resend.From)
	whatsappClient := integration.NewWhatsAppClient(cfg.WhatsApp.APIURL, cfg.WhatsApp.APIToken)

	// Initialize local storage for file uploads
	localStorage, err := storage.NewLocalStorage(cfg.Upload.Dir, "/uploads")
	if err != nil {
		log.Fatalf("failed to initialize local storage: %v", err)
	}
	log.Printf("Local storage initialized at %s", cfg.Upload.Dir)

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
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Uploaded files
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Upload.Dir))))

	// Auth routes
	adminAuthHandler := handler.NewAdminAuthHandler(googleOAuth, sessionMgr)
	adminAuthHandler.RegisterRoutes(mux)

	// Admin routes (protected)
	adminHandler := handler.NewAdminHandler(sessionMgr, resendClient)
	adminHandler.RegisterRoutes(mux)

	// Finance routes (protected, finance or admin only)
	financeHandler := handler.NewFinanceHandler(resendClient)
	financeHandler.RegisterRoutes(mux, func(next http.Handler) http.Handler {
		return handler.RequireAuth(sessionMgr, next)
	})

	// Portal routes (candidate portal)
	portalHandler := handler.NewPortalHandler(sessionMgr, localStorage)
	portalHandler.RegisterRoutes(mux)

	// Public routes (registration and login)
	publicHandler := handler.NewPublicHandler(sessionMgr, resendClient, whatsappClient)
	publicHandler.RegisterRoutes(mux)

	// Test routes for Playwright (only in development)
	mux.HandleFunc("GET /test/portal", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Portal")
		pages.TestPortal(data).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /test/admin", func(w http.ResponseWriter, r *http.Request) {
		data := handler.NewPageData("Test Admin")
		data.UserRole = "admin" // Default to admin for test page
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

	// Test login endpoint - for E2E testing only
	// Creates or retrieves a test user with the specified role and logs them in
	mux.HandleFunc("GET /test/login/{role}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role := r.PathValue("role")
		if role != "admin" && role != "supervisor" && role != "consultant" && role != "finance" && role != "academic" {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}

		// Try to find existing user with this role
		users, err := model.ListUsers(ctx, role, true)
		var user model.User
		if err != nil || len(users) == 0 {
			// Create test user for this role
			testEmail := fmt.Sprintf("test-%s@tazkia.ac.id", role)
			testName := fmt.Sprintf("Test %s User", role)

			// Check if test user already exists (might be inactive)
			existingUser, err := model.FindUserByEmail(ctx, testEmail)
			if err != nil {
				log.Printf("Error finding test user: %v", err)
				http.Error(w, "Failed to find test user", http.StatusInternalServerError)
				return
			}

			if existingUser != nil {
				user = *existingUser
				log.Printf("Found existing test user: %s (%s)", user.Email, user.Role)
			} else {
				// Create new test user
				newUser, err := model.CreateUser(ctx, testEmail, testName, "", role)
				if err != nil {
					log.Printf("Error creating test user: %v", err)
					http.Error(w, "Failed to create test user", http.StatusInternalServerError)
					return
				}
				user = *newUser
				log.Printf("Created test user: %s (%s)", user.Email, user.Role)

				// If consultant, assign to supervisor
				if role == "consultant" {
					supervisors, _ := model.ListUsers(ctx, "supervisor", true)
					if len(supervisors) > 0 {
						_ = model.UpdateUserSupervisor(ctx, user.ID, &supervisors[0].ID)
					}
				}
			}
		} else {
			user = users[0].User
		}

		token, err := sessionMgr.CreateToken(
			user.ID,
			user.Email,
			user.Name,
			user.Role,
		)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		sessionMgr.SetCookie(w, token)
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	// Test login endpoint for candidate - for E2E testing only
	// Creates or retrieves a test candidate and logs them in
	mux.HandleFunc("GET /test/login/candidate", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		testEmail := "test-candidate@example.com"

		// Try to find existing test candidate
		candidate, err := model.FindCandidateByEmail(ctx, testEmail)
		if err != nil {
			log.Printf("Error finding test candidate: %v", err)
			http.Error(w, "Failed to find test candidate", http.StatusInternalServerError)
			return
		}

		// Create test candidate if not exists
		if candidate == nil {
			passwordHash, err := model.HashPassword("testpassword123")
			if err != nil {
				log.Printf("Error hashing password: %v", err)
				http.Error(w, "Failed to create test candidate", http.StatusInternalServerError)
				return
			}

			candidate, err = model.CreateCandidate(ctx, testEmail, "08123456789", passwordHash)
			if err != nil {
				log.Printf("Error creating test candidate: %v", err)
				http.Error(w, "Failed to create test candidate", http.StatusInternalServerError)
				return
			}

			// Update personal info
			if err := model.UpdateCandidatePersonalInfo(ctx, candidate.ID, "Test Candidate", "Test Address 123", "Jakarta", "DKI Jakarta"); err != nil {
				log.Printf("Error updating personal info: %v", err)
			}

			// Update status to prospecting
			if err := model.UpdateCandidateStatus(ctx, candidate.ID, "prospecting"); err != nil {
				log.Printf("Error updating status: %v", err)
			}

			log.Printf("Created test candidate: %s", candidate.ID)
		}

		// Create session
		candidateEmail := ""
		if candidate.Email != nil {
			candidateEmail = *candidate.Email
		}
		candidateName := ""
		if candidate.Name != nil {
			candidateName = *candidate.Name
		}

		token, err := sessionMgr.CreateCandidateToken(candidate.ID, candidateEmail, candidateName)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		sessionMgr.SetCookie(w, token)
		http.Redirect(w, r, "/portal", http.StatusFound)
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
