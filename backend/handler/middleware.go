package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/templates/layouts"
	"github.com/idtazkia/stmik-admission-api/version"
)

// CrossOriginProtection returns CSRF protection middleware using Go 1.25's stdlib
// Uses Sec-Fetch-Site and Origin headers (92-95% browser support)
// Combined with SameSite=Strict cookies for defense in depth
func CrossOriginProtection(next http.Handler) http.Handler {
	cop := &http.CrossOriginProtection{}
	return cop.Handler(next)
}

// NewPageData creates a PageData with common fields populated
// Note: Go 1.25's CSRF protection doesn't require tokens in forms
// It uses browser headers instead, so CSRFToken/CSRFField are empty
func NewPageData(title string) layouts.PageData {
	return layouts.PageData{
		Title:   title,
		Version: version.Short(),
		// CSRFToken and CSRFField not needed with Go 1.25's CrossOriginProtection
		// It uses Sec-Fetch-Site and Origin headers instead
	}
}

// Logging middleware logs HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

// Recovery middleware recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
