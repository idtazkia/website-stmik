package model

import (
	"context"
	"fmt"
	"time"
)

// MGMRewardConfig represents Member-Get-Member reward configuration
type MGMRewardConfig struct {
	ID             string    `json:"id"`
	AcademicYear   string    `json:"academic_year"`
	RewardType     string    `json:"reward_type"`
	ReferrerAmount int64     `json:"referrer_amount"`
	RefereeAmount  *int64    `json:"referee_amount,omitempty"`
	TriggerEvent   string    `json:"trigger_event"`
	Description    *string   `json:"description,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateMGMRewardConfig creates a new MGM reward config
func CreateMGMRewardConfig(ctx context.Context, academicYear, rewardType string, referrerAmount int64, refereeAmount *int64, triggerEvent string, description *string) (*MGMRewardConfig, error) {
	var m MGMRewardConfig
	err := pool.QueryRow(ctx, `
		INSERT INTO mgm_reward_configs (academic_year, reward_type, referrer_amount, referee_amount, trigger_event, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, academic_year, reward_type, referrer_amount, referee_amount, trigger_event, description, is_active, created_at, updated_at
	`, academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description).Scan(
		&m.ID, &m.AcademicYear, &m.RewardType, &m.ReferrerAmount, &m.RefereeAmount,
		&m.TriggerEvent, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MGM reward config: %w", err)
	}
	return &m, nil
}

// FindMGMRewardConfigByID finds a MGM reward config by ID
func FindMGMRewardConfigByID(ctx context.Context, id string) (*MGMRewardConfig, error) {
	var m MGMRewardConfig
	err := pool.QueryRow(ctx, `
		SELECT id, academic_year, reward_type, referrer_amount, referee_amount, trigger_event,
		       description, is_active, created_at, updated_at
		FROM mgm_reward_configs WHERE id = $1
	`, id).Scan(
		&m.ID, &m.AcademicYear, &m.RewardType, &m.ReferrerAmount, &m.RefereeAmount,
		&m.TriggerEvent, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find MGM reward config: %w", err)
	}
	return &m, nil
}

// FindMGMRewardConfigByYear finds active MGM reward configs by academic year
func FindMGMRewardConfigByYear(ctx context.Context, academicYear string) ([]MGMRewardConfig, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, academic_year, reward_type, referrer_amount, referee_amount, trigger_event,
		       description, is_active, created_at, updated_at
		FROM mgm_reward_configs
		WHERE academic_year = $1 AND is_active = true
		ORDER BY trigger_event
	`, academicYear)
	if err != nil {
		return nil, fmt.Errorf("failed to find MGM reward config by year: %w", err)
	}
	defer rows.Close()

	var configs []MGMRewardConfig
	for rows.Next() {
		var m MGMRewardConfig
		err := rows.Scan(
			&m.ID, &m.AcademicYear, &m.RewardType, &m.ReferrerAmount, &m.RefereeAmount,
			&m.TriggerEvent, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan MGM reward config: %w", err)
		}
		configs = append(configs, m)
	}
	return configs, nil
}

// ListMGMRewardConfigs returns all MGM reward configs
func ListMGMRewardConfigs(ctx context.Context) ([]MGMRewardConfig, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, academic_year, reward_type, referrer_amount, referee_amount, trigger_event,
		       description, is_active, created_at, updated_at
		FROM mgm_reward_configs
		ORDER BY academic_year DESC, trigger_event
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list MGM reward configs: %w", err)
	}
	defer rows.Close()

	var configs []MGMRewardConfig
	for rows.Next() {
		var m MGMRewardConfig
		err := rows.Scan(
			&m.ID, &m.AcademicYear, &m.RewardType, &m.ReferrerAmount, &m.RefereeAmount,
			&m.TriggerEvent, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan MGM reward config: %w", err)
		}
		configs = append(configs, m)
	}
	return configs, nil
}

// UpdateMGMRewardConfig updates a MGM reward config
func UpdateMGMRewardConfig(ctx context.Context, id, academicYear, rewardType string, referrerAmount int64, refereeAmount *int64, triggerEvent string, description *string) error {
	_, err := pool.Exec(ctx, `
		UPDATE mgm_reward_configs
		SET academic_year = $1, reward_type = $2, referrer_amount = $3, referee_amount = $4,
		    trigger_event = $5, description = $6, updated_at = NOW()
		WHERE id = $7
	`, academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description, id)
	if err != nil {
		return fmt.Errorf("failed to update MGM reward config: %w", err)
	}
	return nil
}

// ToggleMGMRewardConfigActive toggles the active status
func ToggleMGMRewardConfigActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE mgm_reward_configs SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle MGM reward config active: %w", err)
	}
	return nil
}

// DeleteMGMRewardConfig deletes a MGM reward config
func DeleteMGMRewardConfig(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM mgm_reward_configs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete MGM reward config: %w", err)
	}
	return nil
}
