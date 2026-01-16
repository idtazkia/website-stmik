package model

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
)

// VerificationToken represents an OTP token for email/phone verification
type VerificationToken struct {
	ID          string     `json:"id"`
	CandidateID string     `json:"candidate_id"`
	TokenType   string     `json:"token_type"` // 'email' or 'phone'
	Token       string     `json:"token"`
	ExpiresAt   time.Time  `json:"expires_at"`
	UsedAt      *time.Time `json:"used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

const (
	TokenTypeEmail = "email"
	TokenTypePhone = "phone"
	TokenLength    = 6
	TokenExpiry    = 15 * time.Minute
)

// generateOTP generates a cryptographically secure 6-digit OTP
func generateOTP() (string, error) {
	const digits = "0123456789"
	otp := make([]byte, TokenLength)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %w", err)
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}

// CreateVerificationToken creates a new OTP token for a candidate
func CreateVerificationToken(ctx context.Context, candidateID, tokenType string) (string, error) {
	// Validate token type
	if tokenType != TokenTypeEmail && tokenType != TokenTypePhone {
		return "", fmt.Errorf("invalid token type: %s", tokenType)
	}

	// Generate OTP
	otp, err := generateOTP()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(TokenExpiry)

	_, err = pool.Exec(ctx, `
		INSERT INTO verification_tokens (candidate_id, token_type, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`, candidateID, tokenType, otp, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create verification token: %w", err)
	}

	return otp, nil
}

// VerifyToken verifies an OTP token and marks it as used
func VerifyToken(ctx context.Context, candidateID, tokenType, token string) error {
	var tokenID string
	err := pool.QueryRow(ctx, `
		SELECT id FROM verification_tokens
		WHERE candidate_id = $1 AND token_type = $2 AND token = $3
		  AND expires_at > NOW() AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, candidateID, tokenType, token).Scan(&tokenID)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("invalid or expired token")
	}
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	// Mark token as used
	_, err = pool.Exec(ctx, `
		UPDATE verification_tokens SET used_at = NOW() WHERE id = $1
	`, tokenID)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// InvalidatePendingTokens invalidates all pending tokens for a candidate and token type
func InvalidatePendingTokens(ctx context.Context, candidateID, tokenType string) error {
	_, err := pool.Exec(ctx, `
		UPDATE verification_tokens
		SET used_at = NOW()
		WHERE candidate_id = $1 AND token_type = $2 AND used_at IS NULL
	`, candidateID, tokenType)
	if err != nil {
		return fmt.Errorf("failed to invalidate pending tokens: %w", err)
	}
	return nil
}

// CleanupExpiredTokens removes tokens that have expired more than 24 hours ago
func CleanupExpiredTokens(ctx context.Context) (int64, error) {
	result, err := pool.Exec(ctx, `
		DELETE FROM verification_tokens
		WHERE expires_at < NOW() - INTERVAL '24 hours'
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return result.RowsAffected(), nil
}

// GetLatestToken gets the latest unused token for a candidate (for debugging/admin)
func GetLatestToken(ctx context.Context, candidateID, tokenType string) (*VerificationToken, error) {
	var token VerificationToken
	err := pool.QueryRow(ctx, `
		SELECT id, candidate_id, token_type, token, expires_at, used_at, created_at
		FROM verification_tokens
		WHERE candidate_id = $1 AND token_type = $2 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, candidateID, tokenType).Scan(
		&token.ID, &token.CandidateID, &token.TokenType, &token.Token,
		&token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest token: %w", err)
	}
	return &token, nil
}
