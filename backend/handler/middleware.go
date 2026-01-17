package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/layouts"
	"github.com/idtazkia/stmik-admission-api/version"
)

// Context keys for storing user information
type contextKey string

const (
	userClaimsKey contextKey = "user_claims"
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
	}
}

// NewPageDataWithUser creates a PageData with user info from context
func NewPageDataWithUser(ctx context.Context, title string) layouts.PageData {
	data := layouts.PageData{
		Title:   title,
		Version: version.Short(),
	}
	if claims := GetUserClaims(ctx); claims != nil {
		data.UserName = claims.Name
		data.UserRole = claims.Role

		// For consultants, fetch unread suggestions count
		if claims.Role == "consultant" {
			count, err := model.CountUnreadSuggestions(ctx, claims.UserID)
			if err != nil {
				slog.Warn("failed to count unread suggestions", "error", err)
			} else {
				data.UnreadSuggestions = count
			}
		}
	}
	return data
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

// RequireAuth middleware checks if user is authenticated
func RequireAuth(sessionMgr *auth.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := sessionMgr.GetClaimsFromRequest(r)
		if err != nil {
			slog.Debug("authentication failed", "error", err)
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}

		// Store claims in context for handlers to use
		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks if user has required role
func RequireRole(roles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaims(r.Context())
		if claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		slog.Warn("access denied", "user", claims.Email, "required_roles", roles, "user_role", claims.Role)
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

// GetUserClaims retrieves user claims from context
func GetUserClaims(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(userClaimsKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// RequireAdmin is a shorthand for requiring admin role
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRole([]string{"admin"}, next)
}

// RequireSupervisorOrAdmin is a shorthand for requiring supervisor or admin role
func RequireSupervisorOrAdmin(next http.Handler) http.Handler {
	return RequireRole([]string{"admin", "supervisor"}, next)
}

// Context key for candidate claims
const candidateClaimsKey contextKey = "candidate_claims"

// RequireCandidateAuth middleware checks if a candidate is authenticated
func RequireCandidateAuth(sessionMgr *auth.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := sessionMgr.GetClaimsFromRequest(r)
		if err != nil || claims == nil || !claims.IsCandidate {
			slog.Debug("candidate authentication failed", "error", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Store claims in context for handlers to use
		ctx := context.WithValue(r.Context(), candidateClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCandidateClaims retrieves candidate claims from context
func GetCandidateClaims(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(candidateClaimsKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
