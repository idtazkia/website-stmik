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

	// TODO: Fetch real announcements from database when announcement feature is implemented
	announcements := []portal.AnnouncementItem{
		{Title: "Selamat Datang!", Content: "Terima kasih telah mendaftar di STMIK Tazkia.", Date: candidateData.CreatedAt.Format("02 Jan 2006"), IsNew: true},
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

	data := NewPageData("Pengumuman")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Pengumuman - Coming Soon</h1></body></html>`))
	_ = data
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
