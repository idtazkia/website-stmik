package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// AssignmentAlgorithm represents a candidate assignment algorithm
type AssignmentAlgorithm struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FindAssignmentAlgorithmByID finds an algorithm by ID
func FindAssignmentAlgorithmByID(ctx context.Context, id string) (*AssignmentAlgorithm, error) {
	var alg AssignmentAlgorithm
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM assignment_algorithms WHERE id = $1
	`, id).Scan(
		&alg.ID, &alg.Name, &alg.Code, &alg.Description, &alg.IsActive, &alg.CreatedAt, &alg.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find assignment algorithm: %w", err)
	}
	return &alg, nil
}

// FindActiveAssignmentAlgorithm finds the currently active algorithm
func FindActiveAssignmentAlgorithm(ctx context.Context) (*AssignmentAlgorithm, error) {
	var alg AssignmentAlgorithm
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM assignment_algorithms WHERE is_active = true
	`).Scan(
		&alg.ID, &alg.Name, &alg.Code, &alg.Description, &alg.IsActive, &alg.CreatedAt, &alg.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find active assignment algorithm: %w", err)
	}
	return &alg, nil
}

// ListAssignmentAlgorithms returns all assignment algorithms
func ListAssignmentAlgorithms(ctx context.Context) ([]AssignmentAlgorithm, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM assignment_algorithms
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignment algorithms: %w", err)
	}
	defer rows.Close()

	var algorithms []AssignmentAlgorithm
	for rows.Next() {
		var alg AssignmentAlgorithm
		err := rows.Scan(&alg.ID, &alg.Name, &alg.Code, &alg.Description, &alg.IsActive, &alg.CreatedAt, &alg.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment algorithm: %w", err)
		}
		algorithms = append(algorithms, alg)
	}
	return algorithms, nil
}

// SetAssignmentAlgorithmActive sets one algorithm as active (deactivates all others)
func SetAssignmentAlgorithmActive(ctx context.Context, id string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Deactivate all algorithms
	_, err = tx.Exec(ctx, `UPDATE assignment_algorithms SET is_active = false, updated_at = NOW()`)
	if err != nil {
		return fmt.Errorf("failed to deactivate all algorithms: %w", err)
	}

	// Activate the selected one
	_, err = tx.Exec(ctx, `UPDATE assignment_algorithms SET is_active = true, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to activate algorithm: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// UpdateAssignmentAlgorithm updates an algorithm's description
func UpdateAssignmentAlgorithm(ctx context.Context, id string, description *string) error {
	_, err := pool.Exec(ctx, `
		UPDATE assignment_algorithms SET description = $1, updated_at = NOW() WHERE id = $2
	`, description, id)
	if err != nil {
		return fmt.Errorf("failed to update assignment algorithm: %w", err)
	}
	return nil
}
