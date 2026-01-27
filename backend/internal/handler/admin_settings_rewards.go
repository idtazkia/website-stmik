package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleRewardsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Konfigurasi Reward")

	// Fetch reward configs from database
	dbRewards, err := model.ListRewardConfigs(r.Context())
	if err != nil {
		slog.Error("failed to list reward configs", "error", err)
		http.Error(w, "Failed to load reward configs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	rewards := make([]admin.RewardConfigItem, len(dbRewards))
	for i, r := range dbRewards {
		amountStr := formatRupiah(r.Amount)
		if r.IsPercentage {
			amountStr = fmt.Sprintf("%d%%", r.Amount)
		}
		description := ""
		if r.Description != nil {
			description = *r.Description
		}
		rewards[i] = admin.RewardConfigItem{
			ID:           r.ID,
			ReferrerType: r.ReferrerType,
			RewardType:   r.RewardType,
			Amount:       r.Amount,
			AmountStr:    amountStr,
			IsPercentage: r.IsPercentage,
			TriggerEvent: r.TriggerEvent,
			Description:  description,
			IsActive:     r.IsActive,
		}
	}

	// Fetch MGM reward configs from database
	dbMGMRewards, err := model.ListMGMRewardConfigs(r.Context())
	if err != nil {
		slog.Error("failed to list MGM reward configs", "error", err)
		http.Error(w, "Failed to load MGM reward configs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	mgmRewards := make([]admin.MGMRewardConfigItem, len(dbMGMRewards))
	for i, m := range dbMGMRewards {
		referrerStr := formatRupiah(m.ReferrerAmount)
		refereeStr := ""
		if m.RefereeAmount != nil {
			refereeStr = formatRupiah(*m.RefereeAmount)
		}
		description := ""
		if m.Description != nil {
			description = *m.Description
		}
		mgmRewards[i] = admin.MGMRewardConfigItem{
			ID:             m.ID,
			AcademicYear:   m.AcademicYear,
			RewardType:     m.RewardType,
			ReferrerAmount: m.ReferrerAmount,
			ReferrerStr:    referrerStr,
			RefereeAmount:  m.RefereeAmount,
			RefereeStr:     refereeStr,
			TriggerEvent:   m.TriggerEvent,
			Description:    description,
			IsActive:       m.IsActive,
		}
	}

	admin.SettingsRewards(data, rewards, mgmRewards).Render(r.Context(), w)
}

// Reward Config Handlers

func (h *AdminHandler) handleCreateRewardConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	referrerType := r.FormValue("referrer_type")
	rewardType := r.FormValue("reward_type")
	amountStr := r.FormValue("amount")
	isPercentage := r.FormValue("is_percentage") == "on"
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	reward, err := model.CreateRewardConfig(r.Context(), referrerType, rewardType, amount, isPercentage, triggerEvent, description)
	if err != nil {
		slog.Error("failed to create reward config", "error", err)
		http.Error(w, "Failed to create reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config created", "reward_id", reward.ID)
	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) handleUpdateRewardConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing reward ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	referrerType := r.FormValue("referrer_type")
	rewardType := r.FormValue("reward_type")
	amountStr := r.FormValue("amount")
	isPercentage := r.FormValue("is_percentage") == "on"
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	err = model.UpdateRewardConfig(r.Context(), id, referrerType, rewardType, amount, isPercentage, triggerEvent, description)
	if err != nil {
		slog.Error("failed to update reward config", "error", err)
		http.Error(w, "Failed to update reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config updated", "reward_id", id)

	// Fetch updated reward and render
	reward, err := model.FindRewardConfigByID(r.Context(), id)
	if err != nil || reward == nil {
		http.Error(w, "Failed to fetch updated reward", http.StatusInternalServerError)
		return
	}

	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) handleToggleRewardConfigActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing reward ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleRewardConfigActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle reward config active", "error", err)
		http.Error(w, "Failed to toggle reward status", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config active status toggled", "reward_id", id)

	// Fetch updated reward and render
	reward, err := model.FindRewardConfigByID(r.Context(), id)
	if err != nil || reward == nil {
		http.Error(w, "Failed to fetch updated reward", http.StatusInternalServerError)
		return
	}

	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) renderRewardCard(w http.ResponseWriter, r *http.Request, reward *model.RewardConfig) {
	amountStr := formatRupiah(reward.Amount)
	if reward.IsPercentage {
		amountStr = fmt.Sprintf("%d%%", reward.Amount)
	}
	description := ""
	if reward.Description != nil {
		description = *reward.Description
	}

	item := admin.RewardConfigItem{
		ID:           reward.ID,
		ReferrerType: reward.ReferrerType,
		RewardType:   reward.RewardType,
		Amount:       reward.Amount,
		AmountStr:    amountStr,
		IsPercentage: reward.IsPercentage,
		TriggerEvent: reward.TriggerEvent,
		Description:  description,
		IsActive:     reward.IsActive,
	}

	admin.RewardCard(item).Render(r.Context(), w)
}

// MGM Reward Config Handlers

func (h *AdminHandler) handleCreateMGMRewardConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	academicYear := r.FormValue("academic_year")
	rewardType := r.FormValue("reward_type")
	referrerAmountStr := r.FormValue("referrer_amount")
	refereeAmountStr := r.FormValue("referee_amount")
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	referrerAmount, err := strconv.ParseInt(referrerAmountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid referrer amount", http.StatusBadRequest)
		return
	}

	var refereeAmount *int64
	if refereeAmountStr != "" {
		amt, err := strconv.ParseInt(refereeAmountStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid referee amount", http.StatusBadRequest)
			return
		}
		refereeAmount = &amt
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	mgmReward, err := model.CreateMGMRewardConfig(r.Context(), academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description)
	if err != nil {
		slog.Error("failed to create MGM reward config", "error", err)
		http.Error(w, "Failed to create MGM reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config created", "mgm_reward_id", mgmReward.ID)
	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) handleUpdateMGMRewardConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing MGM reward ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	academicYear := r.FormValue("academic_year")
	rewardType := r.FormValue("reward_type")
	referrerAmountStr := r.FormValue("referrer_amount")
	refereeAmountStr := r.FormValue("referee_amount")
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	referrerAmount, err := strconv.ParseInt(referrerAmountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid referrer amount", http.StatusBadRequest)
		return
	}

	var refereeAmount *int64
	if refereeAmountStr != "" {
		amt, err := strconv.ParseInt(refereeAmountStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid referee amount", http.StatusBadRequest)
			return
		}
		refereeAmount = &amt
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	err = model.UpdateMGMRewardConfig(r.Context(), id, academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description)
	if err != nil {
		slog.Error("failed to update MGM reward config", "error", err)
		http.Error(w, "Failed to update MGM reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config updated", "mgm_reward_id", id)

	// Fetch updated MGM reward and render
	mgmReward, err := model.FindMGMRewardConfigByID(r.Context(), id)
	if err != nil || mgmReward == nil {
		http.Error(w, "Failed to fetch updated MGM reward", http.StatusInternalServerError)
		return
	}

	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) handleToggleMGMRewardConfigActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing MGM reward ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleMGMRewardConfigActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle MGM reward config active", "error", err)
		http.Error(w, "Failed to toggle MGM reward status", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config active status toggled", "mgm_reward_id", id)

	// Fetch updated MGM reward and render
	mgmReward, err := model.FindMGMRewardConfigByID(r.Context(), id)
	if err != nil || mgmReward == nil {
		http.Error(w, "Failed to fetch updated MGM reward", http.StatusInternalServerError)
		return
	}

	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) renderMGMRewardCard(w http.ResponseWriter, r *http.Request, mgmReward *model.MGMRewardConfig) {
	referrerStr := formatRupiah(mgmReward.ReferrerAmount)
	refereeStr := ""
	if mgmReward.RefereeAmount != nil {
		refereeStr = formatRupiah(*mgmReward.RefereeAmount)
	}
	description := ""
	if mgmReward.Description != nil {
		description = *mgmReward.Description
	}

	item := admin.MGMRewardConfigItem{
		ID:             mgmReward.ID,
		AcademicYear:   mgmReward.AcademicYear,
		RewardType:     mgmReward.RewardType,
		ReferrerAmount: mgmReward.ReferrerAmount,
		ReferrerStr:    referrerStr,
		RefereeAmount:  mgmReward.RefereeAmount,
		RefereeStr:     refereeStr,
		TriggerEvent:   mgmReward.TriggerEvent,
		Description:    description,
		IsActive:       mgmReward.IsActive,
	}

	admin.MGMRewardCard(item).Render(r.Context(), w)
}
