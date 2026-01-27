package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// CookieName is the name of the session cookie
	CookieName = "session"
	// DefaultExpiration is the default token expiration time
	DefaultExpiration = 7 * 24 * time.Hour // 7 days
)

// Claims represents JWT claims for user session (staff or candidate)
type Claims struct {
	// Staff fields
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
	Role   string `json:"role,omitempty"`

	// Candidate fields
	CandidateID string `json:"candidate_id,omitempty"`
	IsCandidate bool   `json:"is_candidate,omitempty"`

	jwt.RegisteredClaims
}

// SessionManager handles JWT session management
type SessionManager struct {
	secret     []byte
	expiration time.Duration
	secure     bool // true for HTTPS
}

// NewSessionManager creates a new session manager
func NewSessionManager(secret string, expiration time.Duration, secure bool) *SessionManager {
	if expiration == 0 {
		expiration = DefaultExpiration
	}
	return &SessionManager{
		secret:     []byte(secret),
		expiration: expiration,
		secure:     secure,
	}
}

// CreateToken generates a new JWT token for a user
func (s *SessionManager) CreateToken(userID, email, name, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// CreateCandidateToken generates a new JWT token for a candidate
func (s *SessionManager) CreateCandidateToken(candidateID, email, name string) (string, error) {
	now := time.Now()
	claims := &Claims{
		CandidateID: candidateID,
		Email:       email,
		Name:        name,
		IsCandidate: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken validates and parses a JWT token
func (s *SessionManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// SetCookie sets the session cookie on the response
func (s *SessionManager) SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode, // Lax allows cookie on OAuth redirects
		MaxAge:   int(s.expiration.Seconds()),
	})
}

// ClearCookie removes the session cookie
func (s *SessionManager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// GetTokenFromRequest extracts the session token from request cookies
func (s *SessionManager) GetTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", fmt.Errorf("session cookie not found: %w", err)
	}
	return cookie.Value, nil
}

// GetClaimsFromRequest extracts and validates claims from request
func (s *SessionManager) GetClaimsFromRequest(r *http.Request) (*Claims, error) {
	token, err := s.GetTokenFromRequest(r)
	if err != nil {
		return nil, err
	}
	return s.ValidateToken(token)
}
