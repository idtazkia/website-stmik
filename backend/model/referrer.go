package model

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Referrer represents a person who refers candidates
type Referrer struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	Institution        *string   `json:"institution,omitempty"`
	Phone              *string   `json:"phone,omitempty"`
	Email              *string   `json:"email,omitempty"`
	Code               *string   `json:"code,omitempty"`
	BankName           *string   `json:"bank_name,omitempty"`
	BankAccount        *string   `json:"bank_account,omitempty"`
	AccountHolder      *string   `json:"account_holder,omitempty"`
	CommissionOverride *int64    `json:"commission_override,omitempty"`
	PayoutPreference   string    `json:"payout_preference"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// GenerateReferralCode generates a unique referral code
func GenerateReferralCode(name, referrerType string) string {
	// Create prefix from first 2 letters of type and first letter of name
	typePrefix := strings.ToUpper(referrerType[:2])
	namePrefix := strings.ToUpper(string(name[0]))

	// Generate random suffix
	bytes := make([]byte, 3)
	rand.Read(bytes)
	suffix := strings.ToUpper(hex.EncodeToString(bytes))[:4]

	return fmt.Sprintf("REF-%s%s-%s", typePrefix, namePrefix, suffix)
}

// CreateReferrer creates a new referrer
func CreateReferrer(ctx context.Context, name, referrerType string, institution, phone, email, code, bankName, bankAccount, accountHolder *string, commissionOverride *int64, payoutPreference string) (*Referrer, error) {
	var r Referrer
	err := pool.QueryRow(ctx, `
		INSERT INTO referrers (name, type, institution, phone, email, code, bank_name, bank_account, account_holder, commission_override, payout_preference)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, type, institution, phone, email, code, bank_name, bank_account, account_holder, commission_override, payout_preference, is_active, created_at, updated_at
	`, name, referrerType, institution, phone, email, code, bankName, bankAccount, accountHolder, commissionOverride, payoutPreference).Scan(
		&r.ID, &r.Name, &r.Type, &r.Institution, &r.Phone, &r.Email, &r.Code,
		&r.BankName, &r.BankAccount, &r.AccountHolder, &r.CommissionOverride,
		&r.PayoutPreference, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create referrer: %w", err)
	}
	return &r, nil
}

// FindReferrerByID finds a referrer by ID
func FindReferrerByID(ctx context.Context, id string) (*Referrer, error) {
	var r Referrer
	err := pool.QueryRow(ctx, `
		SELECT id, name, type, institution, phone, email, code, bank_name, bank_account, account_holder,
		       commission_override, payout_preference, is_active, created_at, updated_at
		FROM referrers WHERE id = $1
	`, id).Scan(
		&r.ID, &r.Name, &r.Type, &r.Institution, &r.Phone, &r.Email, &r.Code,
		&r.BankName, &r.BankAccount, &r.AccountHolder, &r.CommissionOverride,
		&r.PayoutPreference, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find referrer: %w", err)
	}
	return &r, nil
}

// FindReferrerByCode finds a referrer by referral code
func FindReferrerByCode(ctx context.Context, code string) (*Referrer, error) {
	var r Referrer
	err := pool.QueryRow(ctx, `
		SELECT id, name, type, institution, phone, email, code, bank_name, bank_account, account_holder,
		       commission_override, payout_preference, is_active, created_at, updated_at
		FROM referrers WHERE code = $1 AND is_active = true
	`, code).Scan(
		&r.ID, &r.Name, &r.Type, &r.Institution, &r.Phone, &r.Email, &r.Code,
		&r.BankName, &r.BankAccount, &r.AccountHolder, &r.CommissionOverride,
		&r.PayoutPreference, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find referrer by code: %w", err)
	}
	return &r, nil
}

// ListReferrers returns all referrers, optionally filtered by type
func ListReferrers(ctx context.Context, referrerType string) ([]Referrer, error) {
	query := `
		SELECT id, name, type, institution, phone, email, code, bank_name, bank_account, account_holder,
		       commission_override, payout_preference, is_active, created_at, updated_at
		FROM referrers
	`
	args := []interface{}{}

	if referrerType != "" {
		query += " WHERE type = $1"
		args = append(args, referrerType)
	}
	query += " ORDER BY name"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list referrers: %w", err)
	}
	defer rows.Close()

	var referrers []Referrer
	for rows.Next() {
		var r Referrer
		err := rows.Scan(
			&r.ID, &r.Name, &r.Type, &r.Institution, &r.Phone, &r.Email, &r.Code,
			&r.BankName, &r.BankAccount, &r.AccountHolder, &r.CommissionOverride,
			&r.PayoutPreference, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referrer: %w", err)
		}
		referrers = append(referrers, r)
	}
	return referrers, nil
}

// SearchReferrersByName searches referrers by name (case-insensitive)
func SearchReferrersByName(ctx context.Context, name string) ([]Referrer, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, name, type, institution, phone, email, code, bank_name, bank_account, account_holder,
		       commission_override, payout_preference, is_active, created_at, updated_at
		FROM referrers
		WHERE name ILIKE $1
		ORDER BY name
	`, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search referrers: %w", err)
	}
	defer rows.Close()

	var referrers []Referrer
	for rows.Next() {
		var r Referrer
		err := rows.Scan(
			&r.ID, &r.Name, &r.Type, &r.Institution, &r.Phone, &r.Email, &r.Code,
			&r.BankName, &r.BankAccount, &r.AccountHolder, &r.CommissionOverride,
			&r.PayoutPreference, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referrer: %w", err)
		}
		referrers = append(referrers, r)
	}
	return referrers, nil
}

// UpdateReferrer updates a referrer
func UpdateReferrer(ctx context.Context, id, name, referrerType string, institution, phone, email, code, bankName, bankAccount, accountHolder *string, commissionOverride *int64, payoutPreference string) error {
	_, err := pool.Exec(ctx, `
		UPDATE referrers
		SET name = $1, type = $2, institution = $3, phone = $4, email = $5, code = $6,
		    bank_name = $7, bank_account = $8, account_holder = $9, commission_override = $10,
		    payout_preference = $11, updated_at = NOW()
		WHERE id = $12
	`, name, referrerType, institution, phone, email, code, bankName, bankAccount, accountHolder, commissionOverride, payoutPreference, id)
	if err != nil {
		return fmt.Errorf("failed to update referrer: %w", err)
	}
	return nil
}

// ToggleReferrerActive toggles the active status
func ToggleReferrerActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE referrers SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle referrer active: %w", err)
	}
	return nil
}

// DeleteReferrer deletes a referrer
func DeleteReferrer(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM referrers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete referrer: %w", err)
	}
	return nil
}

// CountReferrersByType returns counts of referrers by type
func CountReferrersByType(ctx context.Context) (map[string]int, error) {
	rows, err := pool.Query(ctx, `
		SELECT type, COUNT(*) as count
		FROM referrers
		WHERE is_active = true
		GROUP BY type
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to count referrers: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	total := 0
	for rows.Next() {
		var t string
		var c int
		if err := rows.Scan(&t, &c); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[t] = c
		total += c
	}
	counts["total"] = total
	return counts, nil
}
