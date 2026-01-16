package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/portal"
)

// PortalHandler handles all portal routes for candidates
type PortalHandler struct {
	sessionMgr *auth.SessionManager
}

// NewPortalHandler creates a new portal handler
func NewPortalHandler(sessionMgr *auth.SessionManager) *PortalHandler {
	return &PortalHandler{sessionMgr: sessionMgr}
}

// RegisterRoutes registers all portal routes to the mux
func (h *PortalHandler) RegisterRoutes(mux *http.ServeMux) {
	// All portal routes require candidate authentication
	mux.Handle("GET /portal", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleDashboard)))
	mux.Handle("GET /portal/documents", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleDocuments)))
	mux.Handle("GET /portal/payments", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handlePayments)))
	mux.Handle("GET /portal/announcements", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleAnnouncements)))
	mux.Handle("POST /portal/announcements/{id}/read", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleMarkAnnouncementRead)))
	mux.Handle("GET /portal/referral", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleReferral)))
}

func (h *PortalHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get candidate data with details
	candidateData, err := model.GetCandidateDashboardData(r.Context(), claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		slog.Error("candidate not found", "id", claims.CandidateID)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := NewPageData("Dashboard")

	// Build dashboard candidate struct
	candidate := portal.DashboardCandidate{
		Name:         safeString(candidateData.Name),
		Email:        safeString(candidateData.Email),
		Phone:        safeString(candidateData.Phone),
		Program:      safeString(candidateData.ProdiName),
		Status:       candidateData.Status,
		StatusLabel:  getStatusLabel(candidateData.Status),
		RegisteredAt: candidateData.CreatedAt.Format("02 January 2006"),
	}

	if candidateData.ConsultantName != nil {
		candidate.ConsultantName = *candidateData.ConsultantName
	}
	if candidateData.ConsultantEmail != nil {
		candidate.ConsultantPhone = *candidateData.ConsultantEmail // Using email field temporarily
	}

	// Build checklist based on candidate progress
	checklist := buildChecklist(candidateData)

	// Fetch announcements for this candidate (limited to 3 for dashboard)
	dbAnnouncements, err := model.ListAnnouncementsForCandidate(r.Context(), candidateData.ID, candidateData.Status, candidateData.ProdiID)
	if err != nil {
		slog.Error("failed to list announcements", "error", err)
		// Continue with empty announcements, don't fail the dashboard
		dbAnnouncements = nil
	}

	// Convert to template type (show max 3 on dashboard)
	announcements := make([]portal.AnnouncementItem, 0)
	for i, a := range dbAnnouncements {
		if i >= 3 {
			break
		}
		publishedAt := ""
		if a.PublishedAt != nil {
			publishedAt = a.PublishedAt.Format("02 Jan 2006")
		}
		announcements = append(announcements, portal.AnnouncementItem{
			Title:   a.Title,
			Content: truncateContent(a.Content, 150),
			Date:    publishedAt,
			IsNew:   !a.IsRead,
		})
	}

	portal.Dashboard(data, candidate, checklist, announcements).Render(r.Context(), w)
}

func (h *PortalHandler) handleDocuments(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	candidateData, err := model.GetCandidateDashboardData(r.Context(), claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := NewPageData("Dokumen Saya")

	candidate := portal.PortalCandidate{
		ID:                      candidateData.ID,
		Name:                    safeString(candidateData.Name),
		Status:                  candidateData.Status,
		DocumentProgress:        "0 dari 4 dokumen",
		DocumentProgressPercent: "0",
	}

	// TODO: Fetch real documents from database when document upload feature is implemented
	documents := []portal.DocumentItem{
		{Type: "ktp", Name: "KTP", Description: "Kartu Tanda Penduduk yang masih berlaku", Required: true, CanDefer: false, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "not_uploaded"},
		{Type: "photo", Name: "Pas Foto", Description: "Foto 3x4 dengan background merah atau biru", Required: true, CanDefer: false, AllowedFormats: "JPG, PNG", AllowedMimeTypes: "image/jpeg,image/png", MaxSize: "2 MB", Status: "not_uploaded"},
		{Type: "ijazah", Name: "Ijazah", Description: "Ijazah SMA/SMK/sederajat", Required: true, CanDefer: true, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "not_uploaded"},
		{Type: "transcript", Name: "Transkrip Nilai", Description: "Transkrip nilai SMA/SMK/sederajat", Required: true, CanDefer: true, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "not_uploaded"},
	}

	portal.Documents(data, candidate, documents).Render(r.Context(), w)
}

func (h *PortalHandler) handlePayments(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	candidateData, err := model.GetCandidateDashboardData(r.Context(), claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := NewPageData("Pembayaran")
	candidate := portal.PortalCandidate{
		ID:   candidateData.ID,
		Name: safeString(candidateData.Name),
	}

	// TODO: Fetch real payments from database when payment feature is implemented
	payments := []portal.PaymentItem{}
	summary := portal.PaymentSummary{
		TotalDue:     "Rp 0",
		TotalPaid:    "Rp 0",
		TotalPending: "Rp 0",
	}

	portal.Payments(data, candidate, payments, summary).Render(r.Context(), w)
}

func (h *PortalHandler) handleAnnouncements(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	candidateData, err := model.GetCandidateDashboardData(r.Context(), claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Fetch all announcements for this candidate
	dbAnnouncements, err := model.ListAnnouncementsForCandidate(r.Context(), candidateData.ID, candidateData.Status, candidateData.ProdiID)
	if err != nil {
		slog.Error("failed to list announcements", "error", err)
		http.Error(w, "Failed to load announcements", http.StatusInternalServerError)
		return
	}

	// Count unread
	unreadCount, err := model.CountUnreadAnnouncementsForCandidate(r.Context(), candidateData.ID, candidateData.Status, candidateData.ProdiID)
	if err != nil {
		slog.Error("failed to count unread announcements", "error", err)
		unreadCount = 0
	}

	// Convert to template type
	announcements := make([]portal.PortalAnnouncementItem, len(dbAnnouncements))
	for i, a := range dbAnnouncements {
		publishedAt := ""
		if a.PublishedAt != nil {
			publishedAt = a.PublishedAt.Format("02 January 2006")
		}
		announcements[i] = portal.PortalAnnouncementItem{
			ID:          a.ID,
			Title:       a.Title,
			Content:     a.Content,
			PublishedAt: publishedAt,
			IsRead:      a.IsRead,
		}
	}

	data := NewPageData("Pengumuman")
	portal.Announcements(data, announcements, unreadCount).Render(r.Context(), w)
}

func (h *PortalHandler) handleMarkAnnouncementRead(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	announcementID := r.PathValue("id")

	err := model.MarkAnnouncementRead(r.Context(), announcementID, claims.CandidateID)
	if err != nil {
		slog.Error("failed to mark announcement as read", "error", err, "announcement_id", announcementID)
		http.Error(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}

	// Return the updated card (now marked as read)
	ann, err := model.FindAnnouncementByID(r.Context(), announcementID)
	if err != nil || ann == nil {
		http.Error(w, "Announcement not found", http.StatusNotFound)
		return
	}

	publishedAt := ""
	if ann.PublishedAt != nil {
		publishedAt = ann.PublishedAt.Format("02 January 2006")
	}

	item := portal.PortalAnnouncementItem{
		ID:          ann.ID,
		Title:       ann.Title,
		Content:     ann.Content,
		PublishedAt: publishedAt,
		IsRead:      true, // Now marked as read
	}

	portal.AnnouncementCard(item).Render(r.Context(), w)
}

func (h *PortalHandler) handleReferral(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := NewPageData("Referral")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Referral - Coming Soon</h1></body></html>`))
	_ = data
}

// Helper functions

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getStatusLabel(status string) string {
	switch status {
	case "registered":
		return "Terdaftar"
	case "prospecting":
		return "Dalam Proses"
	case "committed":
		return "Berkomitmen"
	case "enrolled":
		return "Diterima"
	case "lost":
		return "Tidak Lanjut"
	default:
		return status
	}
}

func buildChecklist(c *model.CandidateDashboardData) []portal.ChecklistItem {
	checklist := []portal.ChecklistItem{}

	// 1. Email Verification (if email provided)
	if c.Email != nil && *c.Email != "" {
		status := "pending"
		if c.EmailVerified {
			status = "completed"
		}
		checklist = append(checklist, portal.ChecklistItem{
			Title:     "Verifikasi Email",
			Status:    status,
			Action:    "Verifikasi",
			ActionURL: "/portal/verify-email",
		})
	}

	// 2. Phone Verification (if phone provided)
	if c.Phone != nil && *c.Phone != "" {
		status := "pending"
		if c.PhoneVerified {
			status = "completed"
		}
		checklist = append(checklist, portal.ChecklistItem{
			Title:     "Verifikasi WhatsApp",
			Status:    status,
			Action:    "Verifikasi",
			ActionURL: "/portal/verify-phone",
		})
	}

	// 3. Personal Info
	personalComplete := c.Name != nil && *c.Name != "" && c.Address != nil && *c.Address != ""
	status := "pending"
	if personalComplete {
		status = "completed"
	}
	checklist = append(checklist, portal.ChecklistItem{
		Title:     "Lengkapi Data Diri",
		Status:    status,
		Action:    "Lengkapi",
		ActionURL: "/register?step=personal",
	})

	// 4. Education Info
	educationComplete := c.HighSchool != nil && *c.HighSchool != "" && c.ProdiID != nil
	status = "locked"
	if personalComplete {
		status = "pending"
		if educationComplete {
			status = "completed"
		}
	}
	checklist = append(checklist, portal.ChecklistItem{
		Title:     "Data Pendidikan",
		Status:    status,
		Action:    "Lengkapi",
		ActionURL: "/register?step=education",
	})

	// 5. Upload Documents
	status = "locked"
	if educationComplete {
		status = "pending"
	}
	checklist = append(checklist, portal.ChecklistItem{
		Title:     "Upload Dokumen",
		Status:    status,
		Action:    "Upload",
		ActionURL: "/portal/documents",
	})

	// 6. Payment (locked until documents complete)
	checklist = append(checklist, portal.ChecklistItem{
		Title:     "Pembayaran",
		Status:    "locked",
		Action:    "",
		ActionURL: "",
	})

	return checklist
}

func truncateContent(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
