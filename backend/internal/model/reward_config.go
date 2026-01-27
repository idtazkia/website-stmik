package model

import (
	"context"
	"fmt"
	"time"
)

// RewardConfig represents reward configuration for external referrers
type RewardConfig struct {
	ID           string    `json:"id"`
	ReferrerType string    `json:"referrer_type"`
	RewardType   string    `json:"reward_type"`
	Amount       int64     `json:"amount"`
	IsPercentage bool      `json:"is_percentage"`
	TriggerEvent string    `json:"trigger_event"`
	Description  *string   `json:"description,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateRewardConfig creates a new reward config
func CreateRewardConfig(ctx context.Context, referrerType, rewardType string, amount int64, isPercentage bool, triggerEvent string, description *string) (*RewardConfig, error) {
	var r RewardConfig
	err := pool.QueryRow(ctx, `
		INSERT INTO reward_configs (referrer_type, reward_type, amount, is_percentage, trigger_event, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, referrer_type, reward_type, amount, is_percentage, trigger_event, description, is_active, created_at, updated_at
	`, referrerType, rewardType, amount, isPercentage, triggerEvent, description).Scan(
		&r.ID, &r.ReferrerType, &r.RewardType, &r.Amount, &r.IsPercentage,
		&r.TriggerEvent, &r.Description, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reward config: %w", err)
	}
	return &r, nil
}

// FindRewardConfigByID finds a reward config by ID
func FindRewardConfigByID(ctx context.Context, id string) (*RewardConfig, error) {
	var r RewardConfig
	err := pool.QueryRow(ctx, `
		SELECT id, referrer_type, reward_type, amount, is_percentage, trigger_event,
		       description, is_active, created_at, updated_at
		FROM reward_configs WHERE id = $1
	`, id).Scan(
		&r.ID, &r.ReferrerType, &r.RewardType, &r.Amount, &r.IsPercentage,
		&r.TriggerEvent, &r.Description, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find reward config: %w", err)
	}
	return &r, nil
}

// FindRewardConfigByType finds active reward config by referrer type
func FindRewardConfigByType(ctx context.Context, referrerType string) ([]RewardConfig, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, referrer_type, reward_type, amount, is_percentage, trigger_event,
		       description, is_active, created_at, updated_at
		FROM reward_configs
		WHERE referrer_type = $1 AND is_active = true
		ORDER BY trigger_event
	`, referrerType)
	if err != nil {
		return nil, fmt.Errorf("failed to find reward config by type: %w", err)
	}
	defer rows.Close()

	var configs []RewardConfig
	for rows.Next() {
		var r RewardConfig
		err := rows.Scan(
			&r.ID, &r.ReferrerType, &r.RewardType, &r.Amount, &r.IsPercentage,
			&r.TriggerEvent, &r.Description, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward config: %w", err)
		}
		configs = append(configs, r)
	}
	return configs, nil
}

// ListRewardConfigs returns all reward configs
func ListRewardConfigs(ctx context.Context) ([]RewardConfig, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, referrer_type, reward_type, amount, is_percentage, trigger_event,
		       description, is_active, created_at, updated_at
		FROM reward_configs
		ORDER BY referrer_type, trigger_event
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list reward configs: %w", err)
	}
	defer rows.Close()

	var configs []RewardConfig
	for rows.Next() {
		var r RewardConfig
		err := rows.Scan(
			&r.ID, &r.ReferrerType, &r.RewardType, &r.Amount, &r.IsPercentage,
			&r.TriggerEvent, &r.Description, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward config: %w", err)
		}
		configs = append(configs, r)
	}
	return configs, nil
}

// UpdateRewardConfig updates a reward config
func UpdateRewardConfig(ctx context.Context, id, referrerType, rewardType string, amount int64, isPercentage bool, triggerEvent string, description *string) error {
	_, err := pool.Exec(ctx, `
		UPDATE reward_configs
		SET referrer_type = $1, reward_type = $2, amount = $3, is_percentage = $4,
		    trigger_event = $5, description = $6, updated_at = NOW()
		WHERE id = $7
	`, referrerType, rewardType, amount, isPercentage, triggerEvent, description, id)
	if err != nil {
		return fmt.Errorf("failed to update reward config: %w", err)
	}
	return nil
}

// ToggleRewardConfigActive toggles the active status
func ToggleRewardConfigActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE reward_configs SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle reward config active: %w", err)
	}
	return nil
}

// DeleteRewardConfig deletes a reward config
func DeleteRewardConfig(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM reward_configs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete reward config: %w", err)
	}
	return nil
}
