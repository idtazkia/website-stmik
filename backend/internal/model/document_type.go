package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// DocumentType represents a type of document required from candidates
type DocumentType struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	Description       *string   `json:"description,omitempty"`
	IsRequired        bool      `json:"is_required"`
	CanDefer          bool      `json:"can_defer"`
	MaxFileSizeMB     int       `json:"max_file_size_mb"`
	AllowedExtensions []string  `json:"allowed_extensions"`
	DisplayOrder      int       `json:"display_order"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateDocumentType creates a new document type
func CreateDocumentType(ctx context.Context, name, code string, description *string, isRequired, canDefer bool, maxFileSizeMB, displayOrder int) (*DocumentType, error) {
	var dt DocumentType
	err := pool.QueryRow(ctx, `
		INSERT INTO document_types (name, code, description, is_required, can_defer, max_file_size_mb, display_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, code, description, is_required, can_defer, max_file_size_mb, allowed_extensions, display_order, is_active, created_at, updated_at
	`, name, code, description, isRequired, canDefer, maxFileSizeMB, displayOrder).Scan(
		&dt.ID, &dt.Name, &dt.Code, &dt.Description, &dt.IsRequired, &dt.CanDefer,
		&dt.MaxFileSizeMB, &dt.AllowedExtensions, &dt.DisplayOrder, &dt.IsActive, &dt.CreatedAt, &dt.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create document type: %w", err)
	}
	return &dt, nil
}

// FindDocumentTypeByID finds a document type by ID
func FindDocumentTypeByID(ctx context.Context, id string) (*DocumentType, error) {
	var dt DocumentType
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, description, is_required, can_defer, max_file_size_mb, allowed_extensions, display_order, is_active, created_at, updated_at
		FROM document_types WHERE id = $1
	`, id).Scan(
		&dt.ID, &dt.Name, &dt.Code, &dt.Description, &dt.IsRequired, &dt.CanDefer,
		&dt.MaxFileSizeMB, &dt.AllowedExtensions, &dt.DisplayOrder, &dt.IsActive, &dt.CreatedAt, &dt.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find document type: %w", err)
	}
	return &dt, nil
}

// FindDocumentTypeByCode finds a document type by code
func FindDocumentTypeByCode(ctx context.Context, code string) (*DocumentType, error) {
	var dt DocumentType
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, description, is_required, can_defer, max_file_size_mb, allowed_extensions, display_order, is_active, created_at, updated_at
		FROM document_types WHERE code = $1
	`, code).Scan(
		&dt.ID, &dt.Name, &dt.Code, &dt.Description, &dt.IsRequired, &dt.CanDefer,
		&dt.MaxFileSizeMB, &dt.AllowedExtensions, &dt.DisplayOrder, &dt.IsActive, &dt.CreatedAt, &dt.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find document type by code: %w", err)
	}
	return &dt, nil
}

// ListDocumentTypes returns all document types
func ListDocumentTypes(ctx context.Context, activeOnly bool) ([]DocumentType, error) {
	query := `SELECT id, name, code, description, is_required, can_defer, max_file_size_mb, allowed_extensions, display_order, is_active, created_at, updated_at FROM document_types`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY display_order, name"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list document types: %w", err)
	}
	defer rows.Close()

	var types []DocumentType
	for rows.Next() {
		var dt DocumentType
		err := rows.Scan(
			&dt.ID, &dt.Name, &dt.Code, &dt.Description, &dt.IsRequired, &dt.CanDefer,
			&dt.MaxFileSizeMB, &dt.AllowedExtensions, &dt.DisplayOrder, &dt.IsActive, &dt.CreatedAt, &dt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document type: %w", err)
		}
		types = append(types, dt)
	}
	return types, nil
}

// UpdateDocumentType updates a document type
func UpdateDocumentType(ctx context.Context, id, name, code string, description *string, isRequired, canDefer bool, maxFileSizeMB, displayOrder int) error {
	_, err := pool.Exec(ctx, `
		UPDATE document_types
		SET name = $1, code = $2, description = $3, is_required = $4, can_defer = $5, max_file_size_mb = $6, display_order = $7, updated_at = NOW()
		WHERE id = $8
	`, name, code, description, isRequired, canDefer, maxFileSizeMB, displayOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update document type: %w", err)
	}
	return nil
}

// ToggleDocumentTypeActive toggles document type active status
func ToggleDocumentTypeActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE document_types SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle document type active: %w", err)
	}
	return nil
}

// DeleteDocumentType deletes a document type
func DeleteDocumentType(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM document_types WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete document type: %w", err)
	}
	return nil
}
