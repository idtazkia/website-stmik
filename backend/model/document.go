package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Document represents an uploaded document
type Document struct {
	ID              string     `json:"id"`
	CandidateID     string     `json:"candidate_id"`
	DocumentTypeID  string     `json:"document_type_id"`
	FileName        string     `json:"file_name"`
	FilePath        string     `json:"file_path"`
	FileSize        int        `json:"file_size"`
	MimeType        string     `json:"mime_type"`
	Status          string     `json:"status"` // pending, approved, rejected
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	ReviewedBy      *string    `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	IsDeferred      bool       `json:"is_deferred"`
	DeferredReason  *string    `json:"deferred_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DocumentWithType includes document type info
type DocumentWithType struct {
	Document
	TypeName        string  `json:"type_name"`
	TypeCode        string  `json:"type_code"`
	TypeDescription *string `json:"type_description,omitempty"`
	IsRequired      bool    `json:"is_required"`
	CanDefer        bool    `json:"can_defer"`
}

// CreateDocument creates or replaces a document for a candidate
func CreateDocument(ctx context.Context, candidateID, documentTypeID, fileName, filePath string, fileSize int, mimeType string) (*Document, error) {
	var doc Document
	err := pool.QueryRow(ctx, `
		INSERT INTO documents (candidate_id, document_type_id, file_name, file_path, file_size, mime_type, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
		ON CONFLICT (candidate_id, document_type_id) DO UPDATE SET
			file_name = EXCLUDED.file_name,
			file_path = EXCLUDED.file_path,
			file_size = EXCLUDED.file_size,
			mime_type = EXCLUDED.mime_type,
			status = 'pending',
			rejection_reason = NULL,
			reviewed_by = NULL,
			reviewed_at = NULL,
			is_deferred = false,
			deferred_reason = NULL,
			updated_at = NOW()
		RETURNING id, candidate_id, document_type_id, file_name, file_path, file_size, mime_type,
		          status, rejection_reason, reviewed_by, reviewed_at, is_deferred, deferred_reason,
		          created_at, updated_at
	`, candidateID, documentTypeID, fileName, filePath, fileSize, mimeType).Scan(
		&doc.ID, &doc.CandidateID, &doc.DocumentTypeID, &doc.FileName, &doc.FilePath,
		&doc.FileSize, &doc.MimeType, &doc.Status, &doc.RejectionReason, &doc.ReviewedBy,
		&doc.ReviewedAt, &doc.IsDeferred, &doc.DeferredReason, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	return &doc, nil
}

// ListDocumentsByCandidate returns all documents for a candidate
func ListDocumentsByCandidate(ctx context.Context, candidateID string) ([]DocumentWithType, error) {
	rows, err := pool.Query(ctx, `
		SELECT d.id, d.candidate_id, d.document_type_id, d.file_name, d.file_path, d.file_size,
		       d.mime_type, d.status, d.rejection_reason, d.reviewed_by, d.reviewed_at,
		       d.is_deferred, d.deferred_reason, d.created_at, d.updated_at,
		       dt.name, dt.code, dt.description, dt.is_required, dt.can_defer
		FROM documents d
		JOIN document_types dt ON dt.id = d.document_type_id
		WHERE d.candidate_id = $1
		ORDER BY dt.display_order
	`, candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var docs []DocumentWithType
	for rows.Next() {
		var d DocumentWithType
		err := rows.Scan(
			&d.ID, &d.CandidateID, &d.DocumentTypeID, &d.FileName, &d.FilePath, &d.FileSize,
			&d.MimeType, &d.Status, &d.RejectionReason, &d.ReviewedBy, &d.ReviewedAt,
			&d.IsDeferred, &d.DeferredReason, &d.CreatedAt, &d.UpdatedAt,
			&d.TypeName, &d.TypeCode, &d.TypeDescription, &d.IsRequired, &d.CanDefer,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// FindDocumentByID finds a document by ID
func FindDocumentByID(ctx context.Context, id string) (*DocumentWithType, error) {
	var d DocumentWithType
	err := pool.QueryRow(ctx, `
		SELECT d.id, d.candidate_id, d.document_type_id, d.file_name, d.file_path, d.file_size,
		       d.mime_type, d.status, d.rejection_reason, d.reviewed_by, d.reviewed_at,
		       d.is_deferred, d.deferred_reason, d.created_at, d.updated_at,
		       dt.name, dt.code, dt.description, dt.is_required, dt.can_defer
		FROM documents d
		JOIN document_types dt ON dt.id = d.document_type_id
		WHERE d.id = $1
	`, id).Scan(
		&d.ID, &d.CandidateID, &d.DocumentTypeID, &d.FileName, &d.FilePath, &d.FileSize,
		&d.MimeType, &d.Status, &d.RejectionReason, &d.ReviewedBy, &d.ReviewedAt,
		&d.IsDeferred, &d.DeferredReason, &d.CreatedAt, &d.UpdatedAt,
		&d.TypeName, &d.TypeCode, &d.TypeDescription, &d.IsRequired, &d.CanDefer,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return &d, nil
}

// ApproveDocument approves a document
func ApproveDocument(ctx context.Context, documentID, reviewerID string) error {
	_, err := pool.Exec(ctx, `
		UPDATE documents
		SET status = 'approved', reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, documentID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to approve document: %w", err)
	}
	return nil
}

// RejectDocument rejects a document with reason
func RejectDocument(ctx context.Context, documentID, reviewerID, reason string) error {
	_, err := pool.Exec(ctx, `
		UPDATE documents
		SET status = 'rejected', rejection_reason = $2, reviewed_by = $3, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, documentID, reason, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to reject document: %w", err)
	}
	return nil
}

// DeferDocument marks a document as deferred
func DeferDocument(ctx context.Context, candidateID, documentTypeID, reason string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO documents (candidate_id, document_type_id, file_name, file_path, file_size, mime_type, status, is_deferred, deferred_reason)
		VALUES ($1, $2, '', '', 0, '', 'pending', true, $3)
		ON CONFLICT (candidate_id, document_type_id) DO UPDATE SET
			is_deferred = true,
			deferred_reason = $3,
			updated_at = NOW()
	`, candidateID, documentTypeID, reason)
	if err != nil {
		return fmt.Errorf("failed to defer document: %w", err)
	}
	return nil
}

// GetDocumentStats returns document completion stats for a candidate
func GetDocumentStats(ctx context.Context, candidateID string) (uploaded, approved, total int, err error) {
	// Get total required documents
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM document_types WHERE is_active = true AND is_required = true
	`).Scan(&total)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count document types: %w", err)
	}

	// Get uploaded and approved counts
	err = pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status IN ('pending', 'approved') OR is_deferred = true) as uploaded,
			COUNT(*) FILTER (WHERE status = 'approved') as approved
		FROM documents d
		JOIN document_types dt ON dt.id = d.document_type_id
		WHERE d.candidate_id = $1 AND dt.is_required = true
	`, candidateID).Scan(&uploaded, &approved)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return uploaded, approved, total, nil
}

// DocumentReviewFilters for filtering documents in admin review
type DocumentReviewFilters struct {
	Status string
	Type   string
	Search string
}

// DocumentForReview includes document with candidate info for admin review
type DocumentForReview struct {
	ID              string
	CandidateID     string
	CandidateName   string
	ProdiName       string
	TypeCode        string
	TypeName        string
	FileName        string
	FilePath        string
	FileSize        int
	MimeType        string
	Status          string
	RejectionReason *string
	ReviewedBy      *string
	ReviewedAt      *time.Time
	CreatedAt       time.Time
}

// ListDocumentsForReview returns documents for admin review with filters
func ListDocumentsForReview(ctx context.Context, filters DocumentReviewFilters) ([]DocumentForReview, error) {
	query := `
		SELECT d.id, d.candidate_id, c.name, p.name, dt.code, dt.name,
		       d.file_name, d.file_path, d.file_size, d.mime_type,
		       d.status, d.rejection_reason, d.reviewed_by, d.reviewed_at, d.created_at
		FROM documents d
		JOIN document_types dt ON dt.id = d.document_type_id
		JOIN candidates c ON c.id = d.candidate_id
		LEFT JOIN prodis p ON p.id = c.prodi_id
		WHERE d.file_name != '' AND d.is_deferred = false
	`
	args := []interface{}{}
	argNum := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND d.status = $%d", argNum)
		args = append(args, filters.Status)
		argNum++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND dt.code = $%d", argNum)
		args = append(args, filters.Type)
		argNum++
	}

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.email ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}

	query += " ORDER BY d.created_at DESC LIMIT 100"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents for review: %w", err)
	}
	defer rows.Close()

	var docs []DocumentForReview
	for rows.Next() {
		var d DocumentForReview
		var candidateName, prodiName *string
		err := rows.Scan(
			&d.ID, &d.CandidateID, &candidateName, &prodiName, &d.TypeCode, &d.TypeName,
			&d.FileName, &d.FilePath, &d.FileSize, &d.MimeType,
			&d.Status, &d.RejectionReason, &d.ReviewedBy, &d.ReviewedAt, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		if candidateName != nil {
			d.CandidateName = *candidateName
		}
		if prodiName != nil {
			d.ProdiName = *prodiName
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// DocumentReviewStats for admin dashboard
type DocumentReviewStats struct {
	Pending       int
	ApprovedToday int
	RejectedToday int
	Total         int
}

// GetDocumentReviewStats returns stats for admin document review dashboard
func GetDocumentReviewStats(ctx context.Context) (*DocumentReviewStats, error) {
	var stats DocumentReviewStats
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending' AND file_name != '' AND is_deferred = false) as pending,
			COUNT(*) FILTER (WHERE status = 'approved' AND reviewed_at >= CURRENT_DATE) as approved_today,
			COUNT(*) FILTER (WHERE status = 'rejected' AND reviewed_at >= CURRENT_DATE) as rejected_today,
			COUNT(*) FILTER (WHERE file_name != '' AND is_deferred = false) as total
		FROM documents
	`).Scan(&stats.Pending, &stats.ApprovedToday, &stats.RejectedToday, &stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get document review stats: %w", err)
	}
	return &stats, nil
}
