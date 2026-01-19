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
	BillingStatusCancelled           = "cancelled"
)

// Billing represents a billing record for a candidate
type Billing struct {
	ID          string     `json:"id"`
	CandidateID string     `json:"id_candidate"`
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
		INSERT INTO billings (id_candidate, billing_type, description, amount, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, id_candidate, billing_type, description, amount, due_date, status, paid_at, created_at, updated_at
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
		SELECT id, id_candidate, billing_type, description, amount, due_date, status, paid_at, created_at, updated_at
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
		SELECT b.id, b.id_candidate, b.billing_type, b.description, b.amount, b.due_date,
		       b.status, b.paid_at, b.created_at, b.updated_at,
		       p.id as payment_id, p.amount as payment_amount, p.status as payment_status,
		       p.proof_file_path, p.transfer_date
		FROM billings b
		LEFT JOIN LATERAL (
			SELECT id, amount, status, proof_file_path, transfer_date
			FROM payments
			WHERE id_billing = b.id
			ORDER BY created_at DESC
			LIMIT 1
		) p ON true
		WHERE b.id_candidate = $1
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
		WHERE id_candidate = $1
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

// BillingStatusLabel returns the display label for a billing status
func BillingStatusLabel(status string) string {
	switch status {
	case BillingStatusUnpaid:
		return "Belum Dibayar"
	case BillingStatusPendingVerification:
		return "Menunggu Verifikasi"
	case BillingStatusPaid:
		return "Lunas"
	case BillingStatusCancelled:
		return "Dibatalkan"
	default:
		return status
	}
}

// BillingWithCandidate includes candidate info for finance listing
type BillingWithCandidate struct {
	Billing
	CandidateName  *string `json:"candidate_name,omitempty"`
	CandidateEmail string  `json:"candidate_email"`
	ProdiName      *string `json:"prodi_name,omitempty"`
}

// BillingFilters contains filters for listing billings
type BillingFilters struct {
	Search      string // search by candidate name or email
	Status      string // filter by status
	BillingType string // filter by billing type
	Page        int
	PageSize    int
}

// ListAllBillings returns all billings with filters for finance
func ListAllBillings(ctx context.Context, filters BillingFilters) ([]BillingWithCandidate, int, error) {
	// Build query
	baseQuery := `
		FROM billings b
		JOIN candidates c ON c.id = b.id_candidate
		LEFT JOIN prodis p ON p.id = c.id_prodi
		WHERE 1=1
	`
	args := []any{}
	argCount := 0

	if filters.Search != "" {
		// With encryption, search by encrypted email (deterministic) for exact match
		emailEnc, err := encryptEmail(filters.Search)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to encrypt search term: %w", err)
		}
		argCount++
		baseQuery += fmt.Sprintf(" AND c.email = $%d", argCount)
		args = append(args, emailEnc)
	}
	if filters.Status != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND b.status = $%d", argCount)
		args = append(args, filters.Status)
	}
	if filters.BillingType != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND b.billing_type = $%d", argCount)
		args = append(args, filters.BillingType)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count billings: %w", err)
	}

	// Pagination
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	offset := (filters.Page - 1) * filters.PageSize

	// Get billings
	selectQuery := `
		SELECT b.id, b.id_candidate, b.billing_type, b.description, b.amount, b.due_date,
		       b.status, b.paid_at, b.created_at, b.updated_at,
		       c.name, c.email, p.name as prodi_name
	` + baseQuery + fmt.Sprintf(" ORDER BY b.created_at DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, filters.PageSize, offset)

	rows, err := pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list billings: %w", err)
	}
	defer rows.Close()

	var billings []BillingWithCandidate
	for rows.Next() {
		var b BillingWithCandidate
		err := rows.Scan(
			&b.ID, &b.CandidateID, &b.BillingType, &b.Description, &b.Amount, &b.DueDate,
			&b.Status, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt,
			&b.CandidateName, &b.CandidateEmail, &b.ProdiName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan billing: %w", err)
		}

		// Decrypt candidate fields
		b.CandidateName, _ = decryptNullableP(b.CandidateName)
		if emailDec, err := decryptNullableD(&b.CandidateEmail); err == nil && emailDec != nil {
			b.CandidateEmail = *emailDec
		}

		billings = append(billings, b)
	}
	return billings, total, nil
}

// UpdateBilling updates billing amount and due date (only for unpaid billings)
func UpdateBilling(ctx context.Context, id string, amount int, dueDate *time.Time, description *string) error {
	result, err := pool.Exec(ctx, `
		UPDATE billings
		SET amount = $2, due_date = $3, description = $4, updated_at = NOW()
		WHERE id = $1 AND status = 'unpaid'
	`, id, amount, dueDate, description)
	if err != nil {
		return fmt.Errorf("failed to update billing: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("billing not found or not in unpaid status")
	}
	return nil
}

// CancelBilling sets billing status to cancelled (only for unpaid billings)
func CancelBilling(ctx context.Context, id string) error {
	result, err := pool.Exec(ctx, `
		UPDATE billings
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1 AND status = 'unpaid'
	`, id)
	if err != nil {
		return fmt.Errorf("failed to cancel billing: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("billing not found or not in unpaid status")
	}
	return nil
}

// GetBillingStats returns billing statistics for finance dashboard
func GetBillingStats(ctx context.Context) (unpaid, pending, paid, cancelled int, err error) {
	err = pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'unpaid' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'pending_verification' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END), 0)
		FROM billings
	`).Scan(&unpaid, &pending, &paid, &cancelled)
	if err != nil {
		err = fmt.Errorf("failed to get billing stats: %w", err)
	}
	return
}

// CreateTuitionBilling creates tuition billing for a committed candidate
// installments is the number of payment installments (currently only 1 is supported)
func CreateTuitionBilling(ctx context.Context, candidateID string, amount int64, academicYear string, installments int) (*Billing, error) {
	desc := fmt.Sprintf("Biaya Kuliah %s", academicYear)
	// For now, single payment - due date is 30 days from now
	dueDate := time.Now().AddDate(0, 0, 30)
	return CreateBilling(ctx, candidateID, BillingTypeTuition, &desc, int(amount), &dueDate)
}

// CreateDormitoryBillings creates dormitory billings for a committed candidate
// installments determines how many billing records to create (1, 2, or 10)
func CreateDormitoryBillings(ctx context.Context, candidateID string, totalAmount int64, academicYear string, installments int) ([]*Billing, error) {
	if installments <= 0 {
		installments = 1
	}
	if installments > 10 {
		installments = 10
	}

	amountPerInstallment := totalAmount / int64(installments)
	remainder := totalAmount % int64(installments)

	var billings []*Billing
	for i := 1; i <= installments; i++ {
		amount := amountPerInstallment
		// Add remainder to last installment
		if i == installments {
			amount += remainder
		}

		desc := fmt.Sprintf("Biaya Asrama %s", academicYear)
		if installments > 1 {
			desc = fmt.Sprintf("Biaya Asrama %s - Cicilan %d/%d", academicYear, i, installments)
		}

		// Due date: first installment 30 days, subsequent installments 30 days apart
		dueDate := time.Now().AddDate(0, 0, 30*i)

		billing, err := CreateBilling(ctx, candidateID, BillingTypeDormitory, &desc, int(amount), &dueDate)
		if err != nil {
			return billings, fmt.Errorf("failed to create dormitory billing %d: %w", i, err)
		}
		billings = append(billings, billing)
	}

	return billings, nil
}

// GetTuitionFeeForCandidate gets the tuition fee amount for a candidate's prodi and academic year
func GetTuitionFeeForCandidate(ctx context.Context, candidateID, academicYear string) (int64, error) {
	var amount int64
	err := pool.QueryRow(ctx, `
		SELECT fs.amount
		FROM fee_structures fs
		JOIN fee_types ft ON ft.id = fs.id_fee_type
		JOIN candidates c ON c.id_prodi = fs.id_prodi
		WHERE c.id = $1 AND ft.code = 'tuition' AND fs.academic_year = $2 AND fs.is_active = true
	`, candidateID, academicYear).Scan(&amount)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get tuition fee: %w", err)
	}
	return amount, nil
}

// GetDormitoryFee gets the dormitory fee amount for an academic year (global, not per-prodi)
func GetDormitoryFee(ctx context.Context, academicYear string) (int64, error) {
	var amount int64
	err := pool.QueryRow(ctx, `
		SELECT fs.amount
		FROM fee_structures fs
		JOIN fee_types ft ON ft.id = fs.id_fee_type
		WHERE ft.code = 'dormitory' AND fs.id_prodi IS NULL AND fs.academic_year = $1 AND fs.is_active = true
	`, academicYear).Scan(&amount)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get dormitory fee: %w", err)
	}
	return amount, nil
}
