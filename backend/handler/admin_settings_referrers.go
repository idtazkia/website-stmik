package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

func (h *AdminHandler) handleReferrersSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Referrer")

	// Fetch referrers from database
	dbReferrers, err := model.ListReferrers(r.Context(), "")
	if err != nil {
		slog.Error("failed to list referrers", "error", err)
		http.Error(w, "Failed to load referrers", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	referrers := make([]admin.ReferrerItem, len(dbReferrers))
	for i, ref := range dbReferrers {
		institution := ""
		if ref.Institution != nil {
			institution = *ref.Institution
		}
		phone := ""
		if ref.Phone != nil {
			phone = *ref.Phone
		}
		email := ""
		if ref.Email != nil {
			email = *ref.Email
		}
		code := ""
		if ref.Code != nil {
			code = *ref.Code
		}
		bankName := ""
		if ref.BankName != nil {
			bankName = *ref.BankName
		}
		bankAccount := ""
		if ref.BankAccount != nil {
			bankAccount = *ref.BankAccount
		}
		accountHolder := ""
		if ref.AccountHolder != nil {
			accountHolder = *ref.AccountHolder
		}
		commissionStr := ""
		if ref.CommissionOverride != nil {
			commissionStr = formatRupiah(*ref.CommissionOverride)
		}

		referrers[i] = admin.ReferrerItem{
			ID:                 ref.ID,
			Name:               ref.Name,
			Type:               ref.Type,
			Institution:        institution,
			Phone:              phone,
			Email:              email,
			Code:               code,
			BankName:           bankName,
			BankAccount:        bankAccount,
			AccountHolder:      accountHolder,
			CommissionOverride: ref.CommissionOverride,
			CommissionStr:      commissionStr,
			PayoutPreference:   ref.PayoutPreference,
			IsActive:           ref.IsActive,
		}
	}

	// Fetch stats
	counts, err := model.CountReferrersByType(r.Context())
	if err != nil {
		slog.Error("failed to count referrers", "error", err)
	}

	stats := admin.ReferrerStats{
		Total:   counts["total"],
		Alumni:  counts["alumni"],
		Teacher: counts["teacher"],
		Student: counts["student"],
		Partner: counts["partner"],
		Staff:   counts["staff"],
	}

	admin.SettingsReferrers(data, referrers, stats).Render(r.Context(), w)
}

func (h *AdminHandler) handleCreateReferrer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	referrerType := r.FormValue("type")
	institutionStr := r.FormValue("institution")
	phoneStr := r.FormValue("phone")
	emailStr := r.FormValue("email")
	codeStr := r.FormValue("code")
	bankNameStr := r.FormValue("bank_name")
	bankAccountStr := r.FormValue("bank_account")
	accountHolderStr := r.FormValue("account_holder")
	commissionOverrideStr := r.FormValue("commission_override")
	payoutPreference := r.FormValue("payout_preference")

	if name == "" || referrerType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	if payoutPreference == "" {
		payoutPreference = "per_enrollment"
	}

	var institution, phone, email, code, bankName, bankAccount, accountHolder *string
	if institutionStr != "" {
		institution = &institutionStr
	}
	if phoneStr != "" {
		phone = &phoneStr
	}
	if emailStr != "" {
		email = &emailStr
	}
	if codeStr != "" {
		code = &codeStr
	} else {
		// Generate referral code if not provided
		generatedCode := model.GenerateReferralCode(name, referrerType)
		code = &generatedCode
	}
	if bankNameStr != "" {
		bankName = &bankNameStr
	}
	if bankAccountStr != "" {
		bankAccount = &bankAccountStr
	}
	if accountHolderStr != "" {
		accountHolder = &accountHolderStr
	}

	var commissionOverride *int64
	if commissionOverrideStr != "" {
		amt, err := strconv.ParseInt(commissionOverrideStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid commission amount", http.StatusBadRequest)
			return
		}
		commissionOverride = &amt
	}

	referrer, err := model.CreateReferrer(r.Context(), name, referrerType, institution, phone, email, code, bankName, bankAccount, accountHolder, commissionOverride, payoutPreference)
	if err != nil {
		slog.Error("failed to create referrer", "error", err)
		http.Error(w, "Failed to create referrer", http.StatusInternalServerError)
		return
	}

	slog.Info("referrer created", "referrer_id", referrer.ID)
	h.renderReferrerRow(w, r, referrer.ID)
}

func (h *AdminHandler) handleUpdateReferrer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing referrer ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	referrerType := r.FormValue("type")
	institutionStr := r.FormValue("institution")
	phoneStr := r.FormValue("phone")
	emailStr := r.FormValue("email")
	codeStr := r.FormValue("code")
	bankNameStr := r.FormValue("bank_name")
	bankAccountStr := r.FormValue("bank_account")
	accountHolderStr := r.FormValue("account_holder")
	commissionOverrideStr := r.FormValue("commission_override")
	payoutPreference := r.FormValue("payout_preference")

	if name == "" || referrerType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	if payoutPreference == "" {
		payoutPreference = "per_enrollment"
	}

	var institution, phone, email, code, bankName, bankAccount, accountHolder *string
	if institutionStr != "" {
		institution = &institutionStr
	}
	if phoneStr != "" {
		phone = &phoneStr
	}
	if emailStr != "" {
		email = &emailStr
	}
	if codeStr != "" {
		code = &codeStr
	}
	if bankNameStr != "" {
		bankName = &bankNameStr
	}
	if bankAccountStr != "" {
		bankAccount = &bankAccountStr
	}
	if accountHolderStr != "" {
		accountHolder = &accountHolderStr
	}

	var commissionOverride *int64
	if commissionOverrideStr != "" {
		amt, err := strconv.ParseInt(commissionOverrideStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid commission amount", http.StatusBadRequest)
			return
		}
		commissionOverride = &amt
	}

	err := model.UpdateReferrer(r.Context(), id, name, referrerType, institution, phone, email, code, bankName, bankAccount, accountHolder, commissionOverride, payoutPreference)
	if err != nil {
		slog.Error("failed to update referrer", "error", err)
		http.Error(w, "Failed to update referrer", http.StatusInternalServerError)
		return
	}

	slog.Info("referrer updated", "referrer_id", id)
	h.renderReferrerRow(w, r, id)
}

func (h *AdminHandler) handleToggleReferrerActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing referrer ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleReferrerActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle referrer active", "error", err)
		http.Error(w, "Failed to toggle referrer status", http.StatusInternalServerError)
		return
	}

	slog.Info("referrer active status toggled", "referrer_id", id)
	h.renderReferrerRow(w, r, id)
}

func (h *AdminHandler) renderReferrerRow(w http.ResponseWriter, r *http.Request, referrerID string) {
	referrer, err := model.FindReferrerByID(r.Context(), referrerID)
	if err != nil {
		slog.Error("failed to find referrer", "error", err)
		http.Error(w, "Failed to load referrer", http.StatusInternalServerError)
		return
	}

	if referrer == nil {
		http.NotFound(w, r)
		return
	}

	institution := ""
	if referrer.Institution != nil {
		institution = *referrer.Institution
	}
	phone := ""
	if referrer.Phone != nil {
		phone = *referrer.Phone
	}
	email := ""
	if referrer.Email != nil {
		email = *referrer.Email
	}
	code := ""
	if referrer.Code != nil {
		code = *referrer.Code
	}
	bankName := ""
	if referrer.BankName != nil {
		bankName = *referrer.BankName
	}
	bankAccount := ""
	if referrer.BankAccount != nil {
		bankAccount = *referrer.BankAccount
	}
	accountHolder := ""
	if referrer.AccountHolder != nil {
		accountHolder = *referrer.AccountHolder
	}
	commissionStr := ""
	if referrer.CommissionOverride != nil {
		commissionStr = formatRupiah(*referrer.CommissionOverride)
	}

	item := admin.ReferrerItem{
		ID:                 referrer.ID,
		Name:               referrer.Name,
		Type:               referrer.Type,
		Institution:        institution,
		Phone:              phone,
		Email:              email,
		Code:               code,
		BankName:           bankName,
		BankAccount:        bankAccount,
		AccountHolder:      accountHolder,
		CommissionOverride: referrer.CommissionOverride,
		CommissionStr:      commissionStr,
		PayoutPreference:   referrer.PayoutPreference,
		IsActive:           referrer.IsActive,
	}

	admin.ReferrerRow(item).Render(r.Context(), w)
}
