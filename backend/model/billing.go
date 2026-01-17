package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Billing types
const (
	BillingTypeRegistration = "registration"
	BillingTypeTuition      = "tuition"
	BillingTypeDormitory    = "dormitory"
)

// Billing statuses
const (
	BillingStatusUnpaid              = "unpaid"
	BillingStatusPendingVerification = "pending_verification"
	BillingStatusPaid                = "paid"
)

// Billing represents a billing record for a candidate
type Billing struct {
	ID          string     `json:"id"`
	CandidateID string     `json:"candidate_id"`
	BillingType string     `json:"billing_type"`
	Description *string    `json:"description,omitempty"`
	Amount      int        `json:"amount"` // in IDR
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BillingWithPayment includes the latest payment info
type BillingWithPayment struct {
	Billing
	PaymentID       *string    `json:"payment_id,omitempty"`
	PaymentAmount   *int       `json:"payment_amount,omitempty"`
	PaymentStatus   *string    `json:"payment_status,omitempty"`
	PaymentProofURL *string    `json:"payment_proof_url,omitempty"`
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
}

// CreateBilling creates a new billing record
func CreateBilling(ctx context.Context, candidateID, billingType string, description *string, amount int, dueDate *time.Time) (*Billing, error) {
	var b Billing
	err := pool.QueryRow(ctx, `
		INSERT INTO billings (candidate_id, billing_type, description, amount, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, candidate_id, billing_type, description, amount, due_date, status, paid_at, created_at, updated_at
	`, candidateID, billingType, description, amount, dueDate).Scan(
		&b.ID, &b.CandidateID, &b.BillingType, &b.Description, &b.Amount,
		&b.DueDate, &b.Status, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing: %w", err)
	}
	return &b, nil
}

// FindBillingByID finds a billing by ID
func FindBillingByID(ctx context.Context, id string) (*Billing, error) {
	var b Billing
	err := pool.QueryRow(ctx, `
		SELECT id, candidate_id, billing_type, description, amount, due_date, status, paid_at, created_at, updated_at
		FROM billings
		WHERE id = $1
	`, id).Scan(
		&b.ID, &b.CandidateID, &b.BillingType, &b.Description, &b.Amount,
		&b.DueDate, &b.Status, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find billing: %w", err)
	}
	return &b, nil
}

// ListBillingsByCandidate returns all billings for a candidate with latest payment info
func ListBillingsByCandidate(ctx context.Context, candidateID string) ([]BillingWithPayment, error) {
	rows, err := pool.Query(ctx, `
		SELECT b.id, b.candidate_id, b.billing_type, b.description, b.amount, b.due_date,
		       b.status, b.paid_at, b.created_at, b.updated_at,
		       p.id as payment_id, p.amount as payment_amount, p.status as payment_status,
		       p.proof_file_path, p.transfer_date
		FROM billings b
		LEFT JOIN LATERAL (
			SELECT id, amount, status, proof_file_path, transfer_date
			FROM payments
			WHERE billing_id = b.id
			ORDER BY created_at DESC
			LIMIT 1
		) p ON true
		WHERE b.candidate_id = $1
		ORDER BY b.created_at
	`, candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list billings: %w", err)
	}
	defer rows.Close()

	var billings []BillingWithPayment
	for rows.Next() {
		var b BillingWithPayment
		err := rows.Scan(
			&b.ID, &b.CandidateID, &b.BillingType, &b.Description, &b.Amount, &b.DueDate,
			&b.Status, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt,
			&b.PaymentID, &b.PaymentAmount, &b.PaymentStatus, &b.PaymentProofURL, &b.PaymentDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan billing: %w", err)
		}
		billings = append(billings, b)
	}
	return billings, nil
}

// GetBillingSummary returns payment summary for a candidate
func GetBillingSummary(ctx context.Context, candidateID string) (totalDue, totalPaid, totalPending int, err error) {
	err = pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(amount), 0) as total_due,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN status = 'pending_verification' THEN amount ELSE 0 END), 0) as total_pending
		FROM billings
		WHERE candidate_id = $1
	`, candidateID).Scan(&totalDue, &totalPaid, &totalPending)
	if err != nil {
		err = fmt.Errorf("failed to get billing summary: %w", err)
	}
	return
}

// UpdateBillingStatus updates the billing status
func UpdateBillingStatus(ctx context.Context, id, status string) error {
	var paidAt *time.Time
	if status == BillingStatusPaid {
		now := time.Now()
		paidAt = &now
	}

	_, err := pool.Exec(ctx, `
		UPDATE billings
		SET status = $2, paid_at = $3, updated_at = NOW()
		WHERE id = $1
	`, id, status, paidAt)
	if err != nil {
		return fmt.Errorf("failed to update billing status: %w", err)
	}
	return nil
}

// CreateRegistrationBilling creates a registration fee billing for a new candidate
func CreateRegistrationBilling(ctx context.Context, candidateID string, amount int) (*Billing, error) {
	desc := "Biaya Pendaftaran"
	return CreateBilling(ctx, candidateID, BillingTypeRegistration, &desc, amount, nil)
}

// BillingTypeLabel returns the display label for a billing type
func BillingTypeLabel(billingType string) string {
	switch billingType {
	case BillingTypeRegistration:
		return "Biaya Pendaftaran"
	case BillingTypeTuition:
		return "Biaya Kuliah"
	case BillingTypeDormitory:
		return "Biaya Asrama"
	default:
		return billingType
	}
}
