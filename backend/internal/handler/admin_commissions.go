package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleCommissionsReal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Komisi")

	// Parse filter parameters
	statusFilter := r.URL.Query().Get("status")
	referrerFilter := r.URL.Query().Get("referrer_id")

	filters := model.CommissionFilters{
		Status:     statusFilter,
		ReferrerID: referrerFilter,
		Limit:      50,
	}

	// Fetch commissions
	commissions, _, err := model.ListCommissions(ctx, filters)
	if err != nil {
		log.Printf("Error listing commissions: %v", err)
		http.Error(w, "Failed to load commissions", http.StatusInternalServerError)
		return
	}

	// Fetch stats
	stats, err := model.GetCommissionStats(ctx)
	if err != nil {
		log.Printf("Error getting commission stats: %v", err)
		http.Error(w, "Failed to load stats", http.StatusInternalServerError)
		return
	}

	// Convert to template types
	items := make([]admin.CommissionItem, len(commissions))
	for i, c := range commissions {
		candidateName := ""
		if c.CandidateName != nil {
			candidateName = *c.CandidateName
		}

		enrolledAt := c.CreatedAt.Format("2 Jan 2006")
		approvedAt := ""
		if c.ApprovedAt != nil {
			approvedAt = c.ApprovedAt.Format("2 Jan 2006")
		}
		paidAt := ""
		if c.PaidAt != nil {
			paidAt = c.PaidAt.Format("2 Jan 2006")
		}

		items[i] = admin.CommissionItem{
			ID:            c.ID,
			ReferrerName:  c.ReferrerName,
			ReferrerType:  c.ReferrerType,
			CandidateName: candidateName,
			CandidateNIM:  "", // NIM not implemented yet
			Amount:        formatCurrency(c.Amount),
			Status:        c.Status,
			EnrolledAt:    enrolledAt,
			ApprovedAt:    approvedAt,
			PaidAt:        paidAt,
			BankName:      "", // Need to fetch from referrer
			BankAccount:   "", // Need to fetch from referrer
		}
	}

	templateStats := admin.CommissionStats{
		Pending:        fmt.Sprintf("%d", stats.TotalPending),
		PendingAmount:  formatCurrency(stats.AmountPending),
		Approved:       fmt.Sprintf("%d", stats.TotalApproved),
		ApprovedAmount: formatCurrency(stats.AmountApproved),
		Paid:           fmt.Sprintf("%d", stats.TotalPaid),
		PaidAmount:     formatCurrency(stats.AmountPaid),
	}

	admin.Commissions(data, items, templateStats).Render(ctx, w)
}

func (h *AdminHandler) handleApproveCommission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and finance can approve
	if claims.Role != "admin" && claims.Role != "finance" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := model.ApproveCommission(ctx, id, claims.UserID)
	if err != nil {
		log.Printf("Error approving commission: %v", err)
		http.Error(w, "Failed to approve commission", http.StatusInternalServerError)
		return
	}

	log.Printf("Commission %s approved by %s", id, claims.UserID)

	// HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/commissions")
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) handleMarkCommissionPaid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and finance can mark paid
	if claims.Role != "admin" && claims.Role != "finance" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	notes := r.FormValue("notes")
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}

	err := model.MarkCommissionPaid(ctx, id, claims.UserID, notesPtr)
	if err != nil {
		log.Printf("Error marking commission paid: %v", err)
		http.Error(w, "Failed to mark commission paid", http.StatusInternalServerError)
		return
	}

	log.Printf("Commission %s marked paid by %s", id, claims.UserID)

	// HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/commissions")
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) handleBatchApproveCommissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and finance can approve
	if claims.Role != "admin" && claims.Role != "finance" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse IDs from form
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No commissions selected", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsStr, ",")

	count, err := model.BatchApproveCommissions(ctx, ids, claims.UserID)
	if err != nil {
		log.Printf("Error batch approving commissions: %v", err)
		http.Error(w, "Failed to approve commissions", http.StatusInternalServerError)
		return
	}

	log.Printf("Batch approved %d commissions by %s", count, claims.UserID)

	// HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/commissions")
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) handleBatchMarkCommissionsPaid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and finance can mark paid
	if claims.Role != "admin" && claims.Role != "finance" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse IDs from form
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No commissions selected", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsStr, ",")

	notes := r.FormValue("notes")
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}

	count, err := model.BatchMarkCommissionsPaid(ctx, ids, claims.UserID, notesPtr)
	if err != nil {
		log.Printf("Error batch marking commissions paid: %v", err)
		http.Error(w, "Failed to mark commissions paid", http.StatusInternalServerError)
		return
	}

	log.Printf("Batch marked %d commissions paid by %s", count, claims.UserID)

	// HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/commissions")
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) handleExportCommissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and finance can export
	if claims.Role != "admin" && claims.Role != "finance" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get status filter (default to approved for bank transfer)
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "approved"
	}

	// Fetch commissions with bank details
	commissions, err := model.ListCommissionsForExport(ctx, status)
	if err != nil {
		log.Printf("Error listing commissions for export: %v", err)
		http.Error(w, "Failed to export commissions", http.StatusInternalServerError)
		return
	}

	// Set headers for CSV download
	filename := fmt.Sprintf("commissions_%s_%s.csv", status, strings.ReplaceAll(strings.Split(r.Header.Get("Date"), " ")[0], "-", ""))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Write BOM for Excel UTF-8 compatibility
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Write CSV header
	w.Write([]byte("No,Nama Referrer,Tipe,Nama Bank,No Rekening,Atas Nama,Jumlah,Kandidat,Trigger Event,Tanggal Approve\n"))

	// Write data rows
	for i, c := range commissions {
		approvedAt := ""
		if c.ApprovedAt != nil {
			approvedAt = c.ApprovedAt.Format("2006-01-02")
		}

		row := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%d,%s,%s,%s\n",
			i+1,
			escapeCSV(c.ReferrerName),
			escapeCSV(c.ReferrerType),
			escapeCSV(c.BankName),
			escapeCSV(c.BankAccount),
			escapeCSV(c.AccountHolder),
			c.Amount,
			escapeCSV(c.CandidateName),
			escapeCSV(c.TriggerEvent),
			approvedAt,
		)
		w.Write([]byte(row))
	}

	log.Printf("Exported %d commissions (status=%s) by %s", len(commissions), status, claims.UserID)
}

// escapeCSV escapes a string for CSV format
func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

// formatCurrency formats amount as Indonesian Rupiah
func formatCurrency(amount int64) string {
	if amount == 0 {
		return "Rp 0"
	}

	// Format with thousand separators
	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return "Rp " + str
	}

	// Add dots every 3 digits from right
	var result strings.Builder
	for i, c := range str {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(c)
	}

	return "Rp " + result.String()
}
