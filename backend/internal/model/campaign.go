package model

import (
	"context"
	"fmt"
	"time"
)

// Campaign represents a marketing campaign
type Campaign struct {
	ID                      string     `json:"id"`
	Name                    string     `json:"name"`
	Code                    *string    `json:"code,omitempty"`
	Type                    string     `json:"type"`
	Channel                 *string    `json:"channel,omitempty"`
	Description             *string    `json:"description,omitempty"`
	StartDate               *time.Time `json:"start_date,omitempty"`
	EndDate                 *time.Time `json:"end_date,omitempty"`
	RegistrationFeeOverride *int64     `json:"registration_fee_override,omitempty"`
	IsActive                bool       `json:"is_active"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// CreateCampaign creates a new campaign
func CreateCampaign(ctx context.Context, name, campaignType string, channel, description *string, startDate, endDate *time.Time, registrationFeeOverride *int64) (*Campaign, error) {
	var c Campaign
	err := pool.QueryRow(ctx, `
		INSERT INTO campaigns (name, type, channel, description, start_date, end_date, registration_fee_override)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, code, type, channel, description, start_date, end_date, registration_fee_override, is_active, created_at, updated_at
	`, name, campaignType, channel, description, startDate, endDate, registrationFeeOverride).Scan(
		&c.ID, &c.Name, &c.Code, &c.Type, &c.Channel, &c.Description,
		&c.StartDate, &c.EndDate, &c.RegistrationFeeOverride, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}
	return &c, nil
}

// FindCampaignByID finds a campaign by ID
func FindCampaignByID(ctx context.Context, id string) (*Campaign, error) {
	var c Campaign
	err := pool.QueryRow(ctx, `
		SELECT id, name, code, type, channel, description, start_date, end_date,
		       registration_fee_override, is_active, created_at, updated_at
		FROM campaigns WHERE id = $1
	`, id).Scan(
		&c.ID, &c.Name, &c.Code, &c.Type, &c.Channel, &c.Description,
		&c.StartDate, &c.EndDate, &c.RegistrationFeeOverride, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find campaign: %w", err)
	}
	return &c, nil
}

// ListCampaigns returns all campaigns, optionally filtered by active status
func ListCampaigns(ctx context.Context, activeOnly bool) ([]Campaign, error) {
	query := `
		SELECT id, name, code, type, channel, description, start_date, end_date,
		       registration_fee_override, is_active, created_at, updated_at
		FROM campaigns
	`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY created_at DESC"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		err := rows.Scan(
			&c.ID, &c.Name, &c.Code, &c.Type, &c.Channel, &c.Description,
			&c.StartDate, &c.EndDate, &c.RegistrationFeeOverride, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}

// UpdateCampaign updates a campaign
func UpdateCampaign(ctx context.Context, id, name, campaignType string, channel, description *string, startDate, endDate *time.Time, registrationFeeOverride *int64) error {
	_, err := pool.Exec(ctx, `
		UPDATE campaigns
		SET name = $1, type = $2, channel = $3, description = $4,
		    start_date = $5, end_date = $6, registration_fee_override = $7, updated_at = NOW()
		WHERE id = $8
	`, name, campaignType, channel, description, startDate, endDate, registrationFeeOverride, id)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	return nil
}

// ToggleCampaignActive toggles the active status of a campaign
func ToggleCampaignActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE campaigns SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle campaign active: %w", err)
	}
	return nil
}

// DeleteCampaign deletes a campaign
func DeleteCampaign(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM campaigns WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	return nil
}
