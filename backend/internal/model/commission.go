package model

import (
	"context"
	"fmt"
	"time"
)

// Commission represents a commission entry in the ledger
type Commission struct {
	ID           string     `json:"id"`
	ReferrerID   string     `json:"id_referrer"`
	CandidateID  string     `json:"id_candidate"`
	TriggerEvent string     `json:"trigger_event"`
	Amount       int64      `json:"amount"`
	Status       string     `json:"status"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ApprovedBy   *string    `json:"id_approved_by,omitempty"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	PaidBy       *string    `json:"id_paid_by,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CommissionWithDetails includes referrer and candidate names
type CommissionWithDetails struct {
	Commission
	ReferrerName   string  `json:"referrer_name"`
	ReferrerType   string  `json:"referrer_type"`
	CandidateName  *string `json:"candidate_name,omitempty"`
	CandidateEmail *string `json:"candidate_email,omitempty"`
	ApprovedByName *string `json:"id_approved_by_name,omitempty"`
	PaidByName     *string `json:"id_paid_by_name,omitempty"`
}

// CommissionFilters for listing commissions
type CommissionFilters struct {
	ReferrerID string
	Status     string
	DateFrom   *time.Time
	DateTo     *time.Time
	Limit      int
	Offset     int
}

// CommissionStats for dashboard
type CommissionStats struct {
	TotalPending  int   `json:"total_pending"`
	TotalApproved int   `json:"total_approved"`
	TotalPaid     int   `json:"total_paid"`
	AmountPending int64 `json:"amount_pending"`
	AmountApproved int64 `json:"amount_approved"`
	AmountPaid    int64 `json:"amount_paid"`
}

// GetCommissionAmount calculates the commission amount for a referrer
// Uses commission_override if set, otherwise looks up from reward_configs
func GetCommissionAmount(ctx context.Context, referrerID string, triggerEvent string) (int64, error) {
	// First check if referrer has commission_override
	var commissionOverride *int64
	var referrerType string
	err := pool.QueryRow(ctx, `
		SELECT commission_override, type FROM referrers WHERE id = $1
	`, referrerID).Scan(&commissionOverride, &referrerType)
	if err != nil {
		return 0, fmt.Errorf("failed to get referrer: %w", err)
	}

	if commissionOverride != nil {
		return *commissionOverride, nil
	}

	// Look up from reward_configs
	var amount int64
	err = pool.QueryRow(ctx, `
		SELECT amount FROM reward_configs
		WHERE referrer_type = $1 AND trigger_event = $2 AND reward_type = 'cash' AND is_active = true
	`, referrerType, triggerEvent).Scan(&amount)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil // No commission configured
		}
		return 0, fmt.Errorf("failed to get reward config: %w", err)
	}

	return amount, nil
}

// CreateCommission creates a new commission entry
func CreateCommission(ctx context.Context, referrerID, candidateID, triggerEvent string, amount int64) (*Commission, error) {
	var c Commission
	err := pool.QueryRow(ctx, `
		INSERT INTO commission_ledger (id_referrer, id_candidate, trigger_event, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, id_referrer, id_candidate, trigger_event, amount, status, approved_at, id_approved_by, paid_at, id_paid_by, notes, created_at, updated_at
	`, referrerID, candidateID, triggerEvent, amount).Scan(
		&c.ID, &c.ReferrerID, &c.CandidateID, &c.TriggerEvent, &c.Amount, &c.Status,
		&c.ApprovedAt, &c.ApprovedBy, &c.PaidAt, &c.PaidBy, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}
	return &c, nil
}

// CreateCommissionForCandidate auto-creates commission when candidate reaches trigger event
func CreateCommissionForCandidate(ctx context.Context, candidateID, triggerEvent string) error {
	// Get candidate's referrer
	var referrerID *string
	err := pool.QueryRow(ctx, `SELECT id_referrer FROM candidates WHERE id = $1`, candidateID).Scan(&referrerID)
	if err != nil {
		return fmt.Errorf("failed to get candidate: %w", err)
	}

	if referrerID == nil {
		return nil // No referrer, no commission
	}

	// Get commission amount
	amount, err := GetCommissionAmount(ctx, *referrerID, triggerEvent)
	if err != nil {
		return fmt.Errorf("failed to get commission amount: %w", err)
	}

	if amount == 0 {
		return nil // No commission configured
	}

	// Create commission (ignore duplicate errors)
	_, err = CreateCommission(ctx, *referrerID, candidateID, triggerEvent, amount)
	if err != nil {
		// Check if it's a duplicate key error (commission already exists)
		if isDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	return nil
}

// FindCommissionByID finds a commission by ID
func FindCommissionByID(ctx context.Context, id string) (*CommissionWithDetails, error) {
	var c CommissionWithDetails
	err := pool.QueryRow(ctx, `
		SELECT cl.id, cl.id_referrer, cl.id_candidate, cl.trigger_event, cl.amount, cl.status,
		       cl.approved_at, cl.id_approved_by, cl.paid_at, cl.id_paid_by, cl.notes, cl.created_at, cl.updated_at,
		       r.name as referrer_name, r.type as referrer_type,
		       ca.name as candidate_name, ca.email as candidate_email,
		       u1.name as id_approved_by_name, u2.name as id_paid_by_name
		FROM commission_ledger cl
		JOIN referrers r ON cl.id_referrer = r.id
		LEFT JOIN candidates ca ON cl.id_candidate = ca.id
		LEFT JOIN users u1 ON cl.id_approved_by = u1.id
		LEFT JOIN users u2 ON cl.id_paid_by = u2.id
		WHERE cl.id = $1
	`, id).Scan(
		&c.ID, &c.ReferrerID, &c.CandidateID, &c.TriggerEvent, &c.Amount, &c.Status,
		&c.ApprovedAt, &c.ApprovedBy, &c.PaidAt, &c.PaidBy, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
		&c.ReferrerName, &c.ReferrerType, &c.CandidateName, &c.CandidateEmail,
		&c.ApprovedByName, &c.PaidByName,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find commission: %w", err)
	}
	return &c, nil
}

// ListCommissions lists commissions with filters
func ListCommissions(ctx context.Context, filters CommissionFilters) ([]CommissionWithDetails, int, error) {
	// Build query
	baseQuery := `
		FROM commission_ledger cl
		JOIN referrers r ON cl.id_referrer = r.id
		LEFT JOIN candidates ca ON cl.id_candidate = ca.id
		LEFT JOIN users u1 ON cl.id_approved_by = u1.id
		LEFT JOIN users u2 ON cl.id_paid_by = u2.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filters.ReferrerID != "" {
		baseQuery += fmt.Sprintf(" AND cl.id_referrer = $%d", argIdx)
		args = append(args, filters.ReferrerID)
		argIdx++
	}

	if filters.Status != "" {
		baseQuery += fmt.Sprintf(" AND cl.status = $%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}

	if filters.DateFrom != nil {
		baseQuery += fmt.Sprintf(" AND cl.created_at >= $%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}

	if filters.DateTo != nil {
		baseQuery += fmt.Sprintf(" AND cl.created_at <= $%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}

	// Count total
	var total int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count commissions: %w", err)
	}

	// Get commissions
	query := `
		SELECT cl.id, cl.id_referrer, cl.id_candidate, cl.trigger_event, cl.amount, cl.status,
		       cl.approved_at, cl.id_approved_by, cl.paid_at, cl.id_paid_by, cl.notes, cl.created_at, cl.updated_at,
		       r.name as referrer_name, r.type as referrer_type,
		       ca.name as candidate_name, ca.email as candidate_email,
		       u1.name as id_approved_by_name, u2.name as id_paid_by_name
	` + baseQuery + " ORDER BY cl.created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filters.Limit)
	}
	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filters.Offset)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list commissions: %w", err)
	}
	defer rows.Close()

	var commissions []CommissionWithDetails
	for rows.Next() {
		var c CommissionWithDetails
		err := rows.Scan(
			&c.ID, &c.ReferrerID, &c.CandidateID, &c.TriggerEvent, &c.Amount, &c.Status,
			&c.ApprovedAt, &c.ApprovedBy, &c.PaidAt, &c.PaidBy, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
			&c.ReferrerName, &c.ReferrerType, &c.CandidateName, &c.CandidateEmail,
			&c.ApprovedByName, &c.PaidByName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan commission: %w", err)
		}
		commissions = append(commissions, c)
	}

	return commissions, total, nil
}

// ListPendingCommissions lists pending commissions
func ListPendingCommissions(ctx context.Context) ([]CommissionWithDetails, error) {
	commissions, _, err := ListCommissions(ctx, CommissionFilters{Status: "pending", Limit: 100})
	return commissions, err
}

// ListCommissionsByReferrer lists commissions for a specific referrer
func ListCommissionsByReferrer(ctx context.Context, referrerID string) ([]CommissionWithDetails, error) {
	commissions, _, err := ListCommissions(ctx, CommissionFilters{ReferrerID: referrerID, Limit: 100})
	return commissions, err
}

// ApproveCommission approves a pending commission
func ApproveCommission(ctx context.Context, id, approvedBy string) error {
	result, err := pool.Exec(ctx, `
		UPDATE commission_ledger
		SET status = 'approved', approved_at = NOW(), id_approved_by = $1, updated_at = NOW()
		WHERE id = $2 AND status = 'pending'
	`, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to approve commission: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("commission not found or already processed")
	}
	return nil
}

// MarkCommissionPaid marks an approved commission as paid
func MarkCommissionPaid(ctx context.Context, id, paidBy string, notes *string) error {
	result, err := pool.Exec(ctx, `
		UPDATE commission_ledger
		SET status = 'paid', paid_at = NOW(), id_paid_by = $1, notes = COALESCE($2, notes), updated_at = NOW()
		WHERE id = $3 AND status = 'approved'
	`, paidBy, notes, id)
	if err != nil {
		return fmt.Errorf("failed to mark commission paid: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("commission not found or not approved")
	}
	return nil
}

// CancelCommission cancels a pending commission
func CancelCommission(ctx context.Context, id string, notes *string) error {
	result, err := pool.Exec(ctx, `
		UPDATE commission_ledger
		SET status = 'cancelled', notes = COALESCE($1, notes), updated_at = NOW()
		WHERE id = $2 AND status = 'pending'
	`, notes, id)
	if err != nil {
		return fmt.Errorf("failed to cancel commission: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("commission not found or already processed")
	}
	return nil
}

// BatchApproveCommissions approves multiple commissions
func BatchApproveCommissions(ctx context.Context, ids []string, approvedBy string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	result, err := pool.Exec(ctx, `
		UPDATE commission_ledger
		SET status = 'approved', approved_at = NOW(), id_approved_by = $1, updated_at = NOW()
		WHERE id = ANY($2) AND status = 'pending'
	`, approvedBy, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to batch approve commissions: %w", err)
	}
	return int(result.RowsAffected()), nil
}

// BatchMarkCommissionsPaid marks multiple commissions as paid
func BatchMarkCommissionsPaid(ctx context.Context, ids []string, paidBy string, notes *string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	result, err := pool.Exec(ctx, `
		UPDATE commission_ledger
		SET status = 'paid', paid_at = NOW(), id_paid_by = $1, notes = COALESCE($2, notes), updated_at = NOW()
		WHERE id = ANY($3) AND status = 'approved'
	`, paidBy, notes, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to batch mark commissions paid: %w", err)
	}
	return int(result.RowsAffected()), nil
}

// GetCommissionStats returns commission statistics
func GetCommissionStats(ctx context.Context) (*CommissionStats, error) {
	var stats CommissionStats
	err := pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as total_pending,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0) as total_approved,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as amount_pending,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN amount ELSE 0 END), 0) as amount_approved,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) as amount_paid
		FROM commission_ledger
		WHERE status != 'cancelled'
	`).Scan(
		&stats.TotalPending, &stats.TotalApproved, &stats.TotalPaid,
		&stats.AmountPending, &stats.AmountApproved, &stats.AmountPaid,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get commission stats: %w", err)
	}
	return &stats, nil
}

// isDuplicateKeyError checks if error is a duplicate key violation
func isDuplicateKeyError(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CommissionExportData contains data needed for CSV export
type CommissionExportData struct {
	ID             string
	ReferrerName   string
	ReferrerType   string
	BankName       string
	BankAccount    string
	AccountHolder  string
	Amount         int64
	CandidateName  string
	TriggerEvent   string
	ApprovedAt     *time.Time
}

// ListCommissionsForExport returns approved commissions with bank details for CSV export
func ListCommissionsForExport(ctx context.Context, status string) ([]CommissionExportData, error) {
	rows, err := pool.Query(ctx, `
		SELECT cl.id, r.name as referrer_name, r.type as referrer_type,
		       r.bank_name, r.bank_account, r.account_holder,
		       cl.amount, ca.name as candidate_name, cl.trigger_event, cl.approved_at
		FROM commission_ledger cl
		JOIN referrers r ON cl.id_referrer = r.id
		LEFT JOIN candidates ca ON cl.id_candidate = ca.id
		WHERE cl.status = $1
		ORDER BY cl.approved_at ASC
	`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list commissions for export: %w", err)
	}
	defer rows.Close()

	var commissions []CommissionExportData
	for rows.Next() {
		var c CommissionExportData
		var bankName, bankAccount, accountHolder, candidateName *string
		err := rows.Scan(
			&c.ID, &c.ReferrerName, &c.ReferrerType,
			&bankName, &bankAccount, &accountHolder,
			&c.Amount, &candidateName, &c.TriggerEvent, &c.ApprovedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan commission export: %w", err)
		}

		// Decrypt referrer fields using the proper decryption function
		nameDec, _, _, bankNameDec, bankAccountDec, accountHolderDec, _, decryptErr := decryptReferrerFields(
			c.ReferrerName, nil, nil, bankName, bankAccount, accountHolder, nil,
		)
		if decryptErr == nil {
			c.ReferrerName = nameDec
			if bankNameDec != nil {
				c.BankName = *bankNameDec
			}
			if bankAccountDec != nil {
				c.BankAccount = *bankAccountDec
			}
			if accountHolderDec != nil {
				c.AccountHolder = *accountHolderDec
			}
		}

		// Decrypt candidate name
		if candidateName != nil {
			decrypted, err := decryptName(*candidateName)
			if err == nil {
				c.CandidateName = decrypted
			} else {
				c.CandidateName = *candidateName
			}
		}

		commissions = append(commissions, c)
	}
	return commissions, nil
}
