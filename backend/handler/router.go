package handler

import (
	"net/http"
)

// Router holds all the handlers and middleware
type Router struct {
	mux *http.ServeMux
}

// New creates a new router with all routes registered
func New() *Router {
	r := &Router{
		mux: http.NewServeMux(),
	}

	r.registerRoutes()

	return r
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) registerRoutes() {
	// Health check
	r.mux.HandleFunc("GET /health", r.handleHealth)

	// API routes
	// r.mux.HandleFunc("POST /api/v1/auth/login", r.handleLogin)
	// r.mux.HandleFunc("POST /api/v1/auth/register", r.handleRegister)
	// r.mux.HandleFunc("GET /api/v1/auth/google", r.handleGoogleAuth)
	// r.mux.HandleFunc("GET /api/v1/auth/google/callback", r.handleGoogleCallback)

	// Portal routes (registrant)
	// r.mux.HandleFunc("GET /portal", r.requireAuth(r.handlePortalDashboard))
	// r.mux.HandleFunc("GET /portal/application", r.requireAuth(r.handleApplication))
	// r.mux.HandleFunc("POST /portal/application", r.requireAuth(r.handleApplicationSubmit))
	// r.mux.HandleFunc("GET /portal/documents", r.requireAuth(r.handleDocuments))
	// r.mux.HandleFunc("POST /portal/documents/upload", r.requireAuth(r.handleDocumentUpload))

	// Admin routes (staff)
	// r.mux.HandleFunc("GET /admin", r.requireRole("staff")(r.handleAdminDashboard))
	// r.mux.HandleFunc("GET /admin/prospects", r.requireRole("staff")(r.handleProspects))
	// r.mux.HandleFunc("GET /admin/applications", r.requireRole("staff")(r.handleApplications))
	// r.mux.HandleFunc("GET /admin/settings", r.requireRole("admin")(r.handleSettings))
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Middleware placeholders
// func (r *Router) requireAuth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		// TODO: Implement JWT validation
// 		next(w, req)
// 	}
// }

// func (r *Router) requireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
// 	return func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, req *http.Request) {
// 			// TODO: Implement role checking
// 			next(w, req)
// 		}
// 	}
// }
