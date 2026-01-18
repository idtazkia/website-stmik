package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/storage"
	"github.com/idtazkia/stmik-admission-api/templates/portal"
)

// PortalHandler handles all portal routes for candidates
type PortalHandler struct {
	sessionMgr *auth.SessionManager
	storage    storage.Storage
}

// NewPortalHandler creates a new portal handler
func NewPortalHandler(sessionMgr *auth.SessionManager, store storage.Storage) *PortalHandler {
	return &PortalHandler{sessionMgr: sessionMgr, storage: store}
}

// RegisterRoutes registers all portal routes to the mux
func (h *PortalHandler) RegisterRoutes(mux *http.ServeMux) {
	// All portal routes require candidate authentication
	mux.Handle("GET /portal", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleDashboard)))
	mux.Handle("GET /portal/documents", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleDocuments)))
	mux.Handle("POST /portal/documents/upload", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleDocumentUpload)))
	mux.Handle("GET /portal/payments", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handlePayments)))
	mux.Handle("POST /portal/payments/{id}/proof", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handlePaymentUpload)))
	mux.Handle("GET /portal/announcements", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleAnnouncements)))
	mux.Handle("POST /portal/announcements/{id}/read", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleMarkAnnouncementRead)))
	mux.Handle("GET /portal/referral", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleReferral)))
	mux.Handle("GET /portal/verify-email", RequireCandidateAuth(h.sessionMgr, http.HandlerFunc(h.handleVerifyEmailPage)))
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

	ctx := r.Context()
	candidateData, err := model.GetCandidateDashboardData(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get document types (active only)
	docTypes, err := model.ListDocumentTypes(ctx, true)
	if err != nil {
		slog.Error("failed to list document types", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get uploaded documents
	uploadedDocs, err := model.ListDocumentsByCandidate(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to list documents", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a map of uploaded documents by type ID
	uploadedByType := make(map[string]model.DocumentWithType)
	for _, d := range uploadedDocs {
		uploadedByType[d.DocumentTypeID] = d
	}

	// Get document stats
	uploaded, _, total, err := model.GetDocumentStats(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to get document stats", "error", err)
		uploaded, total = 0, len(docTypes)
	}

	data := NewPageData("Dokumen Saya")

	progressPercent := 0
	if total > 0 {
		progressPercent = (uploaded * 100) / total
	}

	candidate := portal.PortalCandidate{
		ID:                      candidateData.ID,
		Name:                    safeString(candidateData.Name),
		Status:                  candidateData.Status,
		DocumentProgress:        fmt.Sprintf("%d dari %d dokumen", uploaded, total),
		DocumentProgressPercent: fmt.Sprintf("%d", progressPercent),
	}

	// Build document items
	documents := make([]portal.DocumentItem, len(docTypes))
	for i, dt := range docTypes {
		allowedFormats := strings.ToUpper(strings.Join(dt.AllowedExtensions, ", "))
		allowedMimeTypes := buildMimeTypes(dt.AllowedExtensions)

		doc := portal.DocumentItem{
			Type:             dt.Code,
			Name:             dt.Name,
			Description:      safeString(dt.Description),
			Required:         dt.IsRequired,
			CanDefer:         dt.CanDefer,
			AllowedFormats:   allowedFormats,
			AllowedMimeTypes: allowedMimeTypes,
			MaxSize:          fmt.Sprintf("%d MB", dt.MaxFileSizeMB),
			Status:           "not_uploaded",
		}

		// Check if document is uploaded
		if uploaded, ok := uploadedByType[dt.ID]; ok {
			if uploaded.IsDeferred {
				doc.Status = "deferred"
			} else {
				doc.Status = uploaded.Status
				doc.FileName = uploaded.FileName
				doc.FileURL = h.storage.GetURL(uploaded.FilePath)
				doc.UploadedAt = uploaded.CreatedAt.Format("02 Jan 2006")
				if uploaded.RejectionReason != nil {
					doc.RejectionReason = *uploaded.RejectionReason
				}
			}
		}

		documents[i] = doc
	}

	portal.Documents(data, candidate, documents).Render(ctx, w)
}

func (h *PortalHandler) handleDocumentUpload(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ctx := r.Context()

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	documentType := r.FormValue("document_type")
	if documentType == "" {
		http.Error(w, "Document type is required", http.StatusBadRequest)
		return
	}

	// Get document type config
	docType, err := model.FindDocumentTypeByCode(ctx, documentType)
	if err != nil {
		slog.Error("failed to find document type", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if docType == nil {
		http.Error(w, "Invalid document type", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error("failed to get file", "error", err)
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	maxSize := int64(docType.MaxFileSizeMB) * 1024 * 1024
	if header.Size > maxSize {
		http.Error(w, fmt.Sprintf("File too large. Maximum size is %d MB", docType.MaxFileSizeMB), http.StatusBadRequest)
		return
	}

	// Validate file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(header.Filename), "."))
	validExt := false
	for _, allowed := range docType.AllowedExtensions {
		if ext == allowed {
			validExt = true
			break
		}
	}
	if !validExt {
		http.Error(w, fmt.Sprintf("Invalid file type. Allowed: %s", strings.Join(docType.AllowedExtensions, ", ")), http.StatusBadRequest)
		return
	}

	// Generate storage path: documents/{candidate_id}/{document_type}_{timestamp}.{ext}
	timestamp := time.Now().Unix()
	storagePath := fmt.Sprintf("documents/%s/%s_%d.%s", claims.CandidateID, documentType, timestamp, ext)

	// Upload file to storage
	if err := h.storage.Upload(ctx, storagePath, file); err != nil {
		slog.Error("failed to upload file", "error", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Get MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getMimeType(ext)
	}

	// Save document record
	_, err = model.CreateDocument(ctx, claims.CandidateID, docType.ID, header.Filename, storagePath, int(header.Size), mimeType)
	if err != nil {
		slog.Error("failed to create document record", "error", err)
		// Try to clean up uploaded file
		_ = h.storage.Delete(ctx, storagePath)
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	slog.Info("document uploaded", "candidate_id", claims.CandidateID, "document_type", documentType, "file", header.Filename)

	// Redirect back to documents page
	http.Redirect(w, r, "/portal/documents", http.StatusSeeOther)
}

func buildMimeTypes(extensions []string) string {
	mimeTypes := make([]string, len(extensions))
	for i, ext := range extensions {
		mimeTypes[i] = getMimeType(ext)
	}
	return strings.Join(mimeTypes, ",")
}

func getMimeType(ext string) string {
	switch strings.ToLower(ext) {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func (h *PortalHandler) handlePayments(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ctx := r.Context()
	candidateData, err := model.GetCandidateDashboardData(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if candidateData == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Fetch billings with payment info
	billings, err := model.ListBillingsByCandidate(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to list billings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch summary
	totalDue, totalPaid, totalPending, err := model.GetBillingSummary(ctx, claims.CandidateID)
	if err != nil {
		slog.Error("failed to get billing summary", "error", err)
		totalDue, totalPaid, totalPending = 0, 0, 0
	}

	data := NewPageData("Pembayaran")
	candidate := portal.PortalCandidate{
		ID:   candidateData.ID,
		Name: safeString(candidateData.Name),
	}

	// Convert billings to template payment items
	payments := make([]portal.PaymentItem, len(billings))
	for i, b := range billings {
		// Determine status for display
		status := b.Status
		if status == "pending_verification" {
			status = "pending"
		}

		dueDate := ""
		if b.DueDate != nil {
			dueDate = b.DueDate.Format("02 Jan 2006")
		}

		paidAt := ""
		if b.PaidAt != nil {
			paidAt = b.PaidAt.Format("02 Jan 2006")
		}

		proofURL := ""
		if b.PaymentProofURL != nil {
			proofURL = h.storage.GetURL(*b.PaymentProofURL)
		}

		payments[i] = portal.PaymentItem{
			ID:          b.ID,
			Type:        model.BillingTypeLabel(b.BillingType),
			Description: safeString(b.Description),
			Amount:      formatRupiah(int64(b.Amount)),
			DueDate:     dueDate,
			Status:      status,
			PaidAt:      paidAt,
			ProofURL:    proofURL,
		}
	}

	summary := portal.PaymentSummary{
		TotalDue:     formatRupiah(int64(totalDue)),
		TotalPaid:    formatRupiah(int64(totalPaid)),
		TotalPending: formatRupiah(int64(totalPending)),
	}

	portal.Payments(data, candidate, payments, summary).Render(ctx, w)
}

func (h *PortalHandler) handlePaymentUpload(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ctx := r.Context()
	billingID := r.PathValue("id")

	// Verify billing belongs to this candidate
	billing, err := model.FindBillingByID(ctx, billingID)
	if err != nil {
		slog.Error("failed to find billing", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if billing == nil || billing.CandidateID != claims.CandidateID {
		http.Error(w, "Billing not found", http.StatusNotFound)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Get transfer date
	transferDateStr := r.FormValue("transfer_date")
	if transferDateStr == "" {
		http.Error(w, "Transfer date is required", http.StatusBadRequest)
		return
	}
	transferDate, err := time.Parse("2006-01-02", transferDateStr)
	if err != nil {
		http.Error(w, "Invalid transfer date format", http.StatusBadRequest)
		return
	}

	// Get transfer amount
	amountStr := r.FormValue("amount")
	if amountStr == "" {
		http.Error(w, "Transfer amount is required", http.StatusBadRequest)
		return
	}
	// Parse amount (remove "Rp" and dots)
	amountStr = strings.ReplaceAll(amountStr, "Rp", "")
	amountStr = strings.ReplaceAll(amountStr, ".", "")
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.TrimSpace(amountStr)
	var amount int
	_, err = fmt.Sscanf(amountStr, "%d", &amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("proof")
	if err != nil {
		slog.Error("failed to get file", "error", err)
		http.Error(w, "Proof file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size (max 5MB)
	maxSize := int64(5) * 1024 * 1024
	if header.Size > maxSize {
		http.Error(w, "File too large. Maximum size is 5 MB", http.StatusBadRequest)
		return
	}

	// Validate file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(header.Filename), "."))
	validExts := []string{"jpg", "jpeg", "png", "pdf"}
	validExt := false
	for _, allowed := range validExts {
		if ext == allowed {
			validExt = true
			break
		}
	}
	if !validExt {
		http.Error(w, "Invalid file type. Allowed: JPG, PNG, PDF", http.StatusBadRequest)
		return
	}

	// Generate storage path: payments/{candidate_id}/{billing_id}_{timestamp}.{ext}
	timestamp := time.Now().Unix()
	storagePath := fmt.Sprintf("payments/%s/%s_%d.%s", claims.CandidateID, billingID, timestamp, ext)

	// Upload file to storage
	if err := h.storage.Upload(ctx, storagePath, file); err != nil {
		slog.Error("failed to upload file", "error", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Get MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getMimeType(ext)
	}

	// Create payment record
	_, err = model.CreatePayment(ctx, billingID, amount, transferDate, storagePath, header.Filename, int(header.Size), mimeType)
	if err != nil {
		slog.Error("failed to create payment record", "error", err)
		// Try to clean up uploaded file
		_ = h.storage.Delete(ctx, storagePath)
		http.Error(w, "Failed to save payment proof", http.StatusInternalServerError)
		return
	}

	slog.Info("payment proof uploaded", "candidate_id", claims.CandidateID, "billing_id", billingID, "file", header.Filename)

	// Redirect back to payments page
	http.Redirect(w, r, "/portal/payments", http.StatusSeeOther)
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

func (h *PortalHandler) handleVerifyEmailPage(w http.ResponseWriter, r *http.Request) {
	claims := GetCandidateClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
	if err != nil || candidate == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if candidate.Email == nil || *candidate.Email == "" {
		http.Redirect(w, r, "/portal", http.StatusFound)
		return
	}

	data := NewPageData("Verifikasi Email")
	verifyData := portal.VerifyEmailData{
		Email:         *candidate.Email,
		EmailVerified: candidate.EmailVerified,
	}

	portal.VerifyEmail(data, verifyData).Render(r.Context(), w)
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
