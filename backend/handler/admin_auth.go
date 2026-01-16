package handler

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

const (
	// OAuthStateCookie is the name of the OAuth state cookie
	OAuthStateCookie = "oauth_state"
)

// AdminAuthHandler handles admin authentication routes
type AdminAuthHandler struct {
	googleOAuth *auth.GoogleOAuth
	sessionMgr  *auth.SessionManager
}

// NewAdminAuthHandler creates a new admin auth handler
func NewAdminAuthHandler(googleOAuth *auth.GoogleOAuth, sessionMgr *auth.SessionManager) *AdminAuthHandler {
	return &AdminAuthHandler{
		googleOAuth: googleOAuth,
		sessionMgr:  sessionMgr,
	}
}

// RegisterRoutes registers auth routes to the mux
func (h *AdminAuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin/login", h.handleLogin)
	mux.HandleFunc("GET /admin/auth/google", h.handleGoogleAuth)
	mux.HandleFunc("GET /admin/auth/google/callback", h.handleGoogleCallback)
	mux.HandleFunc("POST /admin/logout", h.handleLogout)
}

// handleLogin shows the login page
func (h *AdminAuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	if claims, err := h.sessionMgr.GetClaimsFromRequest(r); err == nil && claims != nil {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	data := NewPageData("Login")
	admin.Login(data, "").Render(r.Context(), w)
}

// handleGoogleAuth initiates the Google OAuth flow
func (h *AdminAuthHandler) handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		slog.Error("failed to generate OAuth state", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store state in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     OAuthStateCookie,
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Google
	url := h.googleOAuth.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback handles the OAuth callback from Google
func (h *AdminAuthHandler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state
	stateCookie, err := r.Cookie(OAuthStateCookie)
	if err != nil {
		slog.Warn("OAuth state cookie not found", "error", err)
		renderLoginError(w, r, "Sesi login tidak valid. Silakan coba lagi.")
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		slog.Warn("OAuth state mismatch")
		renderLoginError(w, r, "Sesi login tidak valid. Silakan coba lagi.")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   OAuthStateCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Check for error from Google
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		slog.Warn("Google OAuth error", "error", errMsg)
		renderLoginError(w, r, "Login dibatalkan atau terjadi kesalahan.")
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := h.googleOAuth.ExchangeCode(r.Context(), code)
	if err != nil {
		slog.Error("failed to exchange OAuth code", "error", err)
		renderLoginError(w, r, "Gagal melakukan autentikasi dengan Google.")
		return
	}

	// Get user info from Google
	googleUser, err := h.googleOAuth.GetUserInfo(r.Context(), token)
	if err != nil {
		slog.Error("failed to get Google user info", "error", err)
		renderLoginError(w, r, "Gagal mengambil informasi akun Google.")
		return
	}

	// Validate email domain
	if err := h.googleOAuth.ValidateEmailDomain(googleUser.Email); err != nil {
		slog.Warn("invalid email domain", "email", googleUser.Email, "error", err)
		renderLoginError(w, r, "Email harus menggunakan domain institusi.")
		return
	}

	// Find or create user
	user, err := model.FindUserByGoogleID(r.Context(), googleUser.ID)
	if err != nil {
		slog.Error("failed to find user by Google ID", "error", err)
		renderLoginError(w, r, "Terjadi kesalahan sistem.")
		return
	}

	if user == nil {
		// Try to find by email
		user, err = model.FindUserByEmail(r.Context(), googleUser.Email)
		if err != nil {
			slog.Error("failed to find user by email", "error", err)
			renderLoginError(w, r, "Terjadi kesalahan sistem.")
			return
		}
	}

	if user == nil {
		// Create new user with default role (consultant)
		user, err = model.CreateUser(r.Context(), googleUser.Email, googleUser.Name, googleUser.ID, "consultant")
		if err != nil {
			slog.Error("failed to create user", "error", err)
			renderLoginError(w, r, "Gagal membuat akun pengguna.")
			return
		}
		slog.Info("new user created via Google OAuth", "email", user.Email, "role", user.Role)
	}

	// Check if user is active
	if !user.IsActive {
		slog.Warn("inactive user attempted login", "email", user.Email)
		renderLoginError(w, r, "Akun Anda tidak aktif. Hubungi administrator.")
		return
	}

	// Update last login
	if err := model.UpdateLastLogin(r.Context(), user.ID); err != nil {
		slog.Warn("failed to update last login", "error", err)
		// Non-fatal, continue
	}

	// Create session token
	sessionToken, err := h.sessionMgr.CreateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		slog.Error("failed to create session token", "error", err)
		renderLoginError(w, r, "Gagal membuat sesi login.")
		return
	}

	// Set session cookie
	h.sessionMgr.SetCookie(w, sessionToken)

	slog.Info("user logged in", "email", user.Email, "role", user.Role)

	// Redirect to dashboard
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// handleLogout logs the user out
func (h *AdminAuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	h.sessionMgr.ClearCookie(w)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

// renderLoginError renders the login page with an error message
func renderLoginError(w http.ResponseWriter, r *http.Request, errMsg string) {
	data := NewPageData("Login")
	admin.Login(data, errMsg).Render(r.Context(), w)
}

// generateRandomState generates a random state string for OAuth CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
