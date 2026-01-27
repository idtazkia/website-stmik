package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/integration"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

// FinanceHandler handles finance-related routes
type FinanceHandler struct {
	resend *integration.ResendClient
}

// NewFinanceHandler creates a new FinanceHandler
func NewFinanceHandler(resend *integration.ResendClient) *FinanceHandler {
	return &FinanceHandler{resend: resend}
}

// RegisterRoutes registers finance routes
func (h *FinanceHandler) RegisterRoutes(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	// Billing routes
	mux.Handle("GET /admin/finance/billings", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingList))))
	mux.Handle("GET /admin/finance/billings/create", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingCreateForm))))
	mux.Handle("POST /admin/finance/billings", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingCreate))))
	mux.Handle("GET /admin/finance/billings/{id}", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingDetail))))
	mux.Handle("POST /admin/finance/billings/{id}", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingUpdate))))
	mux.Handle("POST /admin/finance/billings/{id}/cancel", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleBillingCancel))))

	// Payment routes
	mux.Handle("GET /admin/finance/payments", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handlePaymentList))))
	mux.Handle("POST /admin/finance/payments/{id}/approve", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handlePaymentApprove))))
	mux.Handle("POST /admin/finance/payments/{id}/reject", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handlePaymentReject))))

	// API routes for candidate search
	mux.Handle("GET /admin/api/candidates/search", requireAuth(RequireFinanceOrAdmin(http.HandlerFunc(h.handleCandidateSearch))))
}

func (h *FinanceHandler) handleBillingList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query params
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	billingType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	filters := model.BillingFilters{
		Search:      search,
		Status:      status,
		BillingType: billingType,
		Page:        page,
		PageSize:    20,
	}

	// Get billings
	billings, total, err := model.ListAllBillings(ctx, filters)
	if err != nil {
		slog.Error("failed to list billings", "error", err)
		http.Error(w, "Failed to load billings", http.StatusInternalServerError)
		return
	}

	// Get stats
	unpaid, pending, paid, cancelled, err := model.GetBillingStats(ctx)
	if err != nil {
		slog.Error("failed to get billing stats", "error", err)
	}

	data := admin.BillingListData{
		Billings:    billings,
		Total:       total,
		Page:        page,
		PageSize:    20,
		Search:      search,
		Status:      status,
		BillingType: billingType,
		Stats: admin.BillingStats{
			Unpaid:    unpaid,
			Pending:   pending,
			Paid:      paid,
			Cancelled: cancelled,
		},
	}

	pageData := NewPageDataWithUser(ctx, "Manajemen Tagihan")
	admin.FinanceBillings(pageData, data).Render(ctx, w)
}

func (h *FinanceHandler) handleBillingCreateForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if candidate_id is provided (from candidate detail page)
	candidateID := r.URL.Query().Get("candidate_id")
	var formData admin.BillingFormData

	if candidateID != "" {
		// Get candidate info
		candidate, err := model.GetCandidateDetailData(ctx, candidateID)
		if err != nil {
			slog.Error("failed to get candidate", "error", err)
			http.Error(w, "Failed to load candidate", http.StatusInternalServerError)
			return
		}
		if candidate != nil {
			formData.CandidateID = candidateID
			if candidate.Name != nil {
				formData.CandidateName = *candidate.Name
			}
			if candidate.Email != nil {
				formData.CandidateEmail = *candidate.Email
			}
			if candidate.ProdiName != nil {
				formData.ProdiName = *candidate.ProdiName
			}
		}
	}

	pageData := NewPageDataWithUser(ctx, "Buat Tagihan Baru")
	admin.FinanceBillingForm(pageData, formData).Render(ctx, w)
}

func (h *FinanceHandler) handleBillingCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	candidateID := r.FormValue("candidate_id")
	billingType := r.FormValue("billing_type")
	amountStr := r.FormValue("amount")
	dueDateStr := r.FormValue("due_date")
	description := r.FormValue("description")

	// Validate
	if candidateID == "" {
		pageData := NewPageDataWithUser(ctx, "Buat Tagihan Baru")
		formData := admin.BillingFormData{Error: "Kandidat wajib dipilih"}
		admin.FinanceBillingForm(pageData, formData).Render(ctx, w)
		return
	}
	if billingType == "" {
		pageData := NewPageDataWithUser(ctx, "Buat Tagihan Baru")
		formData := admin.BillingFormData{
			CandidateID: candidateID,
			Error:       "Jenis tagihan wajib dipilih",
		}
		admin.FinanceBillingForm(pageData, formData).Render(ctx, w)
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		pageData := NewPageDataWithUser(ctx, "Buat Tagihan Baru")
		formData := admin.BillingFormData{
			CandidateID: candidateID,
			BillingType: billingType,
			Error:       "Jumlah harus berupa angka positif",
		}
		admin.FinanceBillingForm(pageData, formData).Render(ctx, w)
		return
	}

	var dueDate *time.Time
	if dueDateStr != "" {
		t, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			dueDate = &t
		}
	}

	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	// Create billing
	billing, err := model.CreateBilling(ctx, candidateID, billingType, descPtr, amount, dueDate)
	if err != nil {
		slog.Error("failed to create billing", "error", err)
		pageData := NewPageDataWithUser(ctx, "Buat Tagihan Baru")
		formData := admin.BillingFormData{
			CandidateID: candidateID,
			BillingType: billingType,
			Amount:      amount,
			Description: description,
			Error:       "Gagal membuat tagihan: " + err.Error(),
		}
		admin.FinanceBillingForm(pageData, formData).Render(ctx, w)
		return
	}

	slog.Info("billing created", "billing_id", billing.ID, "candidate_id", candidateID, "type", billingType, "amount", amount)

	http.Redirect(w, r, "/admin/finance/billings/"+billing.ID, http.StatusSeeOther)
}

func (h *FinanceHandler) handleBillingDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	billingID := r.PathValue("id")

	// Get billing
	billing, err := model.FindBillingByID(ctx, billingID)
	if err != nil {
		slog.Error("failed to find billing", "error", err)
		http.Error(w, "Failed to load billing", http.StatusInternalServerError)
		return
	}
	if billing == nil {
		http.NotFound(w, r)
		return
	}

	// Get candidate info
	candidate, err := model.GetCandidateDetailData(ctx, billing.CandidateID)
	if err != nil {
		slog.Error("failed to get candidate", "error", err)
	}

	// Get payments
	payments, err := model.ListPaymentsByBilling(ctx, billingID)
	if err != nil {
		slog.Error("failed to get payments", "error", err)
	}

	// Convert to template type
	paymentInfos := make([]admin.PaymentInfo, len(payments))
	for i, p := range payments {
		paymentInfos[i] = admin.PaymentInfo{
			ID:           p.ID,
			Amount:       p.Amount,
			TransferDate: &p.TransferDate,
			ProofURL:     p.ProofFilePath,
			Status:       p.Status,
			RejectedNote: p.RejectionReason,
			ReviewedBy:   p.ReviewedBy,
			ReviewedAt:   p.ReviewedAt,
			CreatedAt:    p.CreatedAt,
		}
	}

	data := admin.BillingDetailData{
		Billing:  *billing,
		Payments: paymentInfos,
	}

	if candidate != nil {
		if candidate.Name != nil {
			data.CandidateName = *candidate.Name
		}
		if candidate.Email != nil {
			data.CandidateEmail = *candidate.Email
		}
		if candidate.ProdiName != nil {
			data.ProdiName = *candidate.ProdiName
		}
	}

	// Check for success message
	if r.URL.Query().Get("success") == "1" {
		data.Success = "Tagihan berhasil diperbarui"
	}

	pageData := NewPageDataWithUser(ctx, "Detail Tagihan")
	admin.FinanceBillingDetail(pageData, data).Render(ctx, w)
}

func (h *FinanceHandler) handleBillingUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	billingID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	dueDateStr := r.FormValue("due_date")
	description := r.FormValue("description")

	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		http.Error(w, "Jumlah harus berupa angka positif", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if dueDateStr != "" {
		t, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			dueDate = &t
		}
	}

	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	if err := model.UpdateBilling(ctx, billingID, amount, dueDate, descPtr); err != nil {
		slog.Error("failed to update billing", "error", err)
		http.Error(w, "Gagal memperbarui tagihan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("billing updated", "billing_id", billingID, "amount", amount)

	http.Redirect(w, r, "/admin/finance/billings/"+billingID+"?success=1", http.StatusSeeOther)
}

func (h *FinanceHandler) handleBillingCancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	billingID := r.PathValue("id")

	if err := model.CancelBilling(ctx, billingID); err != nil {
		slog.Error("failed to cancel billing", "error", err)
		http.Error(w, "Gagal membatalkan tagihan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("billing cancelled", "billing_id", billingID)

	http.Redirect(w, r, "/admin/finance/billings", http.StatusSeeOther)
}

func (h *FinanceHandler) handlePaymentList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get payments
	payments, total, err := model.ListPaymentsWithDetails(ctx, status, page, 20)
	if err != nil {
		slog.Error("failed to list payments", "error", err)
		http.Error(w, "Failed to load payments", http.StatusInternalServerError)
		return
	}

	// Get stats
	pending, approved, rejected, err := model.GetPaymentStats(ctx)
	if err != nil {
		slog.Error("failed to get payment stats", "error", err)
	}

	// Convert to template type
	paymentList := make([]admin.PaymentWithDetails, len(payments))
	for i, p := range payments {
		name := ""
		if p.CandidateName != "" {
			name = p.CandidateName
		}
		paymentList[i] = admin.PaymentWithDetails{
			ID:             p.ID,
			BillingID:      p.BillingID,
			BillingType:    p.BillingType,
			BillingAmount:  p.BillingAmount,
			CandidateID:    p.CandidateID,
			CandidateName:  name,
			CandidateEmail: "", // Not in PaymentWithBilling, could add if needed
			Amount:         p.Amount,
			TransferDate:   &p.TransferDate,
			ProofURL:       p.ProofFilePath,
			Status:         p.Status,
			RejectedNote:   p.RejectionReason,
			ReviewedBy:     p.ReviewedBy,
			ReviewedAt:     p.ReviewedAt,
			CreatedAt:      p.CreatedAt,
		}
	}

	data := admin.PaymentListData{
		Payments: paymentList,
		Total:    total,
		Page:     page,
		PageSize: 20,
		Status:   status,
		Stats: admin.PaymentStats{
			Pending:  pending,
			Approved: approved,
			Rejected: rejected,
		},
	}

	pageData := NewPageDataWithUser(ctx, "Verifikasi Pembayaran")
	admin.FinancePayments(pageData, data).Render(ctx, w)
}

func (h *FinanceHandler) handlePaymentApprove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paymentID := r.PathValue("id")

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get payment details before approving (for email)
	payment, err := model.FindPaymentByID(ctx, paymentID)
	if err != nil || payment == nil {
		slog.Error("failed to find payment", "error", err)
		http.Error(w, "Pembayaran tidak ditemukan", http.StatusNotFound)
		return
	}

	// Get billing and candidate info
	billing, err := model.FindBillingByID(ctx, payment.BillingID)
	if err != nil || billing == nil {
		slog.Error("failed to find billing", "error", err)
		http.Error(w, "Tagihan tidak ditemukan", http.StatusNotFound)
		return
	}

	candidate, err := model.FindCandidateByID(ctx, billing.CandidateID)
	if err != nil {
		slog.Error("failed to find candidate", "error", err)
	}

	// Approve the payment
	if err := model.ApprovePayment(ctx, paymentID, claims.UserID); err != nil {
		slog.Error("failed to approve payment", "error", err)
		http.Error(w, "Gagal menyetujui pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("payment approved", "payment_id", paymentID, "reviewer", claims.UserID)

	// Send confirmation email (non-blocking, log errors but don't fail the request)
	if h.resend != nil && candidate != nil && candidate.Email != nil && *candidate.Email != "" {
		candidateName := ""
		if candidate.Name != nil {
			candidateName = *candidate.Name
		}

		emailData := integration.PaymentConfirmationData{
			CandidateName: candidateName,
			BillingType:   model.BillingTypeLabel(billing.BillingType),
			Amount:        formatRupiah(int64(payment.Amount)),
			TransferDate:  payment.TransferDate.Format("02 January 2006"),
			ApprovedAt:    time.Now().Format("02 January 2006 15:04"),
		}

		go func() {
			if err := h.resend.SendPaymentConfirmation(*candidate.Email, emailData); err != nil {
				slog.Error("failed to send payment confirmation email", "error", err, "email", *candidate.Email)
			} else {
				slog.Info("payment confirmation email sent", "email", *candidate.Email, "payment_id", paymentID)
			}
		}()
	}

	// For HTMX, return updated card
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/finance/payments")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/finance/payments", http.StatusSeeOther)
}

func (h *FinanceHandler) handlePaymentReject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paymentID := r.PathValue("id")

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	reason := r.FormValue("reason")
	if reason == "" {
		http.Error(w, "Alasan penolakan wajib diisi", http.StatusBadRequest)
		return
	}

	// Get payment details before rejecting (for email)
	payment, err := model.FindPaymentByID(ctx, paymentID)
	if err != nil || payment == nil {
		slog.Error("failed to find payment", "error", err)
		http.Error(w, "Pembayaran tidak ditemukan", http.StatusNotFound)
		return
	}

	// Get billing and candidate info
	billing, err := model.FindBillingByID(ctx, payment.BillingID)
	if err != nil || billing == nil {
		slog.Error("failed to find billing", "error", err)
		http.Error(w, "Tagihan tidak ditemukan", http.StatusNotFound)
		return
	}

	candidate, err := model.FindCandidateByID(ctx, billing.CandidateID)
	if err != nil {
		slog.Error("failed to find candidate", "error", err)
	}

	// Reject the payment
	if err := model.RejectPayment(ctx, paymentID, claims.UserID, reason); err != nil {
		slog.Error("failed to reject payment", "error", err)
		http.Error(w, "Gagal menolak pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("payment rejected", "payment_id", paymentID, "reviewer", claims.UserID, "reason", reason)

	// Send rejection email (non-blocking, log errors but don't fail the request)
	if h.resend != nil && candidate != nil && candidate.Email != nil && *candidate.Email != "" {
		candidateName := ""
		if candidate.Name != nil {
			candidateName = *candidate.Name
		}

		emailData := integration.PaymentRejectionData{
			CandidateName: candidateName,
			BillingType:   model.BillingTypeLabel(billing.BillingType),
			Amount:        formatRupiah(int64(payment.Amount)),
			TransferDate:  payment.TransferDate.Format("02 January 2006"),
			Reason:        reason,
		}

		go func() {
			if err := h.resend.SendPaymentRejection(*candidate.Email, emailData); err != nil {
				slog.Error("failed to send payment rejection email", "error", err, "email", *candidate.Email)
			} else {
				slog.Info("payment rejection email sent", "email", *candidate.Email, "payment_id", paymentID)
			}
		}()
	}

	// For HTMX, return updated card
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/finance/payments")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/finance/payments", http.StatusSeeOther)
}

func (h *FinanceHandler) handleCandidateSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")

	if query == "" || len(query) < 2 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	// Search candidates by name or email
	result, err := model.ListCandidates(ctx, model.CandidateListFilters{
		Search: query,
		Limit:  10,
		Offset: 0,
	}, nil, nil)
	if err != nil {
		slog.Error("failed to search candidates", "error", err)
		http.Error(w, "Failed to search candidates", http.StatusInternalServerError)
		return
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	for i, c := range result.Candidates {
		if i > 0 {
			w.Write([]byte(","))
		}
		name := ""
		if c.Name != nil {
			name = *c.Name
		}
		email := ""
		if c.Email != nil {
			email = *c.Email
		}
		w.Write([]byte(`{"id":"` + c.ID + `","email":"` + email + `","name":"` + name + `"}`))
	}
	w.Write([]byte("]"))
}
