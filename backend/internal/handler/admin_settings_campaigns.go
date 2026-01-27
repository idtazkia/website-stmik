package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleCampaignsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kampanye")

	// Fetch campaigns from database
	dbCampaigns, err := model.ListCampaigns(r.Context(), false)
	if err != nil {
		slog.Error("failed to list campaigns", "error", err)
		http.Error(w, "Failed to load campaigns", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	campaigns := make([]admin.CampaignItem, len(dbCampaigns))
	for i, c := range dbCampaigns {
		var channel, description string
		if c.Channel != nil {
			channel = *c.Channel
		}
		if c.Description != nil {
			description = *c.Description
		}

		startDate := ""
		endDate := ""
		if c.StartDate != nil {
			startDate = c.StartDate.Format("2006-01-02")
		}
		if c.EndDate != nil {
			endDate = c.EndDate.Format("2006-01-02")
		}

		feeOverrideStr := ""
		if c.RegistrationFeeOverride != nil {
			feeOverrideStr = formatRupiah(*c.RegistrationFeeOverride)
		}

		campaigns[i] = admin.CampaignItem{
			ID:                      c.ID,
			Name:                    c.Name,
			Type:                    c.Type,
			Channel:                 channel,
			Description:             description,
			StartDate:               startDate,
			EndDate:                 endDate,
			RegistrationFeeOverride: c.RegistrationFeeOverride,
			FeeOverrideStr:          feeOverrideStr,
			IsActive:                c.IsActive,
		}
	}

	admin.SettingsCampaigns(data, campaigns).Render(r.Context(), w)
}

func (h *AdminHandler) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	campaignType := r.FormValue("type")
	channel := r.FormValue("channel")
	description := r.FormValue("description")
	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")
	feeOverrideStr := r.FormValue("registration_fee_override")

	if name == "" || campaignType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	var channelPtr, descPtr *string
	if channel != "" {
		channelPtr = &channel
	}
	if description != "" {
		descPtr = &description
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var feeOverride *int64
	if feeOverrideStr != "" {
		var fee int64
		if _, err := fmt.Sscanf(feeOverrideStr, "%d", &fee); err == nil {
			feeOverride = &fee
		}
	}

	campaign, err := model.CreateCampaign(r.Context(), name, campaignType, channelPtr, descPtr, startDate, endDate, feeOverride)
	if err != nil {
		slog.Error("failed to create campaign", "error", err)
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign created", "campaign_id", campaign.ID)

	// Return new campaign row
	h.renderCampaignRow(w, r, campaign.ID)
}

func (h *AdminHandler) handleUpdateCampaign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	campaignType := r.FormValue("type")
	channel := r.FormValue("channel")
	description := r.FormValue("description")
	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")
	feeOverrideStr := r.FormValue("registration_fee_override")

	if name == "" || campaignType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	var channelPtr, descPtr *string
	if channel != "" {
		channelPtr = &channel
	}
	if description != "" {
		descPtr = &description
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var feeOverride *int64
	if feeOverrideStr != "" {
		var fee int64
		if _, err := fmt.Sscanf(feeOverrideStr, "%d", &fee); err == nil {
			feeOverride = &fee
		}
	}

	if err := model.UpdateCampaign(r.Context(), id, name, campaignType, channelPtr, descPtr, startDate, endDate, feeOverride); err != nil {
		slog.Error("failed to update campaign", "error", err, "campaign_id", id)
		http.Error(w, "Failed to update campaign", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign updated", "campaign_id", id)

	// Return updated campaign row
	h.renderCampaignRow(w, r, id)
}

func (h *AdminHandler) handleToggleCampaignActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleCampaignActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle campaign active", "error", err, "campaign_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign active status toggled", "campaign_id", id)

	// Return updated campaign row
	h.renderCampaignRow(w, r, id)
}

func (h *AdminHandler) renderCampaignRow(w http.ResponseWriter, r *http.Request, campaignID string) {
	campaign, err := model.FindCampaignByID(r.Context(), campaignID)
	if err != nil {
		slog.Error("failed to find campaign", "error", err)
		http.Error(w, "Failed to load campaign", http.StatusInternalServerError)
		return
	}

	if campaign == nil {
		http.NotFound(w, r)
		return
	}

	var channel, description string
	if campaign.Channel != nil {
		channel = *campaign.Channel
	}
	if campaign.Description != nil {
		description = *campaign.Description
	}

	startDate := ""
	endDate := ""
	if campaign.StartDate != nil {
		startDate = campaign.StartDate.Format("2006-01-02")
	}
	if campaign.EndDate != nil {
		endDate = campaign.EndDate.Format("2006-01-02")
	}

	feeOverrideStr := ""
	if campaign.RegistrationFeeOverride != nil {
		feeOverrideStr = formatRupiah(*campaign.RegistrationFeeOverride)
	}

	item := admin.CampaignItem{
		ID:                      campaign.ID,
		Name:                    campaign.Name,
		Type:                    campaign.Type,
		Channel:                 channel,
		Description:             description,
		StartDate:               startDate,
		EndDate:                 endDate,
		RegistrationFeeOverride: campaign.RegistrationFeeOverride,
		FeeOverrideStr:          feeOverrideStr,
		IsActive:                campaign.IsActive,
	}

	admin.CampaignRow(item).Render(r.Context(), w)
}
