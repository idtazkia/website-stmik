package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Announcement represents a system announcement
type Announcement struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	TargetStatus  *string    `json:"target_status,omitempty"`
	TargetProdiID *string    `json:"target_prodi_id,omitempty"`
	IsPublished   bool       `json:"is_published"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AnnouncementWithDetails includes related names for display
type AnnouncementWithDetails struct {
	Announcement
	TargetProdiName *string `json:"target_prodi_name,omitempty"`
	CreatedByName   *string `json:"created_by_name,omitempty"`
	ReadCount       int     `json:"read_count"`
}

// AnnouncementForCandidate includes read status for portal display
type AnnouncementForCandidate struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
}

// CreateAnnouncement creates a new announcement
func CreateAnnouncement(ctx context.Context, title, content string, targetStatus, targetProdiID, createdBy *string) (*Announcement, error) {
	var ann Announcement
	err := pool.QueryRow(ctx, `
		INSERT INTO announcements (title, content, target_status, target_prodi_id, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, target_status, target_prodi_id, is_published, published_at, created_by, created_at, updated_at
	`, title, content, targetStatus, targetProdiID, createdBy).Scan(
		&ann.ID, &ann.Title, &ann.Content, &ann.TargetStatus, &ann.TargetProdiID,
		&ann.IsPublished, &ann.PublishedAt, &ann.CreatedBy, &ann.CreatedAt, &ann.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}
	return &ann, nil
}

// FindAnnouncementByID finds an announcement by ID
func FindAnnouncementByID(ctx context.Context, id string) (*Announcement, error) {
	var ann Announcement
	err := pool.QueryRow(ctx, `
		SELECT id, title, content, target_status, target_prodi_id, is_published, published_at, created_by, created_at, updated_at
		FROM announcements WHERE id = $1
	`, id).Scan(
		&ann.ID, &ann.Title, &ann.Content, &ann.TargetStatus, &ann.TargetProdiID,
		&ann.IsPublished, &ann.PublishedAt, &ann.CreatedBy, &ann.CreatedAt, &ann.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find announcement: %w", err)
	}
	return &ann, nil
}

// ListAnnouncements lists all announcements for admin
func ListAnnouncements(ctx context.Context) ([]AnnouncementWithDetails, error) {
	rows, err := pool.Query(ctx, `
		SELECT a.id, a.title, a.content, a.target_status, a.target_prodi_id, a.is_published, a.published_at, a.created_by, a.created_at, a.updated_at,
		       p.name as target_prodi_name, u.name as created_by_name,
		       (SELECT COUNT(*) FROM announcement_reads ar WHERE ar.announcement_id = a.id) as read_count
		FROM announcements a
		LEFT JOIN prodis p ON p.id = a.target_prodi_id
		LEFT JOIN users u ON u.id = a.created_by
		ORDER BY a.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list announcements: %w", err)
	}
	defer rows.Close()

	var announcements []AnnouncementWithDetails
	for rows.Next() {
		var ann AnnouncementWithDetails
		err := rows.Scan(
			&ann.ID, &ann.Title, &ann.Content, &ann.TargetStatus, &ann.TargetProdiID,
			&ann.IsPublished, &ann.PublishedAt, &ann.CreatedBy, &ann.CreatedAt, &ann.UpdatedAt,
			&ann.TargetProdiName, &ann.CreatedByName, &ann.ReadCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan announcement: %w", err)
		}
		announcements = append(announcements, ann)
	}
	return announcements, nil
}

// ListAnnouncementsForCandidate lists published announcements for a candidate
func ListAnnouncementsForCandidate(ctx context.Context, candidateID, status string, prodiID *string) ([]AnnouncementForCandidate, error) {
	rows, err := pool.Query(ctx, `
		SELECT a.id, a.title, a.content, a.published_at,
		       ar.read_at IS NOT NULL as is_read, ar.read_at
		FROM announcements a
		LEFT JOIN announcement_reads ar ON ar.announcement_id = a.id AND ar.candidate_id = $1
		WHERE a.is_published = true
		  AND (a.target_status IS NULL OR a.target_status = $2)
		  AND (a.target_prodi_id IS NULL OR a.target_prodi_id = $3)
		ORDER BY a.published_at DESC
	`, candidateID, status, prodiID)
	if err != nil {
		return nil, fmt.Errorf("failed to list announcements for candidate: %w", err)
	}
	defer rows.Close()

	var announcements []AnnouncementForCandidate
	for rows.Next() {
		var ann AnnouncementForCandidate
		err := rows.Scan(&ann.ID, &ann.Title, &ann.Content, &ann.PublishedAt, &ann.IsRead, &ann.ReadAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan announcement: %w", err)
		}
		announcements = append(announcements, ann)
	}
	return announcements, nil
}

// CountUnreadAnnouncementsForCandidate counts unread announcements for a candidate
func CountUnreadAnnouncementsForCandidate(ctx context.Context, candidateID, status string, prodiID *string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM announcements a
		WHERE a.is_published = true
		  AND (a.target_status IS NULL OR a.target_status = $2)
		  AND (a.target_prodi_id IS NULL OR a.target_prodi_id = $3)
		  AND NOT EXISTS (SELECT 1 FROM announcement_reads ar WHERE ar.announcement_id = a.id AND ar.candidate_id = $1)
	`, candidateID, status, prodiID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread announcements: %w", err)
	}
	return count, nil
}

// UpdateAnnouncement updates an announcement
func UpdateAnnouncement(ctx context.Context, id, title, content string, targetStatus, targetProdiID *string) error {
	_, err := pool.Exec(ctx, `
		UPDATE announcements
		SET title = $2, content = $3, target_status = $4, target_prodi_id = $5, updated_at = NOW()
		WHERE id = $1
	`, id, title, content, targetStatus, targetProdiID)
	if err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}
	return nil
}

// PublishAnnouncement publishes an announcement
func PublishAnnouncement(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE announcements
		SET is_published = true, published_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND is_published = false
	`, id)
	if err != nil {
		return fmt.Errorf("failed to publish announcement: %w", err)
	}
	return nil
}

// UnpublishAnnouncement unpublishes an announcement
func UnpublishAnnouncement(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE announcements
		SET is_published = false, updated_at = NOW()
		WHERE id = $1 AND is_published = true
	`, id)
	if err != nil {
		return fmt.Errorf("failed to unpublish announcement: %w", err)
	}
	return nil
}

// DeleteAnnouncement deletes an announcement
func DeleteAnnouncement(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM announcements WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}
	return nil
}

// MarkAnnouncementRead marks an announcement as read by a candidate
func MarkAnnouncementRead(ctx context.Context, announcementID, candidateID string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO announcement_reads (announcement_id, candidate_id)
		VALUES ($1, $2)
		ON CONFLICT (announcement_id, candidate_id) DO NOTHING
	`, announcementID, candidateID)
	if err != nil {
		return fmt.Errorf("failed to mark announcement as read: %w", err)
	}
	return nil
}
