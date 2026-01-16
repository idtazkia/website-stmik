package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// handleDocumentTypesSettings handles GET /admin/settings/document-types
func (h *AdminHandler) handleDocumentTypesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Jenis Dokumen")

	dbDocTypes, err := model.ListDocumentTypes(r.Context(), false)
	if err != nil {
		slog.Error("failed to list document types", "error", err)
		http.Error(w, "Failed to load document types", http.StatusInternalServerError)
		return
	}

	docTypes := make([]admin.DocumentTypeItem, len(dbDocTypes))
	for i, dt := range dbDocTypes {
		description := ""
		if dt.Description != nil {
			description = *dt.Description
		}
		docTypes[i] = admin.DocumentTypeItem{
			ID:            dt.ID,
			Name:          dt.Name,
			Code:          dt.Code,
			Description:   description,
			IsRequired:    dt.IsRequired,
			CanDefer:      dt.CanDefer,
			MaxFileSizeMB: dt.MaxFileSizeMB,
			DisplayOrder:  dt.DisplayOrder,
			IsActive:      dt.IsActive,
		}
	}

	admin.SettingsDocuments(data, docTypes).Render(r.Context(), w)
}

// handleCreateDocumentType handles POST /admin/settings/document-types
func (h *AdminHandler) handleCreateDocumentType(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	isRequired := r.FormValue("is_required") == "on"
	canDefer := r.FormValue("can_defer") == "on"
	maxFileSizeMB, _ := strconv.Atoi(r.FormValue("max_file_size_mb"))
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" || code == "" {
		http.Error(w, "Name and code are required", http.StatusBadRequest)
		return
	}

	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5
	}

	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	dt, err := model.CreateDocumentType(r.Context(), name, code, descPtr, isRequired, canDefer, maxFileSizeMB, displayOrder)
	if err != nil {
		slog.Error("failed to create document type", "error", err)
		http.Error(w, "Failed to create document type", http.StatusInternalServerError)
		return
	}

	slog.Info("document type created", "id", dt.ID, "name", dt.Name, "code", dt.Code)

	descStr := ""
	if dt.Description != nil {
		descStr = *dt.Description
	}

	// For create, use DocumentTypeRow which includes the edit modal
	admin.DocumentTypeRow(admin.DocumentTypeItem{
		ID:            dt.ID,
		Name:          dt.Name,
		Code:          dt.Code,
		Description:   descStr,
		IsRequired:    dt.IsRequired,
		CanDefer:      dt.CanDefer,
		MaxFileSizeMB: dt.MaxFileSizeMB,
		DisplayOrder:  dt.DisplayOrder,
		IsActive:      dt.IsActive,
	}).Render(r.Context(), w)
}

// handleUpdateDocumentType handles POST /admin/settings/document-types/{id}
func (h *AdminHandler) handleUpdateDocumentType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing document type ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	isRequired := r.FormValue("is_required") == "on"
	canDefer := r.FormValue("can_defer") == "on"
	maxFileSizeMB, _ := strconv.Atoi(r.FormValue("max_file_size_mb"))
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" || code == "" {
		http.Error(w, "Name and code are required", http.StatusBadRequest)
		return
	}

	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5
	}

	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	err := model.UpdateDocumentType(r.Context(), id, name, code, descPtr, isRequired, canDefer, maxFileSizeMB, displayOrder)
	if err != nil {
		slog.Error("failed to update document type", "error", err)
		http.Error(w, "Failed to update document type", http.StatusInternalServerError)
		return
	}

	slog.Info("document type updated", "id", id, "name", name, "code", code)

	dt, err := model.FindDocumentTypeByID(r.Context(), id)
	if err != nil || dt == nil {
		slog.Error("failed to find document type after update", "error", err)
		http.Error(w, "Failed to load document type", http.StatusInternalServerError)
		return
	}

	descStr := ""
	if dt.Description != nil {
		descStr = *dt.Description
	}

	admin.DocumentTypeRowOnly(admin.DocumentTypeItem{
		ID:            dt.ID,
		Name:          dt.Name,
		Code:          dt.Code,
		Description:   descStr,
		IsRequired:    dt.IsRequired,
		CanDefer:      dt.CanDefer,
		MaxFileSizeMB: dt.MaxFileSizeMB,
		DisplayOrder:  dt.DisplayOrder,
		IsActive:      dt.IsActive,
	}).Render(r.Context(), w)
}

// handleToggleDocumentTypeActive handles POST /admin/settings/document-types/{id}/toggle-active
func (h *AdminHandler) handleToggleDocumentTypeActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing document type ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleDocumentTypeActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle document type active", "error", err)
		http.Error(w, "Failed to toggle document type status", http.StatusInternalServerError)
		return
	}

	slog.Info("document type active status toggled", "id", id)

	dt, err := model.FindDocumentTypeByID(r.Context(), id)
	if err != nil || dt == nil {
		slog.Error("failed to find document type after toggle", "error", err)
		http.Error(w, "Failed to load document type", http.StatusInternalServerError)
		return
	}

	descStr := ""
	if dt.Description != nil {
		descStr = *dt.Description
	}

	admin.DocumentTypeRowOnly(admin.DocumentTypeItem{
		ID:            dt.ID,
		Name:          dt.Name,
		Code:          dt.Code,
		Description:   descStr,
		IsRequired:    dt.IsRequired,
		CanDefer:      dt.CanDefer,
		MaxFileSizeMB: dt.MaxFileSizeMB,
		DisplayOrder:  dt.DisplayOrder,
		IsActive:      dt.IsActive,
	}).Render(r.Context(), w)
}
