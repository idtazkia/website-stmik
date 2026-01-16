package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Prodi represents a study program
type Prodi struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Degree    string    `json:"degree"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProdi creates a new study program
func CreateProdi(ctx context.Context, name, code, degree string) (*Prodi, error) {
	var prodi Prodi
	err := pool.QueryRow(ctx, `
		INSERT INTO prodis (name, code, degree)
		VALUES ($1, $2, $3)
		RETURNING id, name, code, degree, is_active, created_at, updated_at
	`, name, code, degree).Scan(
		&prodi.ID, &prodi.Name, &prodi.Code, &prodi.Degree,
		&prodi.IsActive, &prodi.CreatedAt, &prodi.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create prodi: %w", err)
	}
	return &prodi, nil
}

// FindProdiByID finds a prodi by ID
func FindProdiByID(ctx context.Context, id string) (*Prodi, error) {
	var prodi Prodi
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, degree, is_active, created_at, updated_at
		FROM prodis WHERE id = $1
	`, id).Scan(
		&prodi.ID, &prodi.Name, &prodi.Code, &prodi.Degree,
		&prodi.IsActive, &prodi.CreatedAt, &prodi.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find prodi: %w", err)
	}
	return &prodi, nil
}

// FindProdiByCode finds a prodi by code
func FindProdiByCode(ctx context.Context, code string) (*Prodi, error) {
	var prodi Prodi
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, degree, is_active, created_at, updated_at
		FROM prodis WHERE code = $1
	`, code).Scan(
		&prodi.ID, &prodi.Name, &prodi.Code, &prodi.Degree,
		&prodi.IsActive, &prodi.CreatedAt, &prodi.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find prodi by code: %w", err)
	}
	return &prodi, nil
}

// ListProdis returns all prodis
func ListProdis(ctx context.Context, activeOnly bool) ([]Prodi, error) {
	query := `SELECT id, name, code, degree, is_active, created_at, updated_at FROM prodis`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY code"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list prodis: %w", err)
	}
	defer rows.Close()

	var prodis []Prodi
	for rows.Next() {
		var p Prodi
		err := rows.Scan(&p.ID, &p.Name, &p.Code, &p.Degree, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prodi: %w", err)
		}
		prodis = append(prodis, p)
	}
	return prodis, nil
}

// UpdateProdi updates a prodi
func UpdateProdi(ctx context.Context, id, name, code, degree string) error {
	_, err := pool.Exec(ctx, `
		UPDATE prodis SET name = $1, code = $2, degree = $3 WHERE id = $4
	`, name, code, degree, id)
	if err != nil {
		return fmt.Errorf("failed to update prodi: %w", err)
	}
	return nil
}

// ToggleProdiActive toggles prodi active status
func ToggleProdiActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE prodis SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle prodi active: %w", err)
	}
	return nil
}

// DeleteProdi deletes a prodi (soft delete by deactivating)
func DeleteProdi(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE prodis SET is_active = false WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete prodi: %w", err)
	}
	return nil
}
