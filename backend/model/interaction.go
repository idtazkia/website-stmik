package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Interaction represents a logged interaction with a candidate
type Interaction struct {
	ID                   string     `json:"id"`
	CandidateID          string     `json:"candidate_id"`
	ConsultantID         string     `json:"consultant_id"`
	Channel              string     `json:"channel"` // call, whatsapp, email, campus_visit, home_visit
	CategoryID           *string    `json:"category_id,omitempty"`
	ObstacleID           *string    `json:"obstacle_id,omitempty"`
	Remarks              string     `json:"remarks"`
	NextFollowupDate     *time.Time `json:"next_followup_date,omitempty"`
	NextAction           *string    `json:"next_action,omitempty"`
	SupervisorSuggestion *string    `json:"supervisor_suggestion,omitempty"`
	SuggestionReadAt     *time.Time `json:"suggestion_read_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
}

// InteractionWithDetails includes related names for display
type InteractionWithDetails struct {
	Interaction
	ConsultantName    string  `json:"consultant_name"`
	CategoryName      *string `json:"category_name,omitempty"`
	CategorySentiment *string `json:"category_sentiment,omitempty"`
	ObstacleName      *string `json:"obstacle_name,omitempty"`
}

// CreateInteraction creates a new interaction and updates candidate's last_interaction_at
func CreateInteraction(ctx context.Context, candidateID, consultantID, channel string, categoryID, obstacleID *string, remarks string, nextFollowupDate *time.Time, nextAction *string) (*Interaction, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var interaction Interaction
	err = tx.QueryRow(ctx, `
		INSERT INTO interactions (candidate_id, consultant_id, channel, category_id, obstacle_id, remarks, next_followup_date, next_action)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, candidate_id, consultant_id, channel, category_id, obstacle_id, remarks, next_followup_date, next_action, supervisor_suggestion, suggestion_read_at, created_at
	`, candidateID, consultantID, channel, categoryID, obstacleID, remarks, nextFollowupDate, nextAction).Scan(
		&interaction.ID, &interaction.CandidateID, &interaction.ConsultantID, &interaction.Channel,
		&interaction.CategoryID, &interaction.ObstacleID, &interaction.Remarks,
		&interaction.NextFollowupDate, &interaction.NextAction, &interaction.SupervisorSuggestion,
		&interaction.SuggestionReadAt, &interaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create interaction: %w", err)
	}

	// Update candidate's last_interaction_at
	_, err = tx.Exec(ctx, `
		UPDATE candidates SET last_interaction_at = $1, updated_at = NOW() WHERE id = $2
	`, interaction.CreatedAt, candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to update candidate last_interaction_at: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &interaction, nil
}

// FindInteractionByID finds an interaction by ID
func FindInteractionByID(ctx context.Context, id string) (*Interaction, error) {
	var interaction Interaction
	err := pool.QueryRow(ctx, `
		SELECT id, candidate_id, consultant_id, channel, category_id, obstacle_id, remarks,
		       next_followup_date, next_action, supervisor_suggestion, suggestion_read_at, created_at
		FROM interactions WHERE id = $1
	`, id).Scan(
		&interaction.ID, &interaction.CandidateID, &interaction.ConsultantID, &interaction.Channel,
		&interaction.CategoryID, &interaction.ObstacleID, &interaction.Remarks,
		&interaction.NextFollowupDate, &interaction.NextAction, &interaction.SupervisorSuggestion,
		&interaction.SuggestionReadAt, &interaction.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find interaction: %w", err)
	}
	return &interaction, nil
}

// ListInteractionsByCandidate returns all interactions for a candidate with details
func ListInteractionsByCandidate(ctx context.Context, candidateID string) ([]InteractionWithDetails, error) {
	rows, err := pool.Query(ctx, `
		SELECT i.id, i.candidate_id, i.consultant_id, i.channel, i.category_id, i.obstacle_id, i.remarks,
		       i.next_followup_date, i.next_action, i.supervisor_suggestion, i.suggestion_read_at, i.created_at,
		       u.name as consultant_name,
		       ic.name as category_name, ic.sentiment as category_sentiment,
		       o.name as obstacle_name
		FROM interactions i
		LEFT JOIN users u ON u.id = i.consultant_id
		LEFT JOIN interaction_categories ic ON ic.id = i.category_id
		LEFT JOIN obstacles o ON o.id = i.obstacle_id
		WHERE i.candidate_id = $1
		ORDER BY i.created_at DESC
	`, candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list interactions: %w", err)
	}
	defer rows.Close()

	var interactions []InteractionWithDetails
	for rows.Next() {
		var i InteractionWithDetails
		err := rows.Scan(
			&i.ID, &i.CandidateID, &i.ConsultantID, &i.Channel, &i.CategoryID, &i.ObstacleID, &i.Remarks,
			&i.NextFollowupDate, &i.NextAction, &i.SupervisorSuggestion, &i.SuggestionReadAt, &i.CreatedAt,
			&i.ConsultantName, &i.CategoryName, &i.CategorySentiment, &i.ObstacleName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interaction: %w", err)
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

// ListInteractionsByConsultant returns all interactions logged by a consultant
func ListInteractionsByConsultant(ctx context.Context, consultantID string, limit, offset int) ([]InteractionWithDetails, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := pool.Query(ctx, `
		SELECT i.id, i.candidate_id, i.consultant_id, i.channel, i.category_id, i.obstacle_id, i.remarks,
		       i.next_followup_date, i.next_action, i.supervisor_suggestion, i.suggestion_read_at, i.created_at,
		       u.name as consultant_name,
		       ic.name as category_name, ic.sentiment as category_sentiment,
		       o.name as obstacle_name
		FROM interactions i
		LEFT JOIN users u ON u.id = i.consultant_id
		LEFT JOIN interaction_categories ic ON ic.id = i.category_id
		LEFT JOIN obstacles o ON o.id = i.obstacle_id
		WHERE i.consultant_id = $1
		ORDER BY i.created_at DESC
		LIMIT $2 OFFSET $3
	`, consultantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list interactions by consultant: %w", err)
	}
	defer rows.Close()

	var interactions []InteractionWithDetails
	for rows.Next() {
		var i InteractionWithDetails
		err := rows.Scan(
			&i.ID, &i.CandidateID, &i.ConsultantID, &i.Channel, &i.CategoryID, &i.ObstacleID, &i.Remarks,
			&i.NextFollowupDate, &i.NextAction, &i.SupervisorSuggestion, &i.SuggestionReadAt, &i.CreatedAt,
			&i.ConsultantName, &i.CategoryName, &i.CategorySentiment, &i.ObstacleName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interaction: %w", err)
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

// AddSupervisorSuggestion adds a supervisor suggestion to an interaction
func AddSupervisorSuggestion(ctx context.Context, id, suggestion string) error {
	_, err := pool.Exec(ctx, `
		UPDATE interactions SET supervisor_suggestion = $1, suggestion_read_at = NULL WHERE id = $2
	`, suggestion, id)
	if err != nil {
		return fmt.Errorf("failed to add supervisor suggestion: %w", err)
	}
	return nil
}

// MarkSuggestionRead marks a supervisor suggestion as read
func MarkSuggestionRead(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE interactions SET suggestion_read_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to mark suggestion as read: %w", err)
	}
	return nil
}

// GetPendingFollowups returns candidates with overdue or upcoming followups
func GetPendingFollowups(ctx context.Context, consultantID string, days int) ([]InteractionWithDetails, error) {
	rows, err := pool.Query(ctx, `
		SELECT DISTINCT ON (i.candidate_id)
		       i.id, i.candidate_id, i.consultant_id, i.channel, i.category_id, i.obstacle_id, i.remarks,
		       i.next_followup_date, i.next_action, i.supervisor_suggestion, i.suggestion_read_at, i.created_at,
		       u.name as consultant_name,
		       ic.name as category_name, ic.sentiment as category_sentiment,
		       o.name as obstacle_name
		FROM interactions i
		LEFT JOIN users u ON u.id = i.consultant_id
		LEFT JOIN interaction_categories ic ON ic.id = i.category_id
		LEFT JOIN obstacles o ON o.id = i.obstacle_id
		JOIN candidates c ON c.id = i.candidate_id
		WHERE i.consultant_id = $1
		  AND i.next_followup_date IS NOT NULL
		  AND i.next_followup_date <= CURRENT_DATE + $2
		  AND c.status NOT IN ('enrolled', 'lost')
		ORDER BY i.candidate_id, i.created_at DESC
	`, consultantID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending followups: %w", err)
	}
	defer rows.Close()

	var interactions []InteractionWithDetails
	for rows.Next() {
		var i InteractionWithDetails
		err := rows.Scan(
			&i.ID, &i.CandidateID, &i.ConsultantID, &i.Channel, &i.CategoryID, &i.ObstacleID, &i.Remarks,
			&i.NextFollowupDate, &i.NextAction, &i.SupervisorSuggestion, &i.SuggestionReadAt, &i.CreatedAt,
			&i.ConsultantName, &i.CategoryName, &i.CategorySentiment, &i.ObstacleName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan followup: %w", err)
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

// CountUnreadSuggestions counts unread supervisor suggestions for a consultant
func CountUnreadSuggestions(ctx context.Context, consultantID string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM interactions
		WHERE consultant_id = $1
		  AND supervisor_suggestion IS NOT NULL
		  AND suggestion_read_at IS NULL
	`, consultantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread suggestions: %w", err)
	}
	return count, nil
}
