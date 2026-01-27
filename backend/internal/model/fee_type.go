package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// FeeType represents a type of fee (registration, tuition, dormitory)
type FeeType struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Code               string    `json:"code"`
	IsRecurring        bool      `json:"is_recurring"`
	InstallmentOptions []int     `json:"installment_options"`
	CreatedAt          time.Time `json:"created_at"`
}

// FindFeeTypeByID finds a fee type by ID
func FindFeeTypeByID(ctx context.Context, id string) (*FeeType, error) {
	var ft FeeType
	var optionsJSON []byte
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, is_recurring, installment_options, created_at
		FROM fee_types WHERE id = $1
	`, id).Scan(&ft.ID, &ft.Name, &ft.Code, &ft.IsRecurring, &optionsJSON, &ft.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find fee type: %w", err)
	}
	if err := json.Unmarshal(optionsJSON, &ft.InstallmentOptions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installment options: %w", err)
	}
	return &ft, nil
}

// FindFeeTypeByCode finds a fee type by code
func FindFeeTypeByCode(ctx context.Context, code string) (*FeeType, error) {
	var ft FeeType
	var optionsJSON []byte
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, is_recurring, installment_options, created_at
		FROM fee_types WHERE code = $1
	`, code).Scan(&ft.ID, &ft.Name, &ft.Code, &ft.IsRecurring, &optionsJSON, &ft.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find fee type by code: %w", err)
	}
	if err := json.Unmarshal(optionsJSON, &ft.InstallmentOptions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installment options: %w", err)
	}
	return &ft, nil
}

// ListFeeTypes returns all fee types
func ListFeeTypes(ctx context.Context) ([]FeeType, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, name, code, is_recurring, installment_options, created_at
		FROM fee_types ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list fee types: %w", err)
	}
	defer rows.Close()

	var feeTypes []FeeType
	for rows.Next() {
		var ft FeeType
		var optionsJSON []byte
		err := rows.Scan(&ft.ID, &ft.Name, &ft.Code, &ft.IsRecurring, &optionsJSON, &ft.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee type: %w", err)
		}
		if err := json.Unmarshal(optionsJSON, &ft.InstallmentOptions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal installment options: %w", err)
		}
		feeTypes = append(feeTypes, ft)
	}
	return feeTypes, nil
}
