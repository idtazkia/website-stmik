package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/auth"
	"github.com/idtazkia/stmik-admission-api/mockdata"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// AdminHandler handles all admin routes
type AdminHandler struct {
	sessionMgr *auth.SessionManager
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(sessionMgr *auth.SessionManager) *AdminHandler {
	return &AdminHandler{sessionMgr: sessionMgr}
}

// RegisterRoutes registers all admin routes to the mux
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	// Helper to wrap handlers with auth middleware
	protected := func(handler http.HandlerFunc) http.Handler {
		return RequireAuth(h.sessionMgr, handler)
	}

	// Dashboard
	mux.Handle("GET /admin", protected(h.handleDashboard))
	mux.Handle("GET /admin/", protected(h.handleDashboard))

	// Consultant personal dashboard
	mux.Handle("GET /admin/my-dashboard", protected(h.handleConsultantDashboard))

	// Candidates
	mux.Handle("GET /admin/candidates", protected(h.handleCandidates))
	mux.Handle("GET /admin/candidates/{id}", protected(h.handleCandidateDetail))
	mux.Handle("GET /admin/candidates/{id}/interaction", protected(h.handleInteractionForm))

	// Documents
	mux.Handle("GET /admin/documents", protected(h.handleDocumentReview))

	// Marketing
	mux.Handle("GET /admin/campaigns", protected(h.handleCampaigns))
	mux.Handle("GET /admin/referrers", protected(h.handleReferrers))
	mux.Handle("GET /admin/referral-claims", protected(h.handleReferralClaims))
	mux.Handle("GET /admin/commissions", protected(h.handleCommissions))

	// Reports
	mux.Handle("GET /admin/reports/funnel", protected(h.handleFunnelReport))
	mux.Handle("GET /admin/reports/consultants", protected(h.handleConsultantsReport))
	mux.Handle("GET /admin/reports/campaigns", protected(h.handleCampaignsReport))

	// Settings
	mux.Handle("GET /admin/settings/users", protected(h.handleUsersSettings))
	mux.Handle("POST /admin/settings/users/{id}/role", protected(h.handleUpdateUserRole))
	mux.Handle("POST /admin/settings/users/{id}/supervisor", protected(h.handleUpdateUserSupervisor))
	mux.Handle("POST /admin/settings/users/{id}/toggle-active", protected(h.handleToggleUserActive))
	mux.Handle("GET /admin/settings/programs", protected(h.handleProgramsSettings))
	mux.Handle("POST /admin/settings/programs", protected(h.handleCreateProgram))
	mux.Handle("POST /admin/settings/programs/{id}", protected(h.handleUpdateProgram))
	mux.Handle("POST /admin/settings/programs/{id}/toggle-active", protected(h.handleToggleProgramActive))
	mux.Handle("GET /admin/settings/categories", protected(h.handleCategoriesSettings))
	mux.Handle("POST /admin/settings/categories", protected(h.handleCreateCategory))
	mux.Handle("POST /admin/settings/categories/{id}", protected(h.handleUpdateCategory))
	mux.Handle("POST /admin/settings/categories/{id}/toggle-active", protected(h.handleToggleCategoryActive))
	mux.Handle("POST /admin/settings/obstacles", protected(h.handleCreateObstacle))
	mux.Handle("POST /admin/settings/obstacles/{id}", protected(h.handleUpdateObstacle))
	mux.Handle("POST /admin/settings/obstacles/{id}/toggle-active", protected(h.handleToggleObstacleActive))
	mux.Handle("GET /admin/settings/fees", protected(h.handleFeesSettings))
	mux.Handle("GET /admin/settings/rewards", protected(h.handleRewardsSettings))
}

func (h *AdminHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Dashboard")
	stats := mockdata.GetAdminStats()
	dashboardStats := admin.DashboardStats{
		TotalCandidates:  stats.TotalCandidates,
		RegisteredCount:  stats.RegisteredCount,
		ProspectingCount: stats.ProspectingCount,
		CommittedCount:   stats.CommittedCount,
		EnrolledCount:    stats.EnrolledCount,
		LostCount:        stats.LostCount,
		OverdueFollowups: stats.OverdueFollowups,
		TodayFollowups:   stats.TodayFollowups,
		ThisMonthLeads:   stats.ThisMonthLeads,
	}
	admin.Dashboard(data, dashboardStats).Render(r.Context(), w)
}

func (h *AdminHandler) handleCandidates(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Kandidat")

	// Get filter parameters
	filter := admin.KandidatFilter{
		Status:     r.URL.Query().Get("status"),
		Prodi:      r.URL.Query().Get("prodi"),
		Consultant: r.URL.Query().Get("consultant"),
		Search:     r.URL.Query().Get("search"),
	}

	// Get filtered candidates from mockdata
	mockCandidates := mockdata.FilterCandidates(filter.Status, filter.Prodi, filter.Consultant, filter.Search)

	// Convert mockdata.CandidateView to admin.Candidate
	candidates := make([]admin.Candidate, len(mockCandidates))
	for i, c := range mockCandidates {
		candidates[i] = admin.Candidate{
			ID:             c.ID,
			Name:           c.Name,
			Email:          c.Email,
			Phone:          c.Phone,
			HighSchool:     c.HighSchool,
			ProdiName:      c.ProdiName,
			SourceType:     c.SourceType,
			CampaignName:   c.CampaignName,
			ReferrerName:   c.ReferrerName,
			Status:         c.Status,
			ConsultantName: c.ConsultantName,
			NextFollowup:   c.NextFollowup,
			IsOverdue:      c.IsOverdue,
		}
	}

	admin.KandidatList(data, candidates, filter).Render(r.Context(), w)
}

func (h *AdminHandler) handleCandidateDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	candidate := mockdata.GetCandidateByID(id)
	if candidate == nil {
		http.NotFound(w, r)
		return
	}

	data := NewPageDataWithUser(r.Context(),"Detail Kandidat - " + candidate.Name)

	// Convert to template type
	c := admin.CandidateDetail{
		ID:                  candidate.ID,
		Name:                candidate.Name,
		Email:               candidate.Email,
		Phone:               candidate.Phone,
		WhatsApp:            candidate.WhatsApp,
		Address:             candidate.Address,
		City:                candidate.City,
		Province:            candidate.Province,
		HighSchool:          candidate.HighSchool,
		GraduationYear:      candidate.GraduationYear,
		ProdiName:           candidate.ProdiName,
		SourceType:          candidate.SourceType,
		SourceDetail:        candidate.SourceDetail,
		CampaignName:        candidate.CampaignName,
		ReferrerName:        candidate.ReferrerName,
		Status:              candidate.Status,
		ConsultantName:      candidate.ConsultantName,
		RegistrationFeePaid: candidate.RegistrationFeePaid,
		CreatedAt:           candidate.CreatedAt,
	}

	// Get interactions
	mockInteractions := mockdata.GetInteractionsByCandidateID(id)
	interactions := make([]admin.Interaction, len(mockInteractions))
	for i, inter := range mockInteractions {
		interactions[i] = admin.Interaction{
			ID:                   inter.ID,
			Channel:              inter.Channel,
			Category:             inter.Category,
			CategorySentiment:    inter.CategorySentiment,
			Obstacle:             inter.Obstacle,
			Remarks:              inter.Remarks,
			NextFollowupDate:     inter.NextFollowupDate,
			SupervisorSuggestion: inter.SupervisorSuggestion,
			SuggestionRead:       inter.SuggestionRead,
			ConsultantName:       inter.ConsultantName,
			CreatedAt:            inter.CreatedAt,
		}
	}

	admin.KandidatDetail(data, c, interactions).Render(r.Context(), w)
}

// Placeholder handlers - will be implemented later
func (h *AdminHandler) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Kampanye")
	campaigns := []admin.CampaignItem{
		{ID: "1", Name: "Promo Early Bird", Type: "promo", Channel: "all", Period: "1 Jan - 28 Feb 2026", FeeOverride: "Gratis", Status: "active", Leads: "45", Enrolled: "12", Conversion: "26.7%"},
		{ID: "2", Name: "Education Expo Jakarta", Type: "event", Channel: "expo", Period: "15-17 Jan 2026", FeeOverride: "", Status: "active", Leads: "38", Enrolled: "10", Conversion: "26.3%"},
		{ID: "3", Name: "Instagram Ads Q1", Type: "ads", Channel: "instagram", Period: "1 Jan - 31 Mar 2026", FeeOverride: "", Status: "active", Leads: "52", Enrolled: "8", Conversion: "15.4%"},
		{ID: "4", Name: "Kunjungan Sekolah Q1", Type: "event", Channel: "school_visit", Period: "6 Jan - 31 Mar 2026", FeeOverride: "", Status: "active", Leads: "28", Enrolled: "12", Conversion: "42.9%"},
		{ID: "5", Name: "Google Ads Q4 2025", Type: "ads", Channel: "google", Period: "1 Oct - 31 Dec 2025", FeeOverride: "", Status: "ended", Leads: "35", Enrolled: "8", Conversion: "22.9%"},
	}
	admin.SettingsCampaigns(data, campaigns).Render(r.Context(), w)
}

func (h *AdminHandler) handleReferrers(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Referrer")
	referrers := []admin.ReferrerItem{
		{ID: "1", Name: "Pak Ahmad Fauzi", Type: "guru", Institution: "SMAN 1 Bogor", Phone: "081234567890", Code: "REF-AF01", Commission: "Rp 750.000", Referrals: "8", Enrolled: "5", TotalEarned: "Rp 3.750.000", Status: "active"},
		{ID: "2", Name: "Siti Nurhaliza", Type: "alumni", Institution: "STMIK Tazkia 2022", Phone: "081234567891", Code: "REF-SN02", Commission: "Rp 500.000", Referrals: "4", Enrolled: "2", TotalEarned: "Rp 1.000.000", Status: "active"},
		{ID: "3", Name: "PT Edutech Indonesia", Type: "partner", Institution: "Bimbel Edutech", Phone: "021-7654321", Code: "REF-EDU", Commission: "Rp 1.000.000", Referrals: "12", Enrolled: "6", TotalEarned: "Rp 6.000.000", Status: "active"},
		{ID: "4", Name: "Budi Santoso", Type: "staff", Institution: "STMIK Tazkia", Phone: "081234567893", Code: "REF-BS04", Commission: "Rp 250.000", Referrals: "3", Enrolled: "2", TotalEarned: "Rp 500.000", Status: "active"},
	}
	admin.SettingsReferrers(data, referrers).Render(r.Context(), w)
}

func (h *AdminHandler) handleReferralClaims(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Klaim Referral")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Klaim Referral - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleCommissions(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Komisi")
	commissions := []admin.CommissionItem{
		{ID: "1", ReferrerName: "Pak Ahmad Fauzi", ReferrerType: "guru", CandidateName: "Dimas Pratama", CandidateNIM: "2026SI003", Amount: "Rp 750.000", Status: "pending", EnrolledAt: "10 Jan 2026", BankName: "BSI", BankAccount: "7123456789"},
		{ID: "2", ReferrerName: "Siti Nurhaliza", ReferrerType: "alumni", CandidateName: "Rina Wulandari", CandidateNIM: "2026TI004", Amount: "Rp 500.000", Status: "approved", EnrolledAt: "8 Jan 2026", ApprovedAt: "12 Jan 2026", BankName: "BCA", BankAccount: "1234567890"},
		{ID: "3", ReferrerName: "PT Edutech Indonesia", ReferrerType: "partner", CandidateName: "Bayu Setiawan", CandidateNIM: "2026SI005", Amount: "Rp 1.000.000", Status: "paid", EnrolledAt: "5 Jan 2026", ApprovedAt: "7 Jan 2026", PaidAt: "10 Jan 2026", BankName: "Mandiri", BankAccount: "0987654321"},
	}
	stats := admin.CommissionStats{
		Pending: "3", PendingAmount: "Rp 2.250.000",
		Approved: "2", ApprovedAmount: "Rp 1.500.000",
		Paid: "5", PaidAmount: "Rp 4.000.000",
	}
	admin.Commissions(data, commissions, stats).Render(r.Context(), w)
}

func (h *AdminHandler) handleFunnelReport(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Laporan Funnel")
	stages := []admin.FunnelStage{
		{Name: "Registered", Count: "97", Percentage: "100", Color: "bg-gray-500"},
		{Name: "Prospecting", Count: "80", Percentage: "82", Color: "bg-blue-500"},
		{Name: "Committed", Count: "55", Percentage: "57", Color: "bg-yellow-500"},
		{Name: "Enrolled", Count: "40", Percentage: "41", Color: "bg-green-500"},
	}
	conversions := []admin.FunnelConversion{
		{From: "Registered", To: "Prospecting", Rate: "82.5%", Change: "+5%", IsPositive: true},
		{From: "Prospecting", To: "Committed", Rate: "68.8%", Change: "+2%", IsPositive: true},
		{From: "Committed", To: "Enrolled", Rate: "72.7%", Change: "-3%", IsPositive: false},
	}
	admin.ReportFunnel(data, stages, conversions).Render(r.Context(), w)
}

// handleConsultantsReport is implemented below with real template

func (h *AdminHandler) handleCampaignsReport(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"ROI Kampanye")
	campaigns := []admin.CampaignReportItem{
		{Name: "Promo Early Bird", Type: "promo", Channel: "all", Leads: "45", Prospecting: "38", Committed: "25", Enrolled: "12", Conversion: "26.7%", Cost: "Rp 0", CostPerLead: "Rp 0"},
		{Name: "Education Expo Jakarta", Type: "event", Channel: "expo", Leads: "38", Prospecting: "30", Committed: "18", Enrolled: "10", Conversion: "26.3%", Cost: "Rp 15.000.000", CostPerLead: "Rp 394.737"},
		{Name: "Instagram Ads Q1", Type: "ads", Channel: "instagram", Leads: "52", Prospecting: "35", Committed: "15", Enrolled: "8", Conversion: "15.4%", Cost: "Rp 8.000.000", CostPerLead: "Rp 153.846"},
		{Name: "Kunjungan Sekolah Q1", Type: "event", Channel: "school_visit", Leads: "28", Prospecting: "24", Committed: "18", Enrolled: "12", Conversion: "42.9%", Cost: "Rp 5.000.000", CostPerLead: "Rp 178.571"},
	}
	summary := admin.CampaignReportSummary{
		TotalLeads: "163", TotalEnrolled: "42", AvgConversion: "25.8%",
		TotalCost: "Rp 28.000.000", AvgCostPerLead: "Rp 171.779", BestCampaign: "Kunjungan Sekolah",
	}
	admin.ReportCampaigns(data, campaigns, summary).Render(r.Context(), w)
}

func (h *AdminHandler) handleUsersSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Users")

	// Fetch users from database
	dbUsers, err := model.ListUsers(r.Context(), "", false)
	if err != nil {
		slog.Error("failed to list users", "error", err)
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	// Fetch supervisors for dropdown
	dbSupervisors, err := model.ListSupervisors(r.Context())
	if err != nil {
		slog.Error("failed to list supervisors", "error", err)
		http.Error(w, "Failed to load supervisors", http.StatusInternalServerError)
		return
	}

	// Convert supervisors to template type
	supervisors := make([]admin.SupervisorOption, len(dbSupervisors))
	for i, s := range dbSupervisors {
		supervisors[i] = admin.SupervisorOption{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	// Convert to template type
	users := make([]admin.UserItem, len(dbUsers))
	for i, u := range dbUsers {
		status := "inactive"
		if u.IsActive {
			status = "active"
		}

		supervisor := "-"
		if u.SupervisorName != nil {
			supervisor = *u.SupervisorName
		}

		supervisorID := ""
		if u.SupervisorID != nil {
			supervisorID = *u.SupervisorID
		}

		lastLogin := "Belum pernah"
		if u.LastLoginAt != nil {
			lastLogin = formatRelativeTime(*u.LastLoginAt)
		}

		users[i] = admin.UserItem{
			ID:           u.ID,
			Name:         u.Name,
			Email:        u.Email,
			Role:         u.Role,
			Supervisor:   supervisor,
			SupervisorID: supervisorID,
			Status:       status,
			LastLogin:    lastLogin,
		}
	}

	// Fetch stats
	counts, err := model.CountUsersByRole(r.Context())
	if err != nil {
		slog.Error("failed to count users", "error", err)
	}
	stats := admin.UserStats{
		Total:      counts["admin"] + counts["supervisor"] + counts["consultant"],
		Admin:      counts["admin"],
		Supervisor: counts["supervisor"],
		Consultant: counts["consultant"],
	}

	admin.SettingsUsers(data, users, stats, supervisors).Render(r.Context(), w)
}

// handleUpdateUserRole handles POST /admin/settings/users/{id}/role
func (h *AdminHandler) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	role := r.FormValue("role")
	if role == "" {
		http.Error(w, "Role is required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateUserRole(r.Context(), id, role); err != nil {
		slog.Error("failed to update user role", "error", err, "user_id", id, "role", role)
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	slog.Info("user role updated", "user_id", id, "role", role)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// handleUpdateUserSupervisor handles POST /admin/settings/users/{id}/supervisor
func (h *AdminHandler) handleUpdateUserSupervisor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	supervisorID := r.FormValue("supervisor_id")
	var supIDPtr *string
	if supervisorID != "" {
		supIDPtr = &supervisorID
	}

	if err := model.UpdateUserSupervisor(r.Context(), id, supIDPtr); err != nil {
		slog.Error("failed to update user supervisor", "error", err, "user_id", id, "supervisor_id", supervisorID)
		http.Error(w, "Failed to update supervisor", http.StatusInternalServerError)
		return
	}

	slog.Info("user supervisor updated", "user_id", id, "supervisor_id", supervisorID)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// handleToggleUserActive handles POST /admin/settings/users/{id}/toggle-active
func (h *AdminHandler) handleToggleUserActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleUserActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle user active", "error", err, "user_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("user active status toggled", "user_id", id)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// renderUserRow renders a single user row for HTMX updates
func (h *AdminHandler) renderUserRow(w http.ResponseWriter, r *http.Request, userID string) {
	// Fetch user by ID with supervisor info
	dbUsers, err := model.ListUsers(r.Context(), "", false)
	if err != nil {
		slog.Error("failed to list users", "error", err)
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Fetch supervisors for dropdown
	dbSupervisors, err := model.ListSupervisors(r.Context())
	if err != nil {
		slog.Error("failed to list supervisors", "error", err)
		http.Error(w, "Failed to load supervisors", http.StatusInternalServerError)
		return
	}

	// Convert supervisors to template type
	supervisors := make([]admin.SupervisorOption, len(dbSupervisors))
	for i, s := range dbSupervisors {
		supervisors[i] = admin.SupervisorOption{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	// Find the user
	var userItem admin.UserItem
	found := false
	for _, u := range dbUsers {
		if u.ID == userID {
			status := "inactive"
			if u.IsActive {
				status = "active"
			}

			supervisor := "-"
			if u.SupervisorName != nil {
				supervisor = *u.SupervisorName
			}

			supervisorID := ""
			if u.SupervisorID != nil {
				supervisorID = *u.SupervisorID
			}

			lastLogin := "Belum pernah"
			if u.LastLoginAt != nil {
				lastLogin = formatRelativeTime(*u.LastLoginAt)
			}

			userItem = admin.UserItem{
				ID:           u.ID,
				Name:         u.Name,
				Email:        u.Email,
				Role:         u.Role,
				Supervisor:   supervisor,
				SupervisorID: supervisorID,
				Status:       status,
				LastLogin:    lastLogin,
			}
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	admin.UserRow(userItem, supervisors).Render(r.Context(), w)
}

// formatRelativeTime formats a time as relative to now
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Baru saja"
	}
	if diff < time.Hour {
		return "Beberapa menit lalu"
	}
	if diff < 24*time.Hour {
		return "Hari ini"
	}
	if diff < 48*time.Hour {
		return "Kemarin"
	}
	days := int(diff.Hours() / 24)
	if days < 7 {
		return formatDays(days) + " lalu"
	}
	return t.Format("2 Jan 2006")
}

func formatDays(days int) string {
	return fmt.Sprintf("%d hari", days)
}

func (h *AdminHandler) handleProgramsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Prodi")

	// Fetch programs from database
	dbProdis, err := model.ListProdis(r.Context(), false)
	if err != nil {
		slog.Error("failed to list prodis", "error", err)
		http.Error(w, "Failed to load programs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	programs := make([]admin.ProgramItem, len(dbProdis))
	for i, p := range dbProdis {
		status := "inactive"
		if p.IsActive {
			status = "active"
		}

		programs[i] = admin.ProgramItem{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Level:    p.Degree,
			SPPFee:   "-",       // TODO: fetch from fee_structures
			Status:   status,
			Students: "-",       // TODO: count enrolled candidates
		}
	}

	admin.SettingsPrograms(data, programs).Render(r.Context(), w)
}

// handleCreateProgram handles POST /admin/settings/programs
func (h *AdminHandler) handleCreateProgram(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	degree := r.FormValue("degree")

	if name == "" || code == "" || degree == "" {
		http.Error(w, "Name, code, and degree are required", http.StatusBadRequest)
		return
	}

	prodi, err := model.CreateProdi(r.Context(), name, code, degree)
	if err != nil {
		slog.Error("failed to create prodi", "error", err)
		http.Error(w, "Failed to create program", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi created", "prodi_id", prodi.ID, "code", prodi.Code)

	// Return the new program card
	h.renderProgramCard(w, r, prodi.ID)
}

// handleUpdateProgram handles POST /admin/settings/programs/{id}
func (h *AdminHandler) handleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	degree := r.FormValue("degree")

	if name == "" || code == "" || degree == "" {
		http.Error(w, "Name, code, and degree are required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateProdi(r.Context(), id, name, code, degree); err != nil {
		slog.Error("failed to update prodi", "error", err, "prodi_id", id)
		http.Error(w, "Failed to update program", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi updated", "prodi_id", id, "code", code)

	// Return updated program card
	h.renderProgramCard(w, r, id)
}

// handleToggleProgramActive handles POST /admin/settings/programs/{id}/toggle-active
func (h *AdminHandler) handleToggleProgramActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleProdiActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle prodi active", "error", err, "prodi_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("prodi active status toggled", "prodi_id", id)

	// Return updated program card
	h.renderProgramCard(w, r, id)
}

// renderProgramCard renders a single program card for HTMX updates
func (h *AdminHandler) renderProgramCard(w http.ResponseWriter, r *http.Request, prodiID string) {
	prodi, err := model.FindProdiByID(r.Context(), prodiID)
	if err != nil {
		slog.Error("failed to find prodi", "error", err)
		http.Error(w, "Failed to load program", http.StatusInternalServerError)
		return
	}

	if prodi == nil {
		http.NotFound(w, r)
		return
	}

	status := "inactive"
	if prodi.IsActive {
		status = "active"
	}

	programItem := admin.ProgramItem{
		ID:       prodi.ID,
		Name:     prodi.Name,
		Code:     prodi.Code,
		Level:    prodi.Degree,
		SPPFee:   "-",   // TODO: fetch from fee_structures
		Status:   status,
		Students: "-",   // TODO: count enrolled candidates
	}

	admin.ProgramCard(programItem).Render(r.Context(), w)
}

func (h *AdminHandler) handleFeesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Biaya")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Biaya - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleRewardsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Reward")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Reward - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleCategoriesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kategori")

	// Fetch categories from database
	dbCategories, err := model.ListInteractionCategories(r.Context(), false)
	if err != nil {
		slog.Error("failed to list interaction categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	categories := make([]admin.CategoryItem, len(dbCategories))
	for i, c := range dbCategories {
		categories[i] = admin.CategoryItem{
			ID:        c.ID,
			Name:      c.Name,
			Sentiment: c.Sentiment,
			Count:     "0", // TODO: count interactions using this category
			IsActive:  c.IsActive,
		}
	}

	// Fetch obstacles from database
	dbObstacles, err := model.ListObstacles(r.Context(), false)
	if err != nil {
		slog.Error("failed to list obstacles", "error", err)
		http.Error(w, "Failed to load obstacles", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	obstacles := make([]admin.ObstacleItem, len(dbObstacles))
	for i, o := range dbObstacles {
		obstacles[i] = admin.ObstacleItem{
			ID:       o.ID,
			Name:     o.Name,
			Count:    "0", // TODO: count interactions with this obstacle
			IsActive: o.IsActive,
		}
	}

	admin.SettingsCategories(data, categories, obstacles).Render(r.Context(), w)
}

// handleCreateCategory handles POST /admin/settings/categories
func (h *AdminHandler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	sentiment := r.FormValue("sentiment")

	if name == "" || sentiment == "" {
		http.Error(w, "Name and sentiment are required", http.StatusBadRequest)
		return
	}

	cat, err := model.CreateInteractionCategory(r.Context(), name, sentiment, 0)
	if err != nil {
		slog.Error("failed to create category", "error", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category created", "category_id", cat.ID)

	// Return new category card
	h.renderCategoryCard(w, r, cat.ID)
}

// handleUpdateCategory handles POST /admin/settings/categories/{id}
func (h *AdminHandler) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	sentiment := r.FormValue("sentiment")

	if name == "" || sentiment == "" {
		http.Error(w, "Name and sentiment are required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateInteractionCategory(r.Context(), id, name, sentiment, 0); err != nil {
		slog.Error("failed to update category", "error", err, "category_id", id)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category updated", "category_id", id)

	// Return updated category card
	h.renderCategoryCard(w, r, id)
}

// handleToggleCategoryActive handles POST /admin/settings/categories/{id}/toggle-active
func (h *AdminHandler) handleToggleCategoryActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleInteractionCategoryActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle category active", "error", err, "category_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category active status toggled", "category_id", id)

	// Return updated category card
	h.renderCategoryCard(w, r, id)
}

// renderCategoryCard renders a single category card for HTMX updates
func (h *AdminHandler) renderCategoryCard(w http.ResponseWriter, r *http.Request, categoryID string) {
	cat, err := model.FindInteractionCategoryByID(r.Context(), categoryID)
	if err != nil {
		slog.Error("failed to find category", "error", err)
		http.Error(w, "Failed to load category", http.StatusInternalServerError)
		return
	}

	item := admin.CategoryItem{
		ID:        cat.ID,
		Name:      cat.Name,
		Sentiment: cat.Sentiment,
		Count:     "0", // TODO: count interactions using this category
		IsActive:  cat.IsActive,
	}

	admin.CategoryCard(item).Render(r.Context(), w)
}

// handleCreateObstacle handles POST /admin/settings/obstacles
func (h *AdminHandler) handleCreateObstacle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	obs, err := model.CreateObstacle(r.Context(), name, nil, 0)
	if err != nil {
		slog.Error("failed to create obstacle", "error", err)
		http.Error(w, "Failed to create obstacle", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle created", "obstacle_id", obs.ID)

	// Return new obstacle card
	h.renderObstacleCard(w, r, obs.ID)
}

// handleUpdateObstacle handles POST /admin/settings/obstacles/{id}
func (h *AdminHandler) handleUpdateObstacle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateObstacle(r.Context(), id, name, nil, 0); err != nil {
		slog.Error("failed to update obstacle", "error", err, "obstacle_id", id)
		http.Error(w, "Failed to update obstacle", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle updated", "obstacle_id", id)

	// Return updated obstacle card
	h.renderObstacleCard(w, r, id)
}

// handleToggleObstacleActive handles POST /admin/settings/obstacles/{id}/toggle-active
func (h *AdminHandler) handleToggleObstacleActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleObstacleActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle obstacle active", "error", err, "obstacle_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle active status toggled", "obstacle_id", id)

	// Return updated obstacle card
	h.renderObstacleCard(w, r, id)
}

// renderObstacleCard renders a single obstacle card for HTMX updates
func (h *AdminHandler) renderObstacleCard(w http.ResponseWriter, r *http.Request, obstacleID string) {
	obs, err := model.FindObstacleByID(r.Context(), obstacleID)
	if err != nil {
		slog.Error("failed to find obstacle", "error", err)
		http.Error(w, "Failed to load obstacle", http.StatusInternalServerError)
		return
	}

	item := admin.ObstacleItem{
		ID:       obs.ID,
		Name:     obs.Name,
		Count:    "0", // TODO: count interactions with this obstacle
		IsActive: obs.IsActive,
	}

	admin.ObstacleCard(item).Render(r.Context(), w)
}


func (h *AdminHandler) handleConsultantDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Dashboard Saya")

	stats := admin.ConsultantDashboardStats{
		ConsultantName:      "Siti Rahayu",
		TodayDate:           "15 Januari 2026",
		MyCandidates:        "5",
		MyProspecting:       "2",
		MyCommitted:         "1",
		MyEnrolled:          "0",
		OverdueCount:        "1",
		TodayTasks:          "2",
		UnreadSuggestions:   "1",
		MonthlyNewLeads:     "8",
		MonthlyInteractions: "25",
		MonthlyCommits:      "3",
		MonthlyEnrollments:  "1",
	}

	overdueList := []admin.CandidateSummary{
		{ID: "8", Name: "Citra Dewi", ProdiName: "Teknik Informatika", WhatsApp: "081234567897", Status: "prospecting", LastContact: "3 hari lalu"},
	}

	todayTasks := []admin.CandidateSummary{
		{ID: "2", Name: "Putri Amelia", ProdiName: "Teknik Informatika", WhatsApp: "081234567891", Status: "prospecting", LastContact: "hari ini"},
		{ID: "7", Name: "Bayu Setiawan", ProdiName: "Sistem Informasi", WhatsApp: "081234567896", Status: "registered", LastContact: "kemarin"},
	}

	suggestions := []admin.SupervisorSuggestion{
		{ID: "1", CandidateID: "2", CandidateName: "Putri Amelia", Suggestion: "Tawarkan promo early bird untuk diskon biaya pendaftaran", SupervisorName: "Budi Santoso", CreatedAt: "14 Jan 2026", IsRead: false},
	}

	admin.ConsultantDashboard(data, stats, overdueList, todayTasks, suggestions).Render(r.Context(), w)
}

func (h *AdminHandler) handleConsultantsReport(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Performa Konsultan")

	filter := admin.ReportFilter{
		Period:    r.URL.Query().Get("period"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}
	if filter.Period == "" {
		filter.Period = "this_month"
	}

	consultants := []admin.ConsultantPerformance{
		{Rank: "1", Name: "Siti Rahayu", Email: "konsultan1@tazkia.ac.id", SupervisorName: "Budi Santoso", Leads: "15", Interactions: "45", Commits: "5", Enrollments: "3", ConversionRate: 20, ConversionRateStr: "20", AvgDaysToCommit: "12", InteractionsPerCandidate: "3.0"},
		{Rank: "2", Name: "Ahmad Hidayat", Email: "konsultan2@tazkia.ac.id", SupervisorName: "Budi Santoso", Leads: "12", Interactions: "30", Commits: "3", Enrollments: "2", ConversionRate: 16.7, ConversionRateStr: "16.7", AvgDaysToCommit: "15", InteractionsPerCandidate: "2.5"},
		{Rank: "3", Name: "Dewi Lestari", Email: "konsultan3@tazkia.ac.id", SupervisorName: "Budi Santoso", Leads: "8", Interactions: "20", Commits: "1", Enrollments: "0", ConversionRate: 0, ConversionRateStr: "0", AvgDaysToCommit: "-", InteractionsPerCandidate: "2.5"},
	}

	summary := admin.ReportSummary{
		TotalLeads:        "35",
		TotalInteractions: "95",
		TotalCommits:      "9",
		TotalEnrollments:  "5",
	}

	admin.ConsultantReport(data, filter, consultants, summary).Render(r.Context(), w)
}

func (h *AdminHandler) handleInteractionForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	candidate := mockdata.GetCandidateByID(id)
	if candidate == nil {
		http.NotFound(w, r)
		return
	}

	data := NewPageDataWithUser(r.Context(),"Log Interaksi - " + candidate.Name)

	candidateSummary := admin.CandidateSummary{
		ID:        candidate.ID,
		Name:      candidate.Name,
		ProdiName: candidate.ProdiName,
		WhatsApp:  candidate.WhatsApp,
		Status:    candidate.Status,
	}

	categories := []admin.InteractionCategoryOption{
		{Value: "interested", Label: "Tertarik", Icon: "😊", Sentiment: "positive"},
		{Value: "considering", Label: "Mempertimbangkan", Icon: "🤔", Sentiment: "neutral"},
		{Value: "hesitant", Label: "Ragu-ragu", Icon: "😐", Sentiment: "neutral"},
		{Value: "cold", Label: "Dingin", Icon: "😞", Sentiment: "negative"},
		{Value: "unreachable", Label: "Tidak Terhubung", Icon: "📵", Sentiment: "negative"},
	}

	obstacles := []admin.ObstacleOption{
		{Value: "price", Label: "Biaya terlalu mahal"},
		{Value: "location", Label: "Lokasi jauh"},
		{Value: "parents", Label: "Orang tua belum setuju"},
		{Value: "timing", Label: "Waktu belum tepat"},
		{Value: "competitor", Label: "Memilih kampus lain"},
	}

	admin.InteractionForm(data, candidateSummary, categories, obstacles).Render(r.Context(), w)
}

func (h *AdminHandler) handleDocumentReview(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(),"Review Dokumen")

	filter := admin.DocumentFilter{
		Status: r.URL.Query().Get("status"),
		Type:   r.URL.Query().Get("type"),
		Search: r.URL.Query().Get("search"),
	}

	stats := admin.DocumentStats{
		Pending:       "5",
		ApprovedToday: "12",
		RejectedToday: "2",
		Total:         "45",
	}

	documents := []admin.DocumentReviewItem{
		{ID: "1", CandidateID: "2", CandidateName: "Putri Amelia", ProdiName: "Teknik Informatika", Type: "ktp", TypeName: "KTP", FileName: "ktp_putri.jpg", FileSize: "1.2 MB", FileURL: "/uploads/ktp_putri.jpg", ThumbnailURL: "/uploads/ktp_putri_thumb.jpg", IsImage: true, Status: "pending", UploadedAt: "15 Jan 2026 10:30"},
		{ID: "2", CandidateID: "2", CandidateName: "Putri Amelia", ProdiName: "Teknik Informatika", Type: "photo", TypeName: "Foto", FileName: "foto_putri.jpg", FileSize: "850 KB", FileURL: "/uploads/foto_putri.jpg", ThumbnailURL: "/uploads/foto_putri_thumb.jpg", IsImage: true, Status: "pending", UploadedAt: "15 Jan 2026 10:32"},
		{ID: "3", CandidateID: "3", CandidateName: "Dimas Pratama", ProdiName: "Sistem Informasi", Type: "ijazah", TypeName: "Ijazah", FileName: "ijazah_dimas.pdf", FileSize: "2.5 MB", FileURL: "/uploads/ijazah_dimas.pdf", ThumbnailURL: "", IsImage: false, Status: "pending", UploadedAt: "14 Jan 2026 15:20"},
		{ID: "4", CandidateID: "1", CandidateName: "Muhammad Rizky", ProdiName: "Sistem Informasi", Type: "ktp", TypeName: "KTP", FileName: "ktp_rizky.jpg", FileSize: "900 KB", FileURL: "/uploads/ktp_rizky.jpg", ThumbnailURL: "/uploads/ktp_rizky_thumb.jpg", IsImage: true, Status: "rejected", RejectionReason: "Gambar buram/tidak jelas", UploadedAt: "13 Jan 2026 09:15"},
	}

	admin.DocumentReviewList(data, filter, documents, stats).Render(r.Context(), w)
}
