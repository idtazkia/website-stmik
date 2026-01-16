package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

func (h *AdminHandler) handleProgramsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Prodi")

	// Fetch programs from database
	dbProdis, err := model.ListProdis(r.Context(), false)
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
		http.Error(w, "Failed to load programs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	programs := make([]admin.ProgramItem, len(dbProdis))
	for i, p := range dbProdis {
		status := "inactive"
		if p.IsActive {
			status = "active"
		}

		programs[i] = admin.ProgramItem{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Level:    p.Degree,
			SPPFee:   "-",       // TODO: fetch from fee_structures
			Status:   status,
			Students: "-",       // TODO: count enrolled candidates
		}
	}

	admin.SettingsPrograms(data, programs).Render(r.Context(), w)
}

// handleCreateProgram handles POST /admin/settings/programs
func (h *AdminHandler) handleCreateProgram(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	degree := r.FormValue("degree")

	if name == "" || code == "" || degree == "" {
		http.Error(w, "Name, code, and degree are required", http.StatusBadRequest)
		return
	}

	prodi, err := model.CreateProdi(r.Context(), name, code, degree)
	if err != nil {
		slog.Error("failed to create prodi", "error", err)
		http.Error(w, "Failed to create program", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi created", "prodi_id", prodi.ID, "code", prodi.Code)

	// Return the new program card
	h.renderProgramCard(w, r, prodi.ID)
}

// handleUpdateProgram handles POST /admin/settings/programs/{id}
func (h *AdminHandler) handleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	degree := r.FormValue("degree")

	if name == "" || code == "" || degree == "" {
		http.Error(w, "Name, code, and degree are required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateProdi(r.Context(), id, name, code, degree); err != nil {
		slog.Error("failed to update prodi", "error", err, "prodi_id", id)
		http.Error(w, "Failed to update program", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi updated", "prodi_id", id, "code", code)

	// Return updated program card
	h.renderProgramCard(w, r, id)
}

// handleToggleProgramActive handles POST /admin/settings/programs/{id}/toggle-active
func (h *AdminHandler) handleToggleProgramActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleProdiActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle prodi active", "error", err, "prodi_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi active status toggled", "prodi_id", id)

	// Return updated program card
	h.renderProgramCard(w, r, id)
}

// renderProgramCard renders a single program card for HTMX updates
func (h *AdminHandler) renderProgramCard(w http.ResponseWriter, r *http.Request, prodiID string) {
	prodi, err := model.FindProdiByID(r.Context(), prodiID)
	if err != nil {
		slog.Error("failed to find prodi", "error", err)
		http.Error(w, "Failed to load program", http.StatusInternalServerError)
		return
	}

	if prodi == nil {
		http.NotFound(w, r)
		return
	}

	status := "inactive"
	if prodi.IsActive {
		status = "active"
	}

	programItem := admin.ProgramItem{
		ID:       prodi.ID,
		Name:     prodi.Name,
		Code:     prodi.Code,
		Level:    prodi.Degree,
		SPPFee:   "-",   // TODO: fetch from fee_structures
		Status:   status,
		Students: "-",   // TODO: count enrolled candidates
	}

	admin.ProgramCard(programItem).Render(r.Context(), w)
}
