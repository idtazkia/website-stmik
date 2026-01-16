package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Obstacle represents an obstacle that prevents candidate conversion
type Obstacle struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	SuggestedResponse *string   `json:"suggested_response,omitempty"`
	IsActive          bool      `json:"is_active"`
	DisplayOrder      int       `json:"display_order"`
	CreatedAt         time.Time `json:"created_at"`
}

// CreateObstacle creates a new obstacle
func CreateObstacle(ctx context.Context, name string, suggestedResponse *string, displayOrder int) (*Obstacle, error) {
	var obs Obstacle
	err := pool.QueryRow(ctx, `
		INSERT INTO obstacles (name, suggested_response, display_order)
		VALUES ($1, $2, $3)
		RETURNING id, name, suggested_response, is_active, display_order, created_at
	`, name, suggestedResponse, displayOrder).Scan(
		&obs.ID, &obs.Name, &obs.SuggestedResponse, &obs.IsActive, &obs.DisplayOrder, &obs.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create obstacle: %w", err)
	}
	return &obs, nil
}

// FindObstacleByID finds an obstacle by ID
func FindObstacleByID(ctx context.Context, id string) (*Obstacle, error) {
	var obs Obstacle
	err := pool.QueryRow(ctx, `
		SELECT id, name, suggested_response, is_active, display_order, created_at
		FROM obstacles WHERE id = $1
	`, id).Scan(
		&obs.ID, &obs.Name, &obs.SuggestedResponse, &obs.IsActive, &obs.DisplayOrder, &obs.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find obstacle: %w", err)
	}
	return &obs, nil
}

// ListObstacles returns all obstacles
func ListObstacles(ctx context.Context, activeOnly bool) ([]Obstacle, error) {
	query := `SELECT id, name, suggested_response, is_active, display_order, created_at FROM obstacles`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY display_order, name"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list obstacles: %w", err)
	}
	defer rows.Close()

	var obstacles []Obstacle
	for rows.Next() {
		var obs Obstacle
		err := rows.Scan(&obs.ID, &obs.Name, &obs.SuggestedResponse, &obs.IsActive, &obs.DisplayOrder, &obs.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan obstacle: %w", err)
		}
		obstacles = append(obstacles, obs)
	}
	return obstacles, nil
}

// UpdateObstacle updates an obstacle
func UpdateObstacle(ctx context.Context, id, name string, suggestedResponse *string, displayOrder int) error {
	_, err := pool.Exec(ctx, `
		UPDATE obstacles SET name = $1, suggested_response = $2, display_order = $3 WHERE id = $4
	`, name, suggestedResponse, displayOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update obstacle: %w", err)
	}
	return nil
}

// ToggleObstacleActive toggles obstacle active status
func ToggleObstacleActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE obstacles SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle obstacle active: %w", err)
	}
	return nil
}

// DeleteObstacle deletes an obstacle
func DeleteObstacle(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM obstacles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete obstacle: %w", err)
	}
	return nil
}
