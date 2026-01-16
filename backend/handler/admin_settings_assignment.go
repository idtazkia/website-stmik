package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// handleAssignmentSettings handles GET /admin/settings/assignment
func (h *AdminHandler) handleAssignmentSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Algoritma Penugasan")

	dbAlgorithms, err := model.ListAssignmentAlgorithms(r.Context())
	if err != nil {
		slog.Error("failed to list assignment algorithms", "error", err)
		http.Error(w, "Failed to load assignment algorithms", http.StatusInternalServerError)
		return
	}

	algorithms := make([]admin.AssignmentAlgorithmItem, len(dbAlgorithms))
	for i, alg := range dbAlgorithms {
		description := ""
		if alg.Description != nil {
			description = *alg.Description
		}
		algorithms[i] = admin.AssignmentAlgorithmItem{
			ID:          alg.ID,
			Name:        alg.Name,
			Code:        alg.Code,
			Description: description,
			IsActive:    alg.IsActive,
		}
	}

	admin.SettingsAssignment(data, algorithms).Render(r.Context(), w)
}

// handleActivateAlgorithm handles POST /admin/settings/assignment/{id}/activate
func (h *AdminHandler) handleActivateAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing algorithm ID", http.StatusBadRequest)
		return
	}

	err := model.SetAssignmentAlgorithmActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to activate algorithm", "error", err)
		http.Error(w, "Failed to activate algorithm", http.StatusInternalServerError)
		return
	}

	slog.Info("assignment algorithm activated", "algorithm_id", id)

	// Return updated list
	dbAlgorithms, err := model.ListAssignmentAlgorithms(r.Context())
	if err != nil {
		slog.Error("failed to list assignment algorithms", "error", err)
		http.Error(w, "Failed to load assignment algorithms", http.StatusInternalServerError)
		return
	}

	algorithms := make([]admin.AssignmentAlgorithmItem, len(dbAlgorithms))
	for i, alg := range dbAlgorithms {
		description := ""
		if alg.Description != nil {
			description = *alg.Description
		}
		algorithms[i] = admin.AssignmentAlgorithmItem{
			ID:          alg.ID,
			Name:        alg.Name,
			Code:        alg.Code,
			Description: description,
			IsActive:    alg.IsActive,
		}
	}

	admin.AlgorithmList(algorithms).Render(r.Context(), w)
}
