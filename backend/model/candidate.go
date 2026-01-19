package model

import (
	"context"
	"fmt"
	"log/slog"
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
	AssignedConsultantID  *string    `json:"id_assigned_consultant,omitempty"`
	Status                string     `json:"status"`
	LostReasonID          *string    `json:"id_lost_reason,omitempty"`
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
	var emailEnc, phoneEnc *string

	// Encrypt email and phone before storing
	if email != "" {
		enc, err := encryptEmail(email)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt email: %w", err)
		}
		emailEnc = &enc
	}
	if phone != "" {
		enc, err := encryptPhone(phone)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt phone: %w", err)
		}
		phoneEnc = &enc
	}

	err := pool.QueryRow(ctx, `
		INSERT INTO candidates (email, phone, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		          high_school, graduation_year, id_prodi, id_campaign, id_referrer, id_referred_by_candidate,
		          source_type, source_detail, id_assigned_consultant, status, id_lost_reason, lost_at, created_at, updated_at
	`, emailEnc, phoneEnc, passwordHash).Scan(
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

	// Decrypt fields before returning
	if err := decryptCandidateFields(&candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
	}

	return &candidate, nil
}

// FindCandidateByID finds a candidate by ID
func FindCandidateByID(ctx context.Context, id string) (*Candidate, error) {
	var candidate Candidate
	err := pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, id_prodi, id_campaign, id_referrer, id_referred_by_candidate,
		       source_type, source_detail, id_assigned_consultant, status, id_lost_reason, lost_at, created_at, updated_at
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

	// Decrypt fields
	if err := decryptCandidateFields(&candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
	}

	return &candidate, nil
}

// FindCandidateByEmail finds a candidate by email
func FindCandidateByEmail(ctx context.Context, email string) (*Candidate, error) {
	// Encrypt email for search (deterministic encryption allows equality match)
	emailEnc, err := encryptEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt email for search: %w", err)
	}

	var candidate Candidate
	err = pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, id_prodi, id_campaign, id_referrer, id_referred_by_candidate,
		       source_type, source_detail, id_assigned_consultant, status, id_lost_reason, lost_at, created_at, updated_at
		FROM candidates WHERE email = $1
	`, emailEnc).Scan(
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

	// Decrypt fields
	if err := decryptCandidateFields(&candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
	}

	return &candidate, nil
}

// FindCandidateByPhone finds a candidate by phone
func FindCandidateByPhone(ctx context.Context, phone string) (*Candidate, error) {
	// Encrypt phone for search (deterministic encryption allows equality match)
	phoneEnc, err := encryptPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt phone for search: %w", err)
	}

	var candidate Candidate
	err = pool.QueryRow(ctx, `
		SELECT id, email, email_verified, phone, phone_verified, password_hash, name, address, city, province,
		       high_school, graduation_year, id_prodi, id_campaign, id_referrer, id_referred_by_candidate,
		       source_type, source_detail, id_assigned_consultant, status, id_lost_reason, lost_at, created_at, updated_at
		FROM candidates WHERE phone = $1
	`, phoneEnc).Scan(
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

	// Decrypt fields
	if err := decryptCandidateFields(&candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
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
	// Encrypt fields before storing
	nameEnc, err := encryptName(name)
	if err != nil {
		return fmt.Errorf("failed to encrypt name: %w", err)
	}

	var addressEnc, cityEnc, provinceEnc *string
	if address != "" {
		enc, err := encryptNullableP(&address)
		if err != nil {
			return fmt.Errorf("failed to encrypt address: %w", err)
		}
		addressEnc = enc
	}
	if city != "" {
		enc, err := encryptNullableP(&city)
		if err != nil {
			return fmt.Errorf("failed to encrypt city: %w", err)
		}
		cityEnc = enc
	}
	if province != "" {
		enc, err := encryptNullableP(&province)
		if err != nil {
			return fmt.Errorf("failed to encrypt province: %w", err)
		}
		provinceEnc = enc
	}

	_, err = pool.Exec(ctx, `
		UPDATE candidates
		SET name = $2, address = $3, city = $4, province = $5, updated_at = NOW()
		WHERE id = $1
	`, id, nameEnc, addressEnc, cityEnc, provinceEnc)
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

	// Encrypt high_school before storing
	var highSchoolEnc *string
	if highSchool != "" {
		enc, err := encryptNullableP(&highSchool)
		if err != nil {
			return fmt.Errorf("failed to encrypt high_school: %w", err)
		}
		highSchoolEnc = enc
	}

	_, err := pool.Exec(ctx, `
		UPDATE candidates
		SET high_school = $2, graduation_year = $3, id_prodi = $4, updated_at = NOW()
		WHERE id = $1
	`, id, highSchoolEnc, graduationYear, prodiPtr)
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
		SET source_type = $2, source_detail = $3, id_campaign = $4, id_referrer = $5, id_referred_by_candidate = $6, updated_at = NOW()
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
		UPDATE candidates SET id_assigned_consultant = $2, updated_at = NOW() WHERE id = $1
	`, id, consultantID)
	if err != nil {
		return fmt.Errorf("failed to assign consultant: %w", err)
	}
	return nil
}

// UpdateCandidateStatus updates the candidate status and triggers commission creation if applicable
func UpdateCandidateStatus(ctx context.Context, id, status string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET status = $2, updated_at = NOW() WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Trigger commission creation based on status
	var triggerEvent string
	switch status {
	case "prospecting":
		triggerEvent = "registration"
	case "committed":
		triggerEvent = "commitment"
	case "enrolled":
		triggerEvent = "enrollment"
	}

	if triggerEvent != "" {
		// Create commission if candidate has a referrer
		if err := CreateCommissionForCandidate(ctx, id, triggerEvent); err != nil {
			// Log error but don't fail the status update
			// Commission creation is secondary to status change
			return nil
		}
	}

	return nil
}

// MarkCandidateLost marks a candidate as lost with reason
func MarkCandidateLost(ctx context.Context, id, lostReasonID string) error {
	_, err := pool.Exec(ctx, `
		UPDATE candidates SET status = 'lost', id_lost_reason = $2, lost_at = NOW(), updated_at = NOW() WHERE id = $1
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
		       c.high_school, c.graduation_year, c.id_prodi, c.id_campaign, c.id_referrer, c.id_referred_by_candidate,
		       c.source_type, c.source_detail, c.id_assigned_consultant, c.status, c.id_lost_reason, c.lost_at, c.created_at, c.updated_at,
		       p.name as prodi_name,
		       u.name as consultant_name,
		       camp.name as campaign_name,
		       ref.name as referrer_name,
		       lr.name as lost_reason_name
		FROM candidates c
		LEFT JOIN prodis p ON p.id = c.id_prodi
		LEFT JOIN users u ON u.id = c.id_assigned_consultant
		LEFT JOIN campaigns camp ON camp.id = c.id_campaign
		LEFT JOIN referrers ref ON ref.id = c.id_referrer
		LEFT JOIN lost_reasons lr ON lr.id = c.id_lost_reason
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

	// Decrypt candidate fields
	if err := decryptCandidateFields(&data.Candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
	}

	// Decrypt related encrypted fields (consultant name, referrer name)
	data.ConsultantName, _ = decryptNullableP(data.ConsultantName)
	data.ReferrerName, _ = decryptNullableP(data.ReferrerName)

	return &data, nil
}

// GetCandidateDashboardData gets candidate with related details for dashboard
func GetCandidateDashboardData(ctx context.Context, id string) (*CandidateDashboardData, error) {
	var data CandidateDashboardData
	err := pool.QueryRow(ctx, `
		SELECT c.id, c.email, c.email_verified, c.phone, c.phone_verified, c.password_hash, c.name, c.address, c.city, c.province,
		       c.high_school, c.graduation_year, c.id_prodi, c.id_campaign, c.id_referrer, c.id_referred_by_candidate,
		       c.source_type, c.source_detail, c.id_assigned_consultant, c.status, c.id_lost_reason, c.lost_at, c.created_at, c.updated_at,
		       p.name as prodi_name, u.name as consultant_name, u.email as consultant_email
		FROM candidates c
		LEFT JOIN prodis p ON p.id = c.id_prodi
		LEFT JOIN users u ON u.id = c.id_assigned_consultant
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

	// Decrypt candidate fields
	if err := decryptCandidateFields(&data.Candidate); err != nil {
		return nil, fmt.Errorf("failed to decrypt candidate: %w", err)
	}

	// Decrypt related encrypted fields (consultant name and email)
	data.ConsultantName, _ = decryptNullableP(data.ConsultantName)
	data.ConsultantEmail, _ = decryptNullableD(data.ConsultantEmail)

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
		whereClause += fmt.Sprintf(" AND c.id_assigned_consultant = $%d", argNum)
		args = append(args, *visibilityConsultantID)
		argNum++
	} else if visibilitySupervisorID != nil && *visibilitySupervisorID != "" {
		// Supervisor sees their team's candidates
		whereClause += fmt.Sprintf(" AND c.id_assigned_consultant IN (SELECT id FROM users WHERE id_supervisor = $%d OR id = $%d)", argNum, argNum+1)
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
		whereClause += fmt.Sprintf(" AND c.id_assigned_consultant = $%d", argNum)
		args = append(args, filters.ConsultantID)
		argNum++
	}

	if filters.ProdiID != "" {
		whereClause += fmt.Sprintf(" AND c.id_prodi = $%d", argNum)
		args = append(args, filters.ProdiID)
		argNum++
	}

	if filters.CampaignID != "" {
		whereClause += fmt.Sprintf(" AND c.id_campaign = $%d", argNum)
		args = append(args, filters.CampaignID)
		argNum++
	}

	if filters.SourceType != "" {
		whereClause += fmt.Sprintf(" AND c.source_type = $%d", argNum)
		args = append(args, filters.SourceType)
		argNum++
	}

	if filters.Search != "" {
		// With encryption, we can only do exact match on deterministically encrypted fields (email, phone)
		// Search by encrypting the search term and matching exactly
		emailEnc, err := encryptEmail(filters.Search)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt search term: %w", err)
		}
		phoneEnc, err := encryptPhone(filters.Search)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt search term: %w", err)
		}
		whereClause += fmt.Sprintf(" AND (c.email = $%d OR c.phone = $%d)", argNum, argNum+1)
		args = append(args, emailEnc, phoneEnc)
		argNum += 2
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
		LEFT JOIN prodis p ON p.id = c.id_prodi
		LEFT JOIN users u ON u.id = c.id_assigned_consultant
		LEFT JOIN campaigns camp ON camp.id = c.id_campaign
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

		// Decrypt encrypted fields
		c.Name, _ = decryptNullableP(c.Name)
		c.Email, _ = decryptNullableD(c.Email)
		c.Phone, _ = decryptNullableD(c.Phone)
		c.ConsultantName, _ = decryptNullableP(c.ConsultantName)

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
		whereClause += fmt.Sprintf(" AND id_assigned_consultant = $%d", argNum)
		args = append(args, *visibilityConsultantID)
		argNum++
	} else if visibilitySupervisorID != nil && *visibilitySupervisorID != "" {
		whereClause += fmt.Sprintf(" AND id_assigned_consultant IN (SELECT id FROM users WHERE id_supervisor = $%d OR id = $%d)", argNum, argNum+1)
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
				SELECT MAX(c.created_at) FROM candidates c WHERE c.id_assigned_consultant = u.id
			) NULLS FIRST, u.created_at
			LIMIT 1
		`).Scan(&consultantID)

	case "load_balanced":
		// Get the consultant with the fewest active candidates
		err = pool.QueryRow(ctx, `
			SELECT u.id FROM users u
			LEFT JOIN candidates c ON c.id_assigned_consultant = u.id AND c.status NOT IN ('enrolled', 'lost')
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
				SELECT MAX(c.created_at) FROM candidates c WHERE c.id_assigned_consultant = u.id
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

// ConsultantWithWorkload represents a consultant with their candidate counts
type ConsultantWithWorkload struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ActiveCount    int    `json:"active_count"`    // prospecting + committed
	TotalCount     int    `json:"total_count"`     // all assigned
	SupervisorID   *string `json:"supervisor_id,omitempty"`
	SupervisorName *string `json:"supervisor_name,omitempty"`
}

// ListConsultantsWithWorkload returns all active consultants with their candidate counts
func ListConsultantsWithWorkload(ctx context.Context) ([]ConsultantWithWorkload, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			u.id, u.name, u.email, u.id_supervisor,
			s.name as supervisor_name,
			COUNT(c.id) FILTER (WHERE c.status IN ('registered', 'prospecting', 'committed')) as active_count,
			COUNT(c.id) as total_count
		FROM users u
		LEFT JOIN users s ON s.id = u.id_supervisor
		LEFT JOIN candidates c ON c.id_assigned_consultant = u.id
		WHERE u.role = 'consultant' AND u.is_active = true
		GROUP BY u.id, u.name, u.email, u.id_supervisor, s.name
		ORDER BY u.name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list consultants: %w", err)
	}
	defer rows.Close()

	var consultants []ConsultantWithWorkload
	for rows.Next() {
		var c ConsultantWithWorkload
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.SupervisorID, &c.SupervisorName, &c.ActiveCount, &c.TotalCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan consultant: %w", err)
		}

		// Decrypt encrypted fields
		c.Name, _ = decryptName(c.Name)
		if emailDec, err := decryptNullableD(&c.Email); err == nil && emailDec != nil {
			c.Email = *emailDec
		}
		c.SupervisorName, _ = decryptNullableP(c.SupervisorName)

		consultants = append(consultants, c)
	}
	return consultants, nil
}

// ReassignCandidate reassigns a candidate to a different consultant
func ReassignCandidate(ctx context.Context, candidateID, newConsultantID, reassignedBy string) error {
	// Get current consultant for logging
	var oldConsultantID *string
	err := pool.QueryRow(ctx, `SELECT id_assigned_consultant FROM candidates WHERE id = $1`, candidateID).Scan(&oldConsultantID)
	if err != nil {
		return fmt.Errorf("failed to get current consultant: %w", err)
	}

	// Update assignment
	_, err = pool.Exec(ctx, `
		UPDATE candidates SET id_assigned_consultant = $2, updated_at = NOW() WHERE id = $1
	`, candidateID, newConsultantID)
	if err != nil {
		return fmt.Errorf("failed to reassign candidate: %w", err)
	}

	// Log the reassignment as an interaction
	oldID := ""
	if oldConsultantID != nil {
		oldID = *oldConsultantID
	}
	remarks := fmt.Sprintf("Kandidat dipindahkan dari konsultan sebelumnya")
	if oldID == "" {
		remarks = "Konsultan pertama kali ditugaskan"
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO interactions (id_candidate, id_consultant, channel, remarks)
		VALUES ($1, $2, 'system', $3)
	`, candidateID, reassignedBy, remarks)
	if err != nil {
		// Log error but don't fail the reassignment
		slog.Error("failed to log reassignment interaction", "error", err)
	}

	return nil
}
