package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

func (h *AdminHandler) handleAnnouncementsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Pengumuman")

	// Fetch announcements from database
	dbAnnouncements, err := model.ListAnnouncements(r.Context())
	if err != nil {
		slog.Error("failed to list announcements", "error", err)
		http.Error(w, "Failed to load announcements", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	announcements := make([]admin.AnnouncementItem, len(dbAnnouncements))
	for i, a := range dbAnnouncements {
		var targetStatus, targetProdiID, targetProdiName, createdByName string
		if a.TargetStatus != nil {
			targetStatus = *a.TargetStatus
		}
		if a.TargetProdiID != nil {
			targetProdiID = *a.TargetProdiID
		}
		if a.TargetProdiName != nil {
			targetProdiName = *a.TargetProdiName
		}
		if a.CreatedByName != nil {
			createdByName = *a.CreatedByName
		}

		publishedAt := ""
		if a.PublishedAt != nil {
			publishedAt = a.PublishedAt.Format("02 Jan 2006")
		}

		announcements[i] = admin.AnnouncementItem{
			ID:              a.ID,
			Title:           a.Title,
			Content:         a.Content,
			TargetStatus:    targetStatus,
			TargetProdiID:   targetProdiID,
			TargetProdiName: targetProdiName,
			IsPublished:     a.IsPublished,
			PublishedAt:     publishedAt,
			ReadCount:       a.ReadCount,
			CreatedByName:   createdByName,
			CreatedAt:       a.CreatedAt.Format("02 Jan 2006"),
		}
	}

	// Fetch prodis for dropdown
	dbProdis, err := model.ListProdis(r.Context(), true)
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
		http.Error(w, "Failed to load prodis", http.StatusInternalServerError)
		return
	}

	prodis := make([]admin.ProdiOption, len(dbProdis))
	for i, p := range dbProdis {
		prodis[i] = admin.ProdiOption{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		}
	}

	admin.SettingsAnnouncements(data, announcements, prodis).Render(r.Context(), w)
}

func (h *AdminHandler) handleCreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	targetStatus := r.FormValue("target_status")
	targetProdiID := r.FormValue("target_prodi_id")

	// Get user ID from claims
	claims := GetUserClaims(r.Context())
	var createdBy *string
	if claims != nil {
		createdBy = &claims.UserID
	}

	// Convert empty strings to nil
	var targetStatusPtr, targetProdiIDPtr *string
	if targetStatus != "" {
		targetStatusPtr = &targetStatus
	}
	if targetProdiID != "" {
		targetProdiIDPtr = &targetProdiID
	}

	// Create announcement
	ann, err := model.CreateAnnouncement(r.Context(), title, content, targetStatusPtr, targetProdiIDPtr, createdBy)
	if err != nil {
		slog.Error("failed to create announcement", "error", err)
		http.Error(w, "Failed to create announcement", http.StatusInternalServerError)
		return
	}

	slog.Info("announcement created", "id", ann.ID, "title", title)

	// Return the new row HTML
	var targetProdiName string
	if targetProdiIDPtr != nil {
		prodi, _ := model.FindProdiByID(r.Context(), *targetProdiIDPtr)
		if prodi != nil {
			targetProdiName = prodi.Name
		}
	}

	var createdByName string
	if claims != nil {
		createdByName = claims.Name
	}

	item := admin.AnnouncementItem{
		ID:              ann.ID,
		Title:           ann.Title,
		Content:         ann.Content,
		TargetStatus:    targetStatus,
		TargetProdiID:   targetProdiID,
		TargetProdiName: targetProdiName,
		IsPublished:     ann.IsPublished,
		PublishedAt:     "",
		ReadCount:       0,
		CreatedByName:   createdByName,
		CreatedAt:       ann.CreatedAt.Format("02 Jan 2006"),
	}

	admin.AnnouncementRow(item).Render(r.Context(), w)
}

func (h *AdminHandler) handleEditAnnouncementForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ann, err := model.FindAnnouncementByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to find announcement", "error", err, "id", id)
		http.Error(w, "Failed to find announcement", http.StatusInternalServerError)
		return
	}
	if ann == nil {
		http.NotFound(w, r)
		return
	}

	// Fetch prodis for dropdown
	dbProdis, err := model.ListProdis(r.Context(), true)
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
		http.Error(w, "Failed to load prodis", http.StatusInternalServerError)
		return
	}

	prodis := make([]admin.ProdiOption, len(dbProdis))
	for i, p := range dbProdis {
		prodis[i] = admin.ProdiOption{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		}
	}

	var targetStatus, targetProdiID string
	if ann.TargetStatus != nil {
		targetStatus = *ann.TargetStatus
	}
	if ann.TargetProdiID != nil {
		targetProdiID = *ann.TargetProdiID
	}

	item := admin.AnnouncementItem{
		ID:            ann.ID,
		Title:         ann.Title,
		Content:       ann.Content,
		TargetStatus:  targetStatus,
		TargetProdiID: targetProdiID,
	}

	admin.EditAnnouncementForm(item, prodis).Render(r.Context(), w)
}

func (h *AdminHandler) handleUpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	targetStatus := r.FormValue("target_status")
	targetProdiID := r.FormValue("target_prodi_id")

	// Convert empty strings to nil
	var targetStatusPtr, targetProdiIDPtr *string
	if targetStatus != "" {
		targetStatusPtr = &targetStatus
	}
	if targetProdiID != "" {
		targetProdiIDPtr = &targetProdiID
	}

	// Update announcement
	err := model.UpdateAnnouncement(r.Context(), id, title, content, targetStatusPtr, targetProdiIDPtr)
	if err != nil {
		slog.Error("failed to update announcement", "error", err, "id", id)
		http.Error(w, "Failed to update announcement", http.StatusInternalServerError)
		return
	}

	slog.Info("announcement updated", "id", id)

	// Fetch updated announcement to return
	ann, err := model.FindAnnouncementByID(r.Context(), id)
	if err != nil || ann == nil {
		http.Error(w, "Failed to fetch updated announcement", http.StatusInternalServerError)
		return
	}

	// Get prodi name
	var targetProdiName string
	if targetProdiIDPtr != nil {
		prodi, _ := model.FindProdiByID(r.Context(), *targetProdiIDPtr)
		if prodi != nil {
			targetProdiName = prodi.Name
		}
	}

	publishedAt := ""
	if ann.PublishedAt != nil {
		publishedAt = ann.PublishedAt.Format("02 Jan 2006")
	}

	item := admin.AnnouncementItem{
		ID:              ann.ID,
		Title:           ann.Title,
		Content:         ann.Content,
		TargetStatus:    targetStatus,
		TargetProdiID:   targetProdiID,
		TargetProdiName: targetProdiName,
		IsPublished:     ann.IsPublished,
		PublishedAt:     publishedAt,
		ReadCount:       0, // Could fetch from model if needed
		CreatedAt:       ann.CreatedAt.Format("02 Jan 2006"),
	}

	admin.AnnouncementRow(item).Render(r.Context(), w)
}

func (h *AdminHandler) handlePublishAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := model.PublishAnnouncement(r.Context(), id)
	if err != nil {
		slog.Error("failed to publish announcement", "error", err, "id", id)
		http.Error(w, "Failed to publish announcement", http.StatusInternalServerError)
		return
	}

	slog.Info("announcement published", "id", id)

	// Return updated row
	h.returnAnnouncementRow(w, r, id)
}

func (h *AdminHandler) handleUnpublishAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := model.UnpublishAnnouncement(r.Context(), id)
	if err != nil {
		slog.Error("failed to unpublish announcement", "error", err, "id", id)
		http.Error(w, "Failed to unpublish announcement", http.StatusInternalServerError)
		return
	}

	slog.Info("announcement unpublished", "id", id)

	// Return updated row
	h.returnAnnouncementRow(w, r, id)
}

func (h *AdminHandler) handleDeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := model.DeleteAnnouncement(r.Context(), id)
	if err != nil {
		slog.Error("failed to delete announcement", "error", err, "id", id)
		http.Error(w, "Failed to delete announcement", http.StatusInternalServerError)
		return
	}

	slog.Info("announcement deleted", "id", id)

	// Return empty response (row will be removed)
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) returnAnnouncementRow(w http.ResponseWriter, r *http.Request, id string) {
	ann, err := model.FindAnnouncementByID(r.Context(), id)
	if err != nil || ann == nil {
		http.Error(w, "Failed to fetch announcement", http.StatusInternalServerError)
		return
	}

	var targetStatus, targetProdiID, targetProdiName string
	if ann.TargetStatus != nil {
		targetStatus = *ann.TargetStatus
	}
	if ann.TargetProdiID != nil {
		targetProdiID = *ann.TargetProdiID
		prodi, _ := model.FindProdiByID(r.Context(), targetProdiID)
		if prodi != nil {
			targetProdiName = prodi.Name
		}
	}

	publishedAt := ""
	if ann.PublishedAt != nil {
		publishedAt = ann.PublishedAt.Format("02 Jan 2006")
	}

	item := admin.AnnouncementItem{
		ID:              ann.ID,
		Title:           ann.Title,
		Content:         ann.Content,
		TargetStatus:    targetStatus,
		TargetProdiID:   targetProdiID,
		TargetProdiName: targetProdiName,
		IsPublished:     ann.IsPublished,
		PublishedAt:     publishedAt,
		ReadCount:       0,
		CreatedAt:       ann.CreatedAt.Format("02 Jan 2006"),
	}

	admin.AnnouncementRow(item).Render(r.Context(), w)
}
