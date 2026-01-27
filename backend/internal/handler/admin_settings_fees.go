package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleFeesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Biaya")

	// Get academic year from query param or use current
	academicYear := r.URL.Query().Get("academic_year")
	if academicYear == "" {
		academicYear = "2025/2026"
	}

	// Fetch fee structures from database
	dbFees, err := model.ListFeeStructures(r.Context(), academicYear, false)
	if err != nil {
		slog.Error("failed to list fee structures", "error", err)
		http.Error(w, "Failed to load fee structures", http.StatusInternalServerError)
		return
	}

	// Fetch fee types for dropdown
	dbFeeTypes, err := model.ListFeeTypes(r.Context())
	if err != nil {
		slog.Error("failed to list fee types", "error", err)
		http.Error(w, "Failed to load fee types", http.StatusInternalServerError)
		return
	}

	// Fetch prodis for dropdown
	dbProdis, err := model.ListProdis(r.Context(), true) // only active prodis
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
		http.Error(w, "Failed to load programs", http.StatusInternalServerError)
		return
	}

	// Convert to template types
	feeStructures := make([]admin.FeeStructureItem, len(dbFees))
	for i, f := range dbFees {
		prodiID := ""
		prodiName := ""
		prodiCode := ""
		if f.IDProdi != nil {
			prodiID = *f.IDProdi
		}
		if f.ProdiName != nil {
			prodiName = *f.ProdiName
		}
		if f.ProdiCode != nil {
			prodiCode = *f.ProdiCode
		}
		feeStructures[i] = admin.FeeStructureItem{
			ID:           f.ID,
			IDFeeType:    f.IDFeeType,
			FeeTypeName:  f.FeeTypeName,
			FeeTypeCode:  f.FeeTypeCode,
			IDProdi:      prodiID,
			ProdiName:    prodiName,
			ProdiCode:    prodiCode,
			AcademicYear: f.AcademicYear,
			Amount:       f.Amount,
			AmountStr:    formatRupiah(f.Amount),
			IsActive:     f.IsActive,
		}
	}

	feeTypes := make([]admin.FeeTypeOption, len(dbFeeTypes))
	for i, ft := range dbFeeTypes {
		feeTypes[i] = admin.FeeTypeOption{
			ID:   ft.ID,
			Name: ft.Name,
			Code: ft.Code,
		}
	}

	prodis := make([]admin.ProdiOption, len(dbProdis))
	for i, p := range dbProdis {
		prodis[i] = admin.ProdiOption{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		}
	}

	admin.SettingsFees(data, feeStructures, feeTypes, prodis, academicYear).Render(r.Context(), w)
}

// formatRupiah formats amount as Indonesian Rupiah
func formatRupiah(amount int64) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	if n <= 3 {
		return "Rp " + s
	}
	// Add thousand separators
	var result []byte
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return "Rp " + string(result)
}

// handleCreateFeeStructure handles POST /admin/settings/fees
func (h *AdminHandler) handleCreateFeeStructure(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	feeTypeID := r.FormValue("fee_type_id")
	prodiIDStr := r.FormValue("prodi_id")
	academicYear := r.FormValue("academic_year")
	amountStr := r.FormValue("amount")

	if feeTypeID == "" || academicYear == "" || amountStr == "" {
		http.Error(w, "Fee type, academic year, and amount are required", http.StatusBadRequest)
		return
	}

	var amount int64
	if _, err := fmt.Sscanf(amountStr, "%d", &amount); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var prodiID *string
	if prodiIDStr != "" {
		prodiID = &prodiIDStr
	}

	fee, err := model.CreateFeeStructure(r.Context(), feeTypeID, prodiID, academicYear, amount)
	if err != nil {
		slog.Error("failed to create fee structure", "error", err)
		http.Error(w, "Failed to create fee structure", http.StatusInternalServerError)
		return
	}

	slog.Info("fee structure created", "fee_id", fee.ID)

	// Return the new fee row
	h.renderFeeRow(w, r, fee.ID)
}

// handleUpdateFeeStructure handles POST /admin/settings/fees/{id}
func (h *AdminHandler) handleUpdateFeeStructure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	if amountStr == "" {
		http.Error(w, "Amount is required", http.StatusBadRequest)
		return
	}

	var amount int64
	if _, err := fmt.Sscanf(amountStr, "%d", &amount); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	if err := model.UpdateFeeStructure(r.Context(), id, amount); err != nil {
		slog.Error("failed to update fee structure", "error", err, "fee_id", id)
		http.Error(w, "Failed to update fee structure", http.StatusInternalServerError)
		return
	}

	slog.Info("fee structure updated", "fee_id", id)

	// Return updated fee row
	h.renderFeeRow(w, r, id)
}

// handleToggleFeeStructureActive handles POST /admin/settings/fees/{id}/toggle-active
func (h *AdminHandler) handleToggleFeeStructureActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleFeeStructureActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle fee structure active", "error", err, "fee_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("fee structure active status toggled", "fee_id", id)

	// Return updated fee row
	h.renderFeeRow(w, r, id)
}

// renderFeeRow renders a single fee row for HTMX updates
func (h *AdminHandler) renderFeeRow(w http.ResponseWriter, r *http.Request, feeID string) {
	fee, err := model.FindFeeStructureByID(r.Context(), feeID)
	if err != nil {
		slog.Error("failed to find fee structure", "error", err)
		http.Error(w, "Failed to load fee structure", http.StatusInternalServerError)
		return
	}

	if fee == nil {
		http.NotFound(w, r)
		return
	}

	// Fetch fee type
	feeType, err := model.FindFeeTypeByID(r.Context(), fee.IDFeeType)
	if err != nil {
		slog.Error("failed to find fee type", "error", err)
		http.Error(w, "Failed to load fee type", http.StatusInternalServerError)
		return
	}

	// Fetch prodi if exists
	var prodiName, prodiCode string
	var prodiID string
	if fee.IDProdi != nil {
		prodiID = *fee.IDProdi
		prodi, err := model.FindProdiByID(r.Context(), *fee.IDProdi)
		if err != nil {
			slog.Error("failed to find prodi", "error", err)
		} else if prodi != nil {
			prodiName = prodi.Name
			prodiCode = prodi.Code
		}
	}

	// Fetch fee types and prodis for edit modal
	dbFeeTypes, err := model.ListFeeTypes(r.Context())
	if err != nil {
		slog.Error("failed to list fee types", "error", err)
	}

	dbProdis, err := model.ListProdis(r.Context(), true)
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
	}

	feeTypes := make([]admin.FeeTypeOption, len(dbFeeTypes))
	for i, ft := range dbFeeTypes {
		feeTypes[i] = admin.FeeTypeOption{
			ID:   ft.ID,
			Name: ft.Name,
			Code: ft.Code,
		}
	}

	prodis := make([]admin.ProdiOption, len(dbProdis))
	for i, p := range dbProdis {
		prodis[i] = admin.ProdiOption{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		}
	}

	item := admin.FeeStructureItem{
		ID:           fee.ID,
		IDFeeType:    fee.IDFeeType,
		FeeTypeName:  feeType.Name,
		FeeTypeCode:  feeType.Code,
		IDProdi:      prodiID,
		ProdiName:    prodiName,
		ProdiCode:    prodiCode,
		AcademicYear: fee.AcademicYear,
		Amount:       fee.Amount,
		AmountStr:    formatRupiah(fee.Amount),
		IsActive:     fee.IsActive,
	}

	admin.FeeRow(item, feeTypes, prodis).Render(r.Context(), w)
}
