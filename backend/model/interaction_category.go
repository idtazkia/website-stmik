package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// InteractionCategory represents a category for candidate interactions
type InteractionCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Sentiment    string    `json:"sentiment"` // positive, neutral, negative
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateInteractionCategory creates a new interaction category
func CreateInteractionCategory(ctx context.Context, name, sentiment string, displayOrder int) (*InteractionCategory, error) {
	var cat InteractionCategory
	err := pool.QueryRow(ctx, `
		INSERT INTO interaction_categories (name, sentiment, display_order)
		VALUES ($1, $2, $3)
		RETURNING id, name, sentiment, is_active, display_order, created_at
	`, name, sentiment, displayOrder).Scan(
		&cat.ID, &cat.Name, &cat.Sentiment, &cat.IsActive, &cat.DisplayOrder, &cat.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create interaction category: %w", err)
	}
	return &cat, nil
}

// FindInteractionCategoryByID finds a category by ID
func FindInteractionCategoryByID(ctx context.Context, id string) (*InteractionCategory, error) {
	var cat InteractionCategory
	err := pool.QueryRow(ctx, `
		SELECT id, name, sentiment, is_active, display_order, created_at
		FROM interaction_categories WHERE id = $1
	`, id).Scan(
		&cat.ID, &cat.Name, &cat.Sentiment, &cat.IsActive, &cat.DisplayOrder, &cat.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find interaction category: %w", err)
	}
	return &cat, nil
}

// ListInteractionCategories returns all interaction categories
func ListInteractionCategories(ctx context.Context, activeOnly bool) ([]InteractionCategory, error) {
	query := `SELECT id, name, sentiment, is_active, display_order, created_at FROM interaction_categories`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY display_order, name"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list interaction categories: %w", err)
	}
	defer rows.Close()

	var categories []InteractionCategory
	for rows.Next() {
		var cat InteractionCategory
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Sentiment, &cat.IsActive, &cat.DisplayOrder, &cat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interaction category: %w", err)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// UpdateInteractionCategory updates a category
func UpdateInteractionCategory(ctx context.Context, id, name, sentiment string, displayOrder int) error {
	_, err := pool.Exec(ctx, `
		UPDATE interaction_categories SET name = $1, sentiment = $2, display_order = $3 WHERE id = $4
	`, name, sentiment, displayOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update interaction category: %w", err)
	}
	return nil
}

// ToggleInteractionCategoryActive toggles category active status
func ToggleInteractionCategoryActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE interaction_categories SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle interaction category active: %w", err)
	}
	return nil
}

// DeleteInteractionCategory deletes a category
func DeleteInteractionCategory(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM interaction_categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete interaction category: %w", err)
	}
	return nil
}
