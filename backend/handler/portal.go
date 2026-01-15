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
	mux.HandleFunc("GET /portal/documents", h.handleDocuments)
	mux.HandleFunc("GET /portal/payments", h.handlePayments)
	mux.HandleFunc("GET /portal/announcements", h.handleAnnouncements)
	mux.HandleFunc("GET /portal/referral", h.handleReferral)
}

func (h *PortalHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Dashboard")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Portal Dashboard - Coming Soon</h1></body></html>`))
	_ = data
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
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Pembayaran - Coming Soon</h1></body></html>`))
	_ = data
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
