package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// FeeStructure represents fee amount per type, prodi, and academic year
type FeeStructure struct {
	ID           string    `json:"id"`
	FeeTypeID    string    `json:"fee_type_id"`
	ProdiID      *string   `json:"prodi_id,omitempty"`
	AcademicYear string    `json:"academic_year"`
	Amount       int64     `json:"amount"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FeeStructureWithDetails includes related fee type and prodi names
type FeeStructureWithDetails struct {
	FeeStructure
	FeeTypeName string  `json:"fee_type_name"`
	FeeTypeCode string  `json:"fee_type_code"`
	ProdiName   *string `json:"prodi_name,omitempty"`
	ProdiCode   *string `json:"prodi_code,omitempty"`
}

// CreateFeeStructure creates a new fee structure
func CreateFeeStructure(ctx context.Context, feeTypeID string, prodiID *string, academicYear string, amount int64) (*FeeStructure, error) {
	var fs FeeStructure
	err := pool.QueryRow(ctx, `
		INSERT INTO fee_structures (fee_type_id, prodi_id, academic_year, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fee_type_id, prodi_id, academic_year, amount, is_active, created_at, updated_at
	`, feeTypeID, prodiID, academicYear, amount).Scan(
		&fs.ID, &fs.FeeTypeID, &fs.ProdiID, &fs.AcademicYear, &fs.Amount,
		&fs.IsActive, &fs.CreatedAt, &fs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fee structure: %w", err)
	}
	return &fs, nil
}

// FindFeeStructureByID finds a fee structure by ID
func FindFeeStructureByID(ctx context.Context, id string) (*FeeStructure, error) {
	var fs FeeStructure
	err := pool.QueryRow(ctx, `
		SELECT id, fee_type_id, prodi_id, academic_year, amount, is_active, created_at, updated_at
		FROM fee_structures WHERE id = $1
	`, id).Scan(
		&fs.ID, &fs.FeeTypeID, &fs.ProdiID, &fs.AcademicYear, &fs.Amount,
		&fs.IsActive, &fs.CreatedAt, &fs.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find fee structure: %w", err)
	}
	return &fs, nil
}

// FindFeeByTypeProdiYear finds fee by type, prodi, and year
func FindFeeByTypeProdiYear(ctx context.Context, feeTypeID string, prodiID *string, academicYear string) (*FeeStructure, error) {
	var fs FeeStructure
	var query string
	var args []interface{}

	if prodiID == nil {
		query = `
			SELECT id, fee_type_id, prodi_id, academic_year, amount, is_active, created_at, updated_at
			FROM fee_structures
			WHERE fee_type_id = $1 AND prodi_id IS NULL AND academic_year = $2 AND is_active = true
		`
		args = []interface{}{feeTypeID, academicYear}
	} else {
		query = `
			SELECT id, fee_type_id, prodi_id, academic_year, amount, is_active, created_at, updated_at
			FROM fee_structures
			WHERE fee_type_id = $1 AND prodi_id = $2 AND academic_year = $3 AND is_active = true
		`
		args = []interface{}{feeTypeID, *prodiID, academicYear}
	}

	err := pool.QueryRow(ctx, query, args...).Scan(
		&fs.ID, &fs.FeeTypeID, &fs.ProdiID, &fs.AcademicYear, &fs.Amount,
		&fs.IsActive, &fs.CreatedAt, &fs.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find fee structure: %w", err)
	}
	return &fs, nil
}

// ListFeeStructures returns all fee structures with details
func ListFeeStructures(ctx context.Context, academicYear string, activeOnly bool) ([]FeeStructureWithDetails, error) {
	query := `
		SELECT fs.id, fs.fee_type_id, fs.prodi_id, fs.academic_year, fs.amount, fs.is_active,
			   fs.created_at, fs.updated_at,
			   ft.name as fee_type_name, ft.code as fee_type_code,
			   p.name as prodi_name, p.code as prodi_code
		FROM fee_structures fs
		JOIN fee_types ft ON fs.fee_type_id = ft.id
		LEFT JOIN prodis p ON fs.prodi_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if academicYear != "" {
		argCount++
		query += fmt.Sprintf(" AND fs.academic_year = $%d", argCount)
		args = append(args, academicYear)
	}

	if activeOnly {
		query += " AND fs.is_active = true"
	}

	query += " ORDER BY ft.name, p.code NULLS FIRST"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list fee structures: %w", err)
	}
	defer rows.Close()

	var structures []FeeStructureWithDetails
	for rows.Next() {
		var fs FeeStructureWithDetails
		err := rows.Scan(
			&fs.ID, &fs.FeeTypeID, &fs.ProdiID, &fs.AcademicYear, &fs.Amount,
			&fs.IsActive, &fs.CreatedAt, &fs.UpdatedAt,
			&fs.FeeTypeName, &fs.FeeTypeCode, &fs.ProdiName, &fs.ProdiCode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee structure: %w", err)
		}
		structures = append(structures, fs)
	}
	return structures, nil
}

// UpdateFeeStructure updates a fee structure amount
func UpdateFeeStructure(ctx context.Context, id string, amount int64) error {
	_, err := pool.Exec(ctx, `UPDATE fee_structures SET amount = $1 WHERE id = $2`, amount, id)
	if err != nil {
		return fmt.Errorf("failed to update fee structure: %w", err)
	}
	return nil
}

// ToggleFeeStructureActive toggles fee structure active status
func ToggleFeeStructureActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE fee_structures SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle fee structure active: %w", err)
	}
	return nil
}

// DeleteFeeStructure deletes a fee structure
func DeleteFeeStructure(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM fee_structures WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete fee structure: %w", err)
	}
	return nil
}
