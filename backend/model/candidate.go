package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Candidate represents a prospective student
type Candidate struct {
	ID                    string     `json:"id"`
	Email                 *string    `json:"email,omitempty"`
	EmailVerified         bool       `json:"email_verified"`
	Phone                 *string    `json:"phone,omitempty"`
	PhoneVerified         bool       `json:"phone_verified"`
	PasswordHash          string     `json:"-"`
	Name                  *string    `json:"name,omitempty"`
	Address               *string    `json:"address,omitempty"`
	City                  *string    `json:"city,omitempty"`
	Province              *string    `json:"province,omitempty"`
	HighSchool            *string    `json:"high_school,omitempty"`
	GraduationYear        *int       `json:"graduation_year,omitempty"`
	ProdiID               *string    `json:"prodi_id,omitempty"`
	CampaignID            *string    `json:"campaign_id,omitempty"`
	ReferrerID            *string    `json:"referrer_id,omitempty"`
	ReferredByCandidateID *string    `json:"referred_by_candidate_id,omitempty"`
	SourceType            *string    `json:"source_type,omitempty"`
	SourceDetail          *string    `json:"source_detail,omitempty"`
	AssignedConsultantID  *string    `json:"assigned_consultant_id,omitempty"`
	Status                string     `json:"status"`
	LostReasonID          *string    `json:"lost_reason_id,omitempty"`
	LostAt                *time.Time `json:"lost_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// CandidateWithDetails includes related names for display
type CandidateWithDetails struct {
	Candidate
	ProdiName      *string `json:"prodi_name,omitempty"`
	ConsultantName *string `json:"consultant_name,omitempty"`
	CampaignName   *string `json:"campaign_name,omitempty"`
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateCandidate creates a new candidate with email/phone and hashed password
func CreateCandidate(ctx context.Context, email, phone, passwordHash string) (*Candidate, error) {
	var candidate Candidate
	var emailPtr, phonePtr *string
	if email != "" {
		emailPtr = &email
	}
	if phone != "" {
		phonePtr = &phone
	}

	err := pool.QueryRow(ctx, `
		INSERT INTO candidates (email, phone, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		          high_school, graduation_year, prodi_id, campaign_id, referrer_id, referred_by_candidate_id,
		          source_type, source_detail, assigned_consultant_id, status, lost_reason_id, lost_at, created_at, updated_at
	`, emailPtr, phonePtr, passwordHash).Scan(
		&candidate.ID, &candidate.Email, &candidate.EmailVerified, &candidate.Phone, &candidate.PhoneVerified,
		&candidate.PasswordHash, &candidate.Name, &candidate.Address, &candidate.City, &candidate.Province,
		&candidate.HighSchool, &candidate.GraduationYear, &candidate.ProdiID, &candidate.CampaignID,
		&candidate.ReferrerID, &candidate.ReferredByCandidateID, &candidate.SourceType, &candidate.SourceDetail,
		&candidate.AssignedConsultantID, &candidate.Status, &candidate.LostReasonID, &candidate.LostAt,
		&candidate.CreatedAt, &candidate.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create candidate: %w", err)
	}
	return &candidate, nil
}

// FindCandidateByID finds a candidate by ID
func FindCandidateByID(ctx context.Context, id string) (*Candidate, error) {
	var candidate Candidate
	err := pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, prodi_id, campaign_id, referrer_id, referred_by_candidate_id,
		       source_type, source_detail, assigned_consultant_id, status, lost_reason_id, lost_at, created_at, updated_at
		FROM candidates WHERE id = $1
	`, id).Scan(
		&candidate.ID, &candidate.Email, &candidate.EmailVerified, &candidate.Phone, &candidate.PhoneVerified,
		&candidate.PasswordHash, &candidate.Name, &candidate.Address, &candidate.City, &candidate.Province,
		&candidate.HighSchool, &candidate.GraduationYear, &candidate.ProdiID, &candidate.CampaignID,
		&candidate.ReferrerID, &candidate.ReferredByCandidateID, &candidate.SourceType, &candidate.SourceDetail,
		&candidate.AssignedConsultantID, &candidate.Status, &candidate.LostReasonID, &candidate.LostAt,
		&candidate.CreatedAt, &candidate.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find candidate by id: %w", err)
	}
	return &candidate, nil
}

// FindCandidateByEmail finds a candidate by email
func FindCandidateByEmail(ctx context.Context, email string) (*Candidate, error) {
	var candidate Candidate
	err := pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, prodi_id, campaign_id, referrer_id, referred_by_candidate_id,
		       source_type, source_detail, assigned_consultant_id, status, lost_reason_id, lost_at, created_at, updated_at
		FROM candidates WHERE email = $1
	`, email).Scan(
		&candidate.ID, &candidate.Email, &candidate.EmailVerified, &candidate.Phone, &candidate.PhoneVerified,
		&candidate.PasswordHash, &candidate.Name, &candidate.Address, &candidate.City, &candidate.Province,
		&candidate.HighSchool, &candidate.GraduationYear, &candidate.ProdiID, &candidate.CampaignID,
		&candidate.ReferrerID, &candidate.ReferredByCandidateID, &candidate.SourceType, &candidate.SourceDetail,
		&candidate.AssignedConsultantID, &candidate.Status, &candidate.LostReasonID, &candidate.LostAt,
		&candidate.CreatedAt, &candidate.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find candidate by email: %w", err)
	}
	return &candidate, nil
}

// FindCandidateByPhone finds a candidate by phone
func FindCandidateByPhone(ctx context.Context, phone string) (*Candidate, error) {
	var candidate Candidate
	err := pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, prodi_id, campaign_id, referrer_id, referred_by_candidate_id,
		       source_type, source_detail, assigned_consultant_id, status, lost_reason_id, lost_at, created_at, updated_at
		FROM candidates WHERE phone = $1
	`, phone).Scan(
		&candidate.ID, &candidate.Email, &candidate.EmailVerified, &candidate.Phone, &candidate.PhoneVerified,
		&candidate.PasswordHash, &candidate.Name, &candidate.Address, &candidate.City, &candidate.Province,
		&candidate.HighSchool, &candidate.GraduationYear, &candidate.ProdiID, &candidate.CampaignID,
		&candidate.ReferrerID, &candidate.ReferredByCandidateID, &candidate.SourceType, &candidate.SourceDetail,
		&candidate.AssignedConsultantID, &candidate.Status, &candidate.LostReasonID, &candidate.LostAt,
		&candidate.CreatedAt, &candidate.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find candidate by phone: %w", err)
	}
	return &candidate, nil
}

// AuthenticateCandidate authenticates a candidate by email or phone with password
func AuthenticateCandidate(ctx context.Context, identifier, password string) (*Candidate, error) {
	// Try to find by email first, then by phone
	candidate, err := FindCandidateByEmail(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		candidate, err = FindCandidateByPhone(ctx, identifier)
		if err != nil {
			return nil, err
		}
	}
	if candidate == nil {
		return nil, nil // Not found
	}

	// Check password
	if !CheckPassword(password, candidate.PasswordHash) {
		return nil, nil // Wrong password
	}

	return candidate, nil
}

// UpdateCandidatePersonalInfo updates personal information
func UpdateCandidatePersonalInfo(ctx context.Context, id, name, address, city, province string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates
		SET name = $2, address = $3, city = $4, province = $5, updated_at = NOW()
		WHERE id = $1
	`, id, name, address, city, province)
	if err != nil {
		return fmt.Errorf("failed to update personal info: %w", err)
	}
	return nil
}

// UpdateCandidateEducation updates education information
func UpdateCandidateEducation(ctx context.Context, id, highSchool string, graduationYear int, prodiID string) error {
	var prodiPtr *string
	if prodiID != "" {
		prodiPtr = &prodiID
	}
	_, err := pool.Exec(ctx, `
		UPDATE candidates
		SET high_school = $2, graduation_year = $3, prodi_id = $4, updated_at = NOW()
		WHERE id = $1
	`, id, highSchool, graduationYear, prodiPtr)
	if err != nil {
		return fmt.Errorf("failed to update education: %w", err)
	}
	return nil
}

// UpdateCandidateSourceTracking updates source tracking information
func UpdateCandidateSourceTracking(ctx context.Context, id, sourceType, sourceDetail, campaignID, referrerID, referredByCandidateID string) error {
	var campaignPtr, referrerPtr, referredByPtr *string
	if campaignID != "" {
		campaignPtr = &campaignID
	}
	if referrerID != "" {
		referrerPtr = &referrerID
	}
	if referredByCandidateID != "" {
		referredByPtr = &referredByCandidateID
	}
	var sourceDetailPtr *string
	if sourceDetail != "" {
		sourceDetailPtr = &sourceDetail
	}

	_, err := pool.Exec(ctx, `
		UPDATE candidates
		SET source_type = $2, source_detail = $3, campaign_id = $4, referrer_id = $5, referred_by_candidate_id = $6, updated_at = NOW()
		WHERE id = $1
	`, id, sourceType, sourceDetailPtr, campaignPtr, referrerPtr, referredByPtr)
	if err != nil {
		return fmt.Errorf("failed to update source tracking: %w", err)
	}
	return nil
}

// SetCandidateEmailVerified marks email as verified
func SetCandidateEmailVerified(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET email_verified = true, updated_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}
	return nil
}

// SetCandidatePhoneVerified marks phone as verified
func SetCandidatePhoneVerified(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET phone_verified = true, updated_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to set phone verified: %w", err)
	}
	return nil
}

// AssignCandidateConsultant assigns a consultant to a candidate
func AssignCandidateConsultant(ctx context.Context, id, consultantID string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET assigned_consultant_id = $2, updated_at = NOW() WHERE id = $1
	`, id, consultantID)
	if err != nil {
		return fmt.Errorf("failed to assign consultant: %w", err)
	}
	return nil
}

// UpdateCandidateStatus updates the candidate status
func UpdateCandidateStatus(ctx context.Context, id, status string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET status = $2, updated_at = NOW() WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

// MarkCandidateLost marks a candidate as lost with reason
func MarkCandidateLost(ctx context.Context, id, lostReasonID string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET status = 'lost', lost_reason_id = $2, lost_at = NOW(), updated_at = NOW() WHERE id = $1
	`, id, lostReasonID)
	if err != nil {
		return fmt.Errorf("failed to mark candidate as lost: %w", err)
	}
	return nil
}

// CandidateDashboardData contains all data needed for dashboard display
type CandidateDashboardData struct {
	Candidate
	ProdiName        *string `json:"prodi_name,omitempty"`
	ConsultantName   *string `json:"consultant_name,omitempty"`
	ConsultantEmail  *string `json:"consultant_email,omitempty"`
}

// CandidateDetailData contains all data needed for admin candidate detail view
type CandidateDetailData struct {
	Candidate
	ProdiName      *string `json:"prodi_name,omitempty"`
	ConsultantName *string `json:"consultant_name,omitempty"`
	CampaignName   *string `json:"campaign_name,omitempty"`
	ReferrerName   *string `json:"referrer_name,omitempty"`
	LostReasonName *string `json:"lost_reason_name,omitempty"`
}

// GetCandidateDetailData gets candidate with all related data for admin detail view
func GetCandidateDetailData(ctx context.Context, id string) (*CandidateDetailData, error) {
	var data CandidateDetailData
	err := pool.QueryRow(ctx, `
		SELECT c.id, c.email, c.email_verified, c.phone, c.phone_verified, c.password_hash, c.name, c.address, c.city, c.province,
		       c.high_school, c.graduation_year, c.prodi_id, c.campaign_id, c.referrer_id, c.referred_by_candidate_id,
		       c.source_type, c.source_detail, c.assigned_consultant_id, c.status, c.lost_reason_id, c.lost_at, c.created_at, c.updated_at,
		       p.name as prodi_name,
		       u.name as consultant_name,
		       camp.name as campaign_name,
		       ref.name as referrer_name,
		       lr.name as lost_reason_name
		FROM candidates c
		LEFT JOIN prodis p ON p.id = c.prodi_id
		LEFT JOIN users u ON u.id = c.assigned_consultant_id
		LEFT JOIN campaigns camp ON camp.id = c.campaign_id
		LEFT JOIN referrers ref ON ref.id = c.referrer_id
		LEFT JOIN lost_reasons lr ON lr.id = c.lost_reason_id
		WHERE c.id = $1
	`, id).Scan(
		&data.ID, &data.Email, &data.EmailVerified, &data.Phone, &data.PhoneVerified,
		&data.PasswordHash, &data.Name, &data.Address, &data.City, &data.Province,
		&data.HighSchool, &data.GraduationYear, &data.ProdiID, &data.CampaignID,
		&data.ReferrerID, &data.ReferredByCandidateID, &data.SourceType, &data.SourceDetail,
		&data.AssignedConsultantID, &data.Status, &data.LostReasonID, &data.LostAt,
		&data.CreatedAt, &data.UpdatedAt,
		&data.ProdiName, &data.ConsultantName, &data.CampaignName, &data.ReferrerName, &data.LostReasonName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate detail data: %w", err)
	}
	return &data, nil
}

// GetCandidateDashboardData gets candidate with related details for dashboard
func GetCandidateDashboardData(ctx context.Context, id string) (*CandidateDashboardData, error) {
	var data CandidateDashboardData
	err := pool.QueryRow(ctx, `
		SELECT c.id, c.email, c.email_verified, c.phone, c.phone_verified, c.password_hash, c.name, c.address, c.city, c.province,
		       c.high_school, c.graduation_year, c.prodi_id, c.campaign_id, c.referrer_id, c.referred_by_candidate_id,
		       c.source_type, c.source_detail, c.assigned_consultant_id, c.status, c.lost_reason_id, c.lost_at, c.created_at, c.updated_at,
		       p.name as prodi_name, u.name as consultant_name, u.email as consultant_email
		FROM candidates c
		LEFT JOIN prodis p ON p.id = c.prodi_id
		LEFT JOIN users u ON u.id = c.assigned_consultant_id
		WHERE c.id = $1
	`, id).Scan(
		&data.ID, &data.Email, &data.EmailVerified, &data.Phone, &data.PhoneVerified,
		&data.PasswordHash, &data.Name, &data.Address, &data.City, &data.Province,
		&data.HighSchool, &data.GraduationYear, &data.ProdiID, &data.CampaignID,
		&data.ReferrerID, &data.ReferredByCandidateID, &data.SourceType, &data.SourceDetail,
		&data.AssignedConsultantID, &data.Status, &data.LostReasonID, &data.LostAt,
		&data.CreatedAt, &data.UpdatedAt,
		&data.ProdiName, &data.ConsultantName, &data.ConsultantEmail,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate dashboard data: %w", err)
	}
	return &data, nil
}

// CandidateListItem is a trimmed version for list display
type CandidateListItem struct {
	ID               string     `json:"id"`
	Name             *string    `json:"name,omitempty"`
	Email            *string    `json:"email,omitempty"`
	Phone            *string    `json:"phone,omitempty"`
	Status           string     `json:"status"`
	ProdiName        *string    `json:"prodi_name,omitempty"`
	ConsultantName   *string    `json:"consultant_name,omitempty"`
	CampaignName     *string    `json:"campaign_name,omitempty"`
	SourceType       *string    `json:"source_type,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	LastInteraction  *time.Time `json:"last_interaction,omitempty"`
	NextFollowup     *time.Time `json:"next_followup,omitempty"`
}

// CandidateListFilters holds filter parameters for listing candidates
type CandidateListFilters struct {
	Status       string // Filter by status
	ConsultantID string // Filter by assigned consultant
	ProdiID      string // Filter by prodi
	CampaignID   string // Filter by campaign
	SourceType   string // Filter by source type
	Search       string // Search by name, email, or phone
	SortBy       string // Sort column: created_at, name, status, last_interaction
	SortOrder    string // asc or desc
	Limit        int    // Pagination limit
	Offset       int    // Pagination offset
}

// CandidateListResult contains list results with pagination info
type CandidateListResult struct {
	Candidates []CandidateListItem `json:"candidates"`
	Total      int                 `json:"total"`
	Limit      int                 `json:"limit"`
	Offset     int                 `json:"offset"`
}

// ListCandidates lists candidates with filters and pagination
// If consultantID is provided in visibility, only shows that consultant's candidates
// If supervisorID is provided in visibility, shows all candidates under that supervisor's team
func ListCandidates(ctx context.Context, filters CandidateListFilters, visibilityConsultantID, visibilitySupervisorID *string) (*CandidateListResult, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	// Role-based visibility
	if visibilityConsultantID != nil && *visibilityConsultantID != "" {
		// Consultant sees only their candidates
		whereClause += fmt.Sprintf(" AND c.assigned_consultant_id = $%d", argNum)
		args = append(args, *visibilityConsultantID)
		argNum++
	} else if visibilitySupervisorID != nil && *visibilitySupervisorID != "" {
		// Supervisor sees their team's candidates
		whereClause += fmt.Sprintf(" AND c.assigned_consultant_id IN (SELECT id FROM users WHERE id_supervisor = $%d OR id = $%d)", argNum, argNum+1)
		args = append(args, *visibilitySupervisorID, *visibilitySupervisorID)
		argNum += 2
	}
	// Admin sees all (no filter)

	// Apply filters
	if filters.Status != "" {
		whereClause += fmt.Sprintf(" AND c.status = $%d", argNum)
		args = append(args, filters.Status)
		argNum++
	}

	if filters.ConsultantID != "" {
		whereClause += fmt.Sprintf(" AND c.assigned_consultant_id = $%d", argNum)
		args = append(args, filters.ConsultantID)
		argNum++
	}

	if filters.ProdiID != "" {
		whereClause += fmt.Sprintf(" AND c.prodi_id = $%d", argNum)
		args = append(args, filters.ProdiID)
		argNum++
	}

	if filters.CampaignID != "" {
		whereClause += fmt.Sprintf(" AND c.campaign_id = $%d", argNum)
		args = append(args, filters.CampaignID)
		argNum++
	}

	if filters.SourceType != "" {
		whereClause += fmt.Sprintf(" AND c.source_type = $%d", argNum)
		args = append(args, filters.SourceType)
		argNum++
	}

	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		whereClause += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.email ILIKE $%d OR c.phone ILIKE $%d)", argNum, argNum, argNum)
		args = append(args, search)
		argNum++
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM candidates c
		%s
	`, whereClause)

	var total int
	err := pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count candidates: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "c.created_at" // Default
	switch filters.SortBy {
	case "name":
		orderBy = "c.name"
	case "status":
		orderBy = "c.status"
	case "created_at":
		orderBy = "c.created_at"
	}
	if filters.SortOrder == "asc" {
		orderBy += " ASC NULLS LAST"
	} else {
		orderBy += " DESC NULLS LAST"
	}

	// Apply pagination
	limit := filters.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT c.id, c.name, c.email, c.phone, c.status, c.source_type, c.created_at,
		       p.name as prodi_name, u.name as consultant_name, camp.name as campaign_name
		FROM candidates c
		LEFT JOIN prodis p ON p.id = c.prodi_id
		LEFT JOIN users u ON u.id = c.assigned_consultant_id
		LEFT JOIN campaigns camp ON camp.id = c.campaign_id
		%s
		ORDER BY %s
		LIMIT %d OFFSET %d
	`, whereClause, orderBy, limit, offset)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list candidates: %w", err)
	}
	defer rows.Close()

	candidates := []CandidateListItem{}
	for rows.Next() {
		var c CandidateListItem
		err := rows.Scan(
			&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &c.SourceType, &c.CreatedAt,
			&c.ProdiName, &c.ConsultantName, &c.CampaignName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan candidate: %w", err)
		}
		candidates = append(candidates, c)
	}

	return &CandidateListResult{
		Candidates: candidates,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	}, nil
}

// CandidateStatusStats holds counts by status
type CandidateStatusStats struct {
	Total       int `json:"total"`
	Registered  int `json:"registered"`
	Prospecting int `json:"prospecting"`
	Committed   int `json:"committed"`
	Enrolled    int `json:"enrolled"`
	Lost        int `json:"lost"`
}

// GetCandidateStatusStats returns counts for each status
// If visibilityConsultantID is set, only counts that consultant's candidates
// If visibilitySupervisorID is set, counts team's candidates
func GetCandidateStatusStats(ctx context.Context, visibilityConsultantID, visibilitySupervisorID *string) (*CandidateStatusStats, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	// Role-based visibility
	if visibilityConsultantID != nil && *visibilityConsultantID != "" {
		whereClause += fmt.Sprintf(" AND assigned_consultant_id = $%d", argNum)
		args = append(args, *visibilityConsultantID)
		argNum++
	} else if visibilitySupervisorID != nil && *visibilitySupervisorID != "" {
		whereClause += fmt.Sprintf(" AND assigned_consultant_id IN (SELECT id FROM users WHERE id_supervisor = $%d OR id = $%d)", argNum, argNum+1)
		args = append(args, *visibilitySupervisorID, *visibilitySupervisorID)
		argNum += 2
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'registered') as registered,
			COUNT(*) FILTER (WHERE status = 'prospecting') as prospecting,
			COUNT(*) FILTER (WHERE status = 'committed') as committed,
			COUNT(*) FILTER (WHERE status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE status = 'lost') as lost
		FROM candidates
		%s
	`, whereClause)

	var stats CandidateStatusStats
	err := pool.QueryRow(ctx, query, args...).Scan(
		&stats.Total, &stats.Registered, &stats.Prospecting, &stats.Committed, &stats.Enrolled, &stats.Lost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate stats: %w", err)
	}
	return &stats, nil
}

// GetNextConsultantForAssignment gets the next consultant using the active algorithm
func GetNextConsultantForAssignment(ctx context.Context) (*string, error) {
	// Get active algorithm
	algo, err := FindActiveAssignmentAlgorithm(ctx)
	if err != nil {
		return nil, err
	}
	if algo == nil {
		return nil, nil // No active algorithm
	}

	var consultantID string

	switch algo.Code {
	case "round_robin":
		// Get the consultant with the oldest last assignment
		err = pool.QueryRow(ctx, `
			SELECT u.id FROM users u
			WHERE u.role = 'consultant' AND u.is_active = true
			ORDER BY (
				SELECT MAX(c.created_at) FROM candidates c WHERE c.assigned_consultant_id = u.id
			) NULLS FIRST, u.created_at
			LIMIT 1
		`).Scan(&consultantID)

	case "load_balanced":
		// Get the consultant with the fewest active candidates
		err = pool.QueryRow(ctx, `
			SELECT u.id FROM users u
			LEFT JOIN candidates c ON c.assigned_consultant_id = u.id AND c.status NOT IN ('enrolled', 'lost')
			WHERE u.role = 'consultant' AND u.is_active = true
			GROUP BY u.id
			ORDER BY COUNT(c.id), u.created_at
			LIMIT 1
		`).Scan(&consultantID)

	default:
		// Default to round robin
		err = pool.QueryRow(ctx, `
			SELECT u.id FROM users u
			WHERE u.role = 'consultant' AND u.is_active = true
			ORDER BY (
				SELECT MAX(c.created_at) FROM candidates c WHERE c.assigned_consultant_id = u.id
			) NULLS FIRST, u.created_at
			LIMIT 1
		`).Scan(&consultantID)
	}

	if err == pgx.ErrNoRows {
		return nil, nil // No consultant available
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get next consultant: %w", err)
	}
	return &consultantID, nil
}
