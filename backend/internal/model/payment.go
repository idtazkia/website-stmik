package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Payment statuses
const (
	PaymentStatusPending  = "pending"
	PaymentStatusApproved = "approved"
	PaymentStatusRejected = "rejected"
)

// Payment represents a payment proof record
type Payment struct {
	ID              string     `json:"id"`
	BillingID       string     `json:"id_billing"`
	Amount          int        `json:"amount"`
	TransferDate    time.Time  `json:"transfer_date"`
	ProofFilePath   string     `json:"proof_file_path"`
	ProofFileName   string     `json:"proof_file_name"`
	ProofFileSize   int        `json:"proof_file_size"`
	ProofMimeType   string     `json:"proof_mime_type"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	ReviewedBy      *string    `json:"id_reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PaymentWithBilling includes billing info
type PaymentWithBilling struct {
	Payment
	BillingType     string  `json:"billing_type"`
	BillingAmount   int     `json:"billing_amount"`
	CandidateID     string  `json:"candidate_id"`
	CandidateName   string  `json:"candidate_name"`
	CandidateProdi  *string `json:"candidate_prodi,omitempty"`
}

// CreatePayment creates a new payment proof record
func CreatePayment(ctx context.Context, billingID string, amount int, transferDate time.Time, proofFilePath, proofFileName string, proofFileSize int, proofMimeType string) (*Payment, error) {
	var p Payment
	err := pool.QueryRow(ctx, `
		INSERT INTO payments (id_billing, amount, transfer_date, proof_file_path, proof_file_name, proof_file_size, proof_mime_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, id_billing, amount, transfer_date, proof_file_path, proof_file_name, proof_file_size,
		          proof_mime_type, status, rejection_reason, id_reviewed_by, reviewed_at, created_at, updated_at
	`, billingID, amount, transferDate, proofFilePath, proofFileName, proofFileSize, proofMimeType).Scan(
		&p.ID, &p.BillingID, &p.Amount, &p.TransferDate, &p.ProofFilePath, &p.ProofFileName, &p.ProofFileSize,
		&p.ProofMimeType, &p.Status, &p.RejectionReason, &p.ReviewedBy, &p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update billing status to pending_verification
	_, err = pool.Exec(ctx, `
		UPDATE billings SET status = 'pending_verification', updated_at = NOW() WHERE id = $1
	`, billingID)
	if err != nil {
		return nil, fmt.Errorf("failed to update billing status: %w", err)
	}

	return &p, nil
}

// FindPaymentByID finds a payment by ID
func FindPaymentByID(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	err := pool.QueryRow(ctx, `
		SELECT id, id_billing, amount, transfer_date, proof_file_path, proof_file_name, proof_file_size,
		       proof_mime_type, status, rejection_reason, id_reviewed_by, reviewed_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.BillingID, &p.Amount, &p.TransferDate, &p.ProofFilePath, &p.ProofFileName, &p.ProofFileSize,
		&p.ProofMimeType, &p.Status, &p.RejectionReason, &p.ReviewedBy, &p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}
	return &p, nil
}

// PaymentReviewFilters for listing payments to review
type PaymentReviewFilters struct {
	Status string
	Search string
}

// ListPaymentsForReview returns payments for admin review
func ListPaymentsForReview(ctx context.Context, filters PaymentReviewFilters) ([]PaymentWithBilling, error) {
	query := `
		SELECT p.id, p.id_billing, p.amount, p.transfer_date, p.proof_file_path, p.proof_file_name,
		       p.proof_file_size, p.proof_mime_type, p.status, p.rejection_reason, p.id_reviewed_by,
		       p.reviewed_at, p.created_at, p.updated_at,
		       b.billing_type, b.amount as billing_amount, c.id as candidate_id, c.name as candidate_name,
		       pr.name as candidate_prodi
		FROM payments p
		JOIN billings b ON b.id = p.id_billing
		JOIN candidates c ON c.id = b.id_candidate
		LEFT JOIN prodis pr ON pr.id = c.id_prodi
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND p.status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Search != "" {
		// With encryption, search by encrypted email (deterministic) for exact match
		emailEnc, err := encryptEmail(filters.Search)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt search term: %w", err)
		}
		query += fmt.Sprintf(" AND c.email = $%d", argIndex)
		args = append(args, emailEnc)
		argIndex++
	}

	query += " ORDER BY p.created_at DESC LIMIT 100"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []PaymentWithBilling
	for rows.Next() {
		var p PaymentWithBilling
		err := rows.Scan(
			&p.ID, &p.BillingID, &p.Amount, &p.TransferDate, &p.ProofFilePath, &p.ProofFileName,
			&p.ProofFileSize, &p.ProofMimeType, &p.Status, &p.RejectionReason, &p.ReviewedBy,
			&p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
			&p.BillingType, &p.BillingAmount, &p.CandidateID, &p.CandidateName, &p.CandidateProdi,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		// Decrypt candidate name
		decName, _ := decryptName(p.CandidateName)
		p.CandidateName = decName

		payments = append(payments, p)
	}
	return payments, nil
}

// PaymentReviewStats for dashboard
type PaymentReviewStats struct {
	Pending       int
	ApprovedToday int
	RejectedToday int
	Total         int
}

// GetPaymentReviewStats returns stats for admin dashboard
func GetPaymentReviewStats(ctx context.Context) (*PaymentReviewStats, error) {
	var stats PaymentReviewStats
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'approved' AND reviewed_at >= CURRENT_DATE) as approved_today,
			COUNT(*) FILTER (WHERE status = 'rejected' AND reviewed_at >= CURRENT_DATE) as rejected_today,
			COUNT(*) as total
		FROM payments
	`).Scan(&stats.Pending, &stats.ApprovedToday, &stats.RejectedToday, &stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment stats: %w", err)
	}
	return &stats, nil
}

// GetPaymentStats returns total count by status for finance
func GetPaymentStats(ctx context.Context) (pending, approved, rejected int, err error) {
	err = pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END), 0)
		FROM payments
	`).Scan(&pending, &approved, &rejected)
	if err != nil {
		err = fmt.Errorf("failed to get payment stats: %w", err)
	}
	return
}

// ListPaymentsWithDetails returns payments with candidate and billing info, with pagination
func ListPaymentsWithDetails(ctx context.Context, status string, page, pageSize int) ([]PaymentWithBilling, int, error) {
	// Build query
	baseQuery := `
		FROM payments p
		JOIN billings b ON b.id = p.id_billing
		JOIN candidates c ON c.id = b.id_candidate
		LEFT JOIN prodis pr ON pr.id = c.id_prodi
		WHERE 1=1
	`
	args := []any{}
	argCount := 0

	if status != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.status = $%d", argCount)
		args = append(args, status)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Pagination
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	// Get payments
	selectQuery := `
		SELECT p.id, p.id_billing, p.amount, p.transfer_date, p.proof_file_path, p.proof_file_name,
		       p.proof_file_size, p.proof_mime_type, p.status, p.rejection_reason, p.id_reviewed_by,
		       p.reviewed_at, p.created_at, p.updated_at,
		       b.billing_type, b.amount as billing_amount, c.id as candidate_id, c.name as candidate_name,
		       pr.name as candidate_prodi
	` + baseQuery + fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pageSize, offset)

	rows, err := pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []PaymentWithBilling
	for rows.Next() {
		var p PaymentWithBilling
		err := rows.Scan(
			&p.ID, &p.BillingID, &p.Amount, &p.TransferDate, &p.ProofFilePath, &p.ProofFileName,
			&p.ProofFileSize, &p.ProofMimeType, &p.Status, &p.RejectionReason, &p.ReviewedBy,
			&p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
			&p.BillingType, &p.BillingAmount, &p.CandidateID, &p.CandidateName, &p.CandidateProdi,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment: %w", err)
		}

		// Decrypt candidate name
		decName, _ := decryptName(p.CandidateName)
		p.CandidateName = decName

		payments = append(payments, p)
	}
	return payments, total, nil
}

// ListPaymentsByBilling returns all payments for a billing
func ListPaymentsByBilling(ctx context.Context, billingID string) ([]Payment, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, id_billing, amount, transfer_date, proof_file_path, proof_file_name, proof_file_size,
		       proof_mime_type, status, rejection_reason, id_reviewed_by, reviewed_at, created_at, updated_at
		FROM payments
		WHERE id_billing = $1
		ORDER BY created_at DESC
	`, billingID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(
			&p.ID, &p.BillingID, &p.Amount, &p.TransferDate, &p.ProofFilePath, &p.ProofFileName, &p.ProofFileSize,
			&p.ProofMimeType, &p.Status, &p.RejectionReason, &p.ReviewedBy, &p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// ApprovePayment approves a payment and updates billing status to paid
func ApprovePayment(ctx context.Context, paymentID, reviewerID string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update payment status
	var billingID string
	err = tx.QueryRow(ctx, `
		UPDATE payments
		SET status = 'approved', id_reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
		RETURNING id_billing
	`, paymentID, reviewerID).Scan(&billingID)
	if err != nil {
		return fmt.Errorf("failed to approve payment: %w", err)
	}

	// Update billing status to paid
	_, err = tx.Exec(ctx, `
		UPDATE billings
		SET status = 'paid', paid_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, billingID)
	if err != nil {
		return fmt.Errorf("failed to update billing: %w", err)
	}

	return tx.Commit(ctx)
}

// RejectPayment rejects a payment and reverts billing status to unpaid
func RejectPayment(ctx context.Context, paymentID, reviewerID, reason string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update payment status
	var billingID string
	err = tx.QueryRow(ctx, `
		UPDATE payments
		SET status = 'rejected', rejection_reason = $3, id_reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
		RETURNING id_billing
	`, paymentID, reviewerID, reason).Scan(&billingID)
	if err != nil {
		return fmt.Errorf("failed to reject payment: %w", err)
	}

	// Revert billing status to unpaid
	_, err = tx.Exec(ctx, `
		UPDATE billings
		SET status = 'unpaid', updated_at = NOW()
		WHERE id = $1
	`, billingID)
	if err != nil {
		return fmt.Errorf("failed to update billing: %w", err)
	}

	return tx.Commit(ctx)
}
