package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// handleLostReasonsSettings handles GET /admin/settings/lost-reasons
func (h *AdminHandler) handleLostReasonsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Alasan Kehilangan")

	dbReasons, err := model.ListLostReasons(r.Context(), false)
	if err != nil {
		slog.Error("failed to list lost reasons", "error", err)
		http.Error(w, "Failed to load lost reasons", http.StatusInternalServerError)
		return
	}

	reasons := make([]admin.LostReasonItem, len(dbReasons))
	for i, lr := range dbReasons {
		description := ""
		if lr.Description != nil {
			description = *lr.Description
		}
		reasons[i] = admin.LostReasonItem{
			ID:           lr.ID,
			Name:         lr.Name,
			Description:  description,
			DisplayOrder: lr.DisplayOrder,
			IsActive:     lr.IsActive,
		}
	}

	admin.SettingsLostReasons(data, reasons).Render(r.Context(), w)
}

// handleCreateLostReason handles POST /admin/settings/lost-reasons
func (h *AdminHandler) handleCreateLostReason(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	descriptionStr := r.FormValue("description")
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	lr, err := model.CreateLostReason(r.Context(), name, description, displayOrder)
	if err != nil {
		slog.Error("failed to create lost reason", "error", err)
		http.Error(w, "Failed to create lost reason", http.StatusInternalServerError)
		return
	}

	slog.Info("lost reason created", "id", lr.ID, "name", lr.Name)

	descStr := ""
	if lr.Description != nil {
		descStr = *lr.Description
	}

	admin.LostReasonCard(admin.LostReasonItem{
		ID:           lr.ID,
		Name:         lr.Name,
		Description:  descStr,
		DisplayOrder: lr.DisplayOrder,
		IsActive:     lr.IsActive,
	}).Render(r.Context(), w)
}

// handleUpdateLostReason handles POST /admin/settings/lost-reasons/{id}
func (h *AdminHandler) handleUpdateLostReason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing lost reason ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	descriptionStr := r.FormValue("description")
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	err := model.UpdateLostReason(r.Context(), id, name, description, displayOrder)
	if err != nil {
		slog.Error("failed to update lost reason", "error", err)
		http.Error(w, "Failed to update lost reason", http.StatusInternalServerError)
		return
	}

	slog.Info("lost reason updated", "id", id, "name", name)

	lr, err := model.FindLostReasonByID(r.Context(), id)
	if err != nil || lr == nil {
		slog.Error("failed to find lost reason after update", "error", err)
		http.Error(w, "Failed to load lost reason", http.StatusInternalServerError)
		return
	}

	descStr := ""
	if lr.Description != nil {
		descStr = *lr.Description
	}

	admin.LostReasonCard(admin.LostReasonItem{
		ID:           lr.ID,
		Name:         lr.Name,
		Description:  descStr,
		DisplayOrder: lr.DisplayOrder,
		IsActive:     lr.IsActive,
	}).Render(r.Context(), w)
}

// handleToggleLostReasonActive handles POST /admin/settings/lost-reasons/{id}/toggle-active
func (h *AdminHandler) handleToggleLostReasonActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing lost reason ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleLostReasonActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle lost reason active", "error", err)
		http.Error(w, "Failed to toggle lost reason status", http.StatusInternalServerError)
		return
	}

	slog.Info("lost reason active status toggled", "id", id)

	lr, err := model.FindLostReasonByID(r.Context(), id)
	if err != nil || lr == nil {
		slog.Error("failed to find lost reason after toggle", "error", err)
		http.Error(w, "Failed to load lost reason", http.StatusInternalServerError)
		return
	}

	descStr := ""
	if lr.Description != nil {
		descStr = *lr.Description
	}

	admin.LostReasonCard(admin.LostReasonItem{
		ID:           lr.ID,
		Name:         lr.Name,
		Description:  descStr,
		DisplayOrder: lr.DisplayOrder,
		IsActive:     lr.IsActive,
	}).Render(r.Context(), w)
}
