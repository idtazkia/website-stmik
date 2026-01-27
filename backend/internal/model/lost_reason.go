package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// LostReason represents a reason for losing a candidate
type LostReason struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateLostReason creates a new lost reason
func CreateLostReason(ctx context.Context, name string, description *string, displayOrder int) (*LostReason, error) {
	var lr LostReason
	err := pool.QueryRow(ctx, `
		INSERT INTO lost_reasons (name, description, display_order)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, is_active, display_order, created_at
	`, name, description, displayOrder).Scan(
		&lr.ID, &lr.Name, &lr.Description, &lr.IsActive, &lr.DisplayOrder, &lr.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lost reason: %w", err)
	}
	return &lr, nil
}

// FindLostReasonByID finds a lost reason by ID
func FindLostReasonByID(ctx context.Context, id string) (*LostReason, error) {
	var lr LostReason
	err := pool.QueryRow(ctx, `
		SELECT id, name, description, is_active, display_order, created_at
		FROM lost_reasons WHERE id = $1
	`, id).Scan(
		&lr.ID, &lr.Name, &lr.Description, &lr.IsActive, &lr.DisplayOrder, &lr.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find lost reason: %w", err)
	}
	return &lr, nil
}

// ListLostReasons returns all lost reasons
func ListLostReasons(ctx context.Context, activeOnly bool) ([]LostReason, error) {
	query := `SELECT id, name, description, is_active, display_order, created_at FROM lost_reasons`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY display_order, name"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list lost reasons: %w", err)
	}
	defer rows.Close()

	var reasons []LostReason
	for rows.Next() {
		var lr LostReason
		err := rows.Scan(&lr.ID, &lr.Name, &lr.Description, &lr.IsActive, &lr.DisplayOrder, &lr.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lost reason: %w", err)
		}
		reasons = append(reasons, lr)
	}
	return reasons, nil
}

// UpdateLostReason updates a lost reason
func UpdateLostReason(ctx context.Context, id, name string, description *string, displayOrder int) error {
	_, err := pool.Exec(ctx, `
		UPDATE lost_reasons SET name = $1, description = $2, display_order = $3 WHERE id = $4
	`, name, description, displayOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update lost reason: %w", err)
	}
	return nil
}

// ToggleLostReasonActive toggles lost reason active status
func ToggleLostReasonActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE lost_reasons SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle lost reason active: %w", err)
	}
	return nil
}

// DeleteLostReason deletes a lost reason
func DeleteLostReason(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM lost_reasons WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete lost reason: %w", err)
	}
	return nil
}
