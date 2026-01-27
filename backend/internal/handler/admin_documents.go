package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/idtazkia/stmik-admission-api/internal/integration"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleDocumentReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Review Dokumen")

	// Parse filters from query string
	filters := model.DocumentReviewFilters{
		Status: r.URL.Query().Get("status"),
		Type:   r.URL.Query().Get("type"),
		Search: r.URL.Query().Get("search"),
	}

	// Default to pending status if no filter
	if filters.Status == "" {
		filters.Status = "pending"
	}

	// Fetch documents
	docs, err := model.ListDocumentsForReview(ctx, filters)
	if err != nil {
		slog.Error("failed to list documents for review", "error", err)
		http.Error(w, "Failed to load documents", http.StatusInternalServerError)
		return
	}

	// Fetch stats
	stats, err := model.GetDocumentReviewStats(ctx)
	if err != nil {
		slog.Error("failed to get document review stats", "error", err)
		stats = &model.DocumentReviewStats{}
	}

	// Convert to template types
	documentItems := make([]admin.DocumentReviewItem, len(docs))
	for i, d := range docs {
		rejectionReason := ""
		if d.RejectionReason != nil {
			rejectionReason = *d.RejectionReason
		}

		isImage := strings.HasPrefix(d.MimeType, "image/")
		fileURL := "/uploads/" + d.FilePath
		thumbnailURL := fileURL // For images, use the same URL

		documentItems[i] = admin.DocumentReviewItem{
			ID:              d.ID,
			CandidateID:     d.CandidateID,
			CandidateName:   d.CandidateName,
			ProdiName:       d.ProdiName,
			Type:            d.TypeCode,
			TypeName:        d.TypeName,
			FileName:        d.FileName,
			FileSize:        formatFileSize(d.FileSize),
			FileURL:         fileURL,
			ThumbnailURL:    thumbnailURL,
			IsImage:         isImage,
			Status:          d.Status,
			RejectionReason: rejectionReason,
			UploadedAt:      d.CreatedAt.Format("2 Jan 2006 15:04"),
		}
	}

	// Build template data
	templateFilter := admin.DocumentFilter{
		Status: filters.Status,
		Type:   filters.Type,
		Search: filters.Search,
	}

	templateStats := admin.DocumentStats{
		Pending:       fmt.Sprintf("%d", stats.Pending),
		ApprovedToday: fmt.Sprintf("%d", stats.ApprovedToday),
		RejectedToday: fmt.Sprintf("%d", stats.RejectedToday),
		Total:         fmt.Sprintf("%d", stats.Total),
	}

	admin.DocumentReviewList(data, templateFilter, documentItems, templateStats).Render(ctx, w)
}

func (h *AdminHandler) handleApproveDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	documentID := r.FormValue("document_id")
	if documentID == "" {
		http.Error(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	// Get document details before approving (for email)
	doc, err := model.FindDocumentByID(ctx, documentID)
	if err != nil || doc == nil {
		slog.Error("failed to find document", "error", err, "document_id", documentID)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Get candidate info
	candidate, err := model.FindCandidateByID(ctx, doc.CandidateID)
	if err != nil {
		slog.Error("failed to find candidate", "error", err)
	}

	// Get document type name
	docType, err := model.FindDocumentTypeByID(ctx, doc.DocumentTypeID)
	if err != nil {
		slog.Error("failed to find document type", "error", err)
	}

	// Approve the document
	err = model.ApproveDocument(ctx, documentID, claims.UserID)
	if err != nil {
		slog.Error("failed to approve document", "error", err, "document_id", documentID)
		http.Error(w, "Failed to approve document", http.StatusInternalServerError)
		return
	}

	slog.Info("document approved", "document_id", documentID, "reviewer_id", claims.UserID)

	// Send approval email (non-blocking)
	if h.resend != nil && candidate != nil && candidate.Email != nil && *candidate.Email != "" {
		candidateName := ""
		if candidate.Name != nil {
			candidateName = *candidate.Name
		}
		docTypeName := "Dokumen"
		if docType != nil {
			docTypeName = docType.Name
		}

		emailData := integration.DocumentStatusData{
			CandidateName: candidateName,
			DocumentType:  docTypeName,
			Status:        "approved",
		}

		go func() {
			if err := h.resend.SendDocumentApproved(*candidate.Email, emailData); err != nil {
				slog.Error("failed to send document approval email", "error", err, "email", *candidate.Email)
			} else {
				slog.Info("document approval email sent", "email", *candidate.Email, "document_id", documentID)
			}
		}()
	}

	// Redirect back to documents page
	http.Redirect(w, r, "/admin/documents", http.StatusSeeOther)
}

func (h *AdminHandler) handleRejectDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	documentID := r.FormValue("document_id")
	if documentID == "" {
		http.Error(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	rejectionReason := r.FormValue("rejection_reason")
	rejectionNotes := r.FormValue("rejection_notes")

	// Build full rejection message
	reason := rejectionReasonLabel(rejectionReason)
	if rejectionNotes != "" {
		reason = reason + ": " + rejectionNotes
	}

	// Get document details before rejecting (for email)
	doc, err := model.FindDocumentByID(ctx, documentID)
	if err != nil || doc == nil {
		slog.Error("failed to find document", "error", err, "document_id", documentID)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Get candidate info
	candidate, err := model.FindCandidateByID(ctx, doc.CandidateID)
	if err != nil {
		slog.Error("failed to find candidate", "error", err)
	}

	// Get document type name
	docType, err := model.FindDocumentTypeByID(ctx, doc.DocumentTypeID)
	if err != nil {
		slog.Error("failed to find document type", "error", err)
	}

	// Reject the document
	err = model.RejectDocument(ctx, documentID, claims.UserID, reason)
	if err != nil {
		slog.Error("failed to reject document", "error", err, "document_id", documentID)
		http.Error(w, "Failed to reject document", http.StatusInternalServerError)
		return
	}

	slog.Info("document rejected", "document_id", documentID, "reviewer_id", claims.UserID, "reason", reason)

	// Send rejection email (non-blocking)
	if h.resend != nil && candidate != nil && candidate.Email != nil && *candidate.Email != "" {
		candidateName := ""
		if candidate.Name != nil {
			candidateName = *candidate.Name
		}
		docTypeName := "Dokumen"
		if docType != nil {
			docTypeName = docType.Name
		}

		emailData := integration.DocumentStatusData{
			CandidateName: candidateName,
			DocumentType:  docTypeName,
			Status:        "rejected",
			Reason:        reason,
		}

		go func() {
			if err := h.resend.SendDocumentRejected(*candidate.Email, emailData); err != nil {
				slog.Error("failed to send document rejection email", "error", err, "email", *candidate.Email)
			} else {
				slog.Info("document rejection email sent", "email", *candidate.Email, "document_id", documentID)
			}
		}()
	}

	// Redirect back to documents page
	http.Redirect(w, r, "/admin/documents", http.StatusSeeOther)
}

func formatFileSize(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}

func rejectionReasonLabel(reason string) string {
	switch reason {
	case "blur":
		return "Gambar buram/tidak jelas"
	case "incomplete":
		return "Dokumen tidak lengkap/terpotong"
	case "wrong_document":
		return "Dokumen tidak sesuai"
	case "expired":
		return "Dokumen sudah tidak berlaku"
	case "wrong_format":
		return "Format file tidak sesuai"
	case "other":
		return "Lainnya"
	default:
		return reason
	}
}
