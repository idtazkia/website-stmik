package handler

import (
	"net/http"

	"github.com/idtazkia/stmik-admission-api/templates/portal"
)

// PortalHandler handles all portal routes for candidates
type PortalHandler struct{}

// NewPortalHandler creates a new portal handler
func NewPortalHandler() *PortalHandler {
	return &PortalHandler{}
}

// RegisterRoutes registers all portal routes to the mux
func (h *PortalHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /portal", h.handleDashboard)
	mux.HandleFunc("GET /portal/register", h.handleRegistration)
	mux.HandleFunc("GET /portal/documents", h.handleDocuments)
	mux.HandleFunc("GET /portal/payments", h.handlePayments)
	mux.HandleFunc("GET /portal/announcements", h.handleAnnouncements)
	mux.HandleFunc("GET /portal/referral", h.handleReferral)
}

func (h *PortalHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Dashboard")
	candidate := portal.DashboardCandidate{
		Name:            "Putri Amelia",
		Email:           "putri@gmail.com",
		Phone:           "081234567891",
		Program:         "Teknik Informatika",
		Status:          "prospecting",
		StatusLabel:     "Dalam Proses",
		RegisteredAt:    "10 Januari 2026",
		ConsultantName:  "Siti Rahayu",
		ConsultantPhone: "6281234567890",
	}
	checklist := []portal.ChecklistItem{
		{Title: "Verifikasi Email", Status: "completed", Action: "", ActionURL: ""},
		{Title: "Verifikasi WhatsApp", Status: "completed", Action: "", ActionURL: ""},
		{Title: "Lengkapi Data Diri", Status: "completed", Action: "", ActionURL: ""},
		{Title: "Upload Dokumen", Status: "pending", Action: "Upload", ActionURL: "/portal/documents"},
		{Title: "Pembayaran", Status: "locked", Action: "", ActionURL: ""},
	}
	announcements := []portal.AnnouncementItem{
		{Title: "Promo Early Bird Diperpanjang!", Content: "Dapatkan gratis biaya pendaftaran hingga 28 Februari 2026.", Date: "15 Jan 2026", IsNew: true},
		{Title: "Jadwal Tes Masuk Gelombang 1", Content: "Tes masuk gelombang 1 akan dilaksanakan pada 1 Maret 2026.", Date: "12 Jan 2026", IsNew: false},
	}
	portal.Dashboard(data, candidate, checklist, announcements).Render(r.Context(), w)
}

func (h *PortalHandler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Pendaftaran")
	steps := []portal.RegistrationStep{
		{Number: "1", Title: "Akun", Status: "completed"},
		{Number: "2", Title: "Verifikasi", Status: "completed"},
		{Number: "3", Title: "Data Diri", Status: "current"},
		{Number: "4", Title: "Prodi", Status: "pending"},
	}
	programs := []portal.ProgramOption{
		{Code: "SI", Name: "Sistem Informasi", Fee: "Rp 7.500.000/semester"},
		{Code: "TI", Name: "Teknik Informatika", Fee: "Rp 8.000.000/semester"},
	}
	portal.Registration(data, steps, programs, "biodata").Render(r.Context(), w)
}

func (h *PortalHandler) handleDocuments(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Dokumen Saya")

	candidate := portal.PortalCandidate{
		ID:                      "2",
		Name:                    "Putri Amelia",
		Status:                  "prospecting",
		DocumentProgress:        "2 dari 4 dokumen",
		DocumentProgressPercent: "50",
	}

	documents := []portal.DocumentItem{
		{Type: "ktp", Name: "KTP", Description: "Kartu Tanda Penduduk yang masih berlaku", Required: true, CanDefer: false, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "approved", FileName: "ktp_putri.jpg", FileURL: "/uploads/ktp_putri.jpg", UploadedAt: "14 Jan 2026"},
		{Type: "photo", Name: "Pas Foto", Description: "Foto 3x4 dengan background merah atau biru", Required: true, CanDefer: false, AllowedFormats: "JPG, PNG", AllowedMimeTypes: "image/jpeg,image/png", MaxSize: "2 MB", Status: "pending", FileName: "foto_putri.jpg", FileURL: "/uploads/foto_putri.jpg", UploadedAt: "15 Jan 2026"},
		{Type: "ijazah", Name: "Ijazah", Description: "Ijazah SMA/SMK/sederajat", Required: true, CanDefer: true, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "not_uploaded"},
		{Type: "transcript", Name: "Transkrip Nilai", Description: "Transkrip nilai SMA/SMK/sederajat", Required: true, CanDefer: true, AllowedFormats: "JPG, PNG, PDF", AllowedMimeTypes: "image/jpeg,image/png,application/pdf", MaxSize: "5 MB", Status: "not_uploaded"},
	}

	portal.Documents(data, candidate, documents).Render(r.Context(), w)
}

func (h *PortalHandler) handlePayments(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Pembayaran")
	candidate := portal.PortalCandidate{
		ID:   "2",
		Name: "Putri Amelia",
	}
	payments := []portal.PaymentItem{
		{ID: "1", Type: "Biaya Pendaftaran", Description: "Biaya pendaftaran mahasiswa baru TA 2025/2026", Amount: "Rp 500.000", DueDate: "20 Jan 2026", Status: "paid", PaidAt: "12 Jan 2026"},
		{ID: "2", Type: "SPP Semester 1", Description: "SPP Teknik Informatika semester 1", Amount: "Rp 8.000.000", DueDate: "1 Mar 2026", Status: "unpaid"},
		{ID: "3", Type: "Asrama Semester 1", Description: "Biaya asrama semester 1 (cicilan 1 dari 2)", Amount: "Rp 6.000.000", DueDate: "1 Mar 2026", Status: "pending", ProofURL: "/uploads/bukti_asrama.jpg"},
	}
	summary := portal.PaymentSummary{
		TotalDue:     "Rp 14.500.000",
		TotalPaid:    "Rp 500.000",
		TotalPending: "Rp 6.000.000",
	}
	portal.Payments(data, candidate, payments, summary).Render(r.Context(), w)
}

func (h *PortalHandler) handleAnnouncements(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Pengumuman")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Pengumuman - Coming Soon</h1></body></html>`))
	_ = data
}

func (h *PortalHandler) handleReferral(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Referral")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Referral - Coming Soon</h1></body></html>`))
	_ = data
}
