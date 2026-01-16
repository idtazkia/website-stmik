package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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
	mux.Handle("POST /admin/settings/fees", protected(h.handleCreateFeeStructure))
	mux.Handle("POST /admin/settings/fees/{id}", protected(h.handleUpdateFeeStructure))
	mux.Handle("POST /admin/settings/fees/{id}/toggle-active", protected(h.handleToggleFeeStructureActive))
	mux.Handle("GET /admin/settings/campaigns", protected(h.handleCampaignsSettings))
	mux.Handle("POST /admin/settings/campaigns", protected(h.handleCreateCampaign))
	mux.Handle("POST /admin/settings/campaigns/{id}", protected(h.handleUpdateCampaign))
	mux.Handle("POST /admin/settings/campaigns/{id}/toggle-active", protected(h.handleToggleCampaignActive))
	mux.Handle("GET /admin/settings/rewards", protected(h.handleRewardsSettings))
	mux.Handle("POST /admin/settings/rewards", protected(h.handleCreateRewardConfig))
	mux.Handle("POST /admin/settings/rewards/{id}", protected(h.handleUpdateRewardConfig))
	mux.Handle("POST /admin/settings/rewards/{id}/toggle-active", protected(h.handleToggleRewardConfigActive))
	mux.Handle("POST /admin/settings/mgm-rewards", protected(h.handleCreateMGMRewardConfig))
	mux.Handle("POST /admin/settings/mgm-rewards/{id}", protected(h.handleUpdateMGMRewardConfig))
	mux.Handle("POST /admin/settings/mgm-rewards/{id}/toggle-active", protected(h.handleToggleMGMRewardConfigActive))
	mux.Handle("GET /admin/settings/referrers", protected(h.handleReferrersSettings))
	mux.Handle("POST /admin/settings/referrers", protected(h.handleCreateReferrer))
	mux.Handle("POST /admin/settings/referrers/{id}", protected(h.handleUpdateReferrer))
	mux.Handle("POST /admin/settings/referrers/{id}/toggle-active", protected(h.handleToggleReferrerActive))
	mux.Handle("GET /admin/settings/assignment", protected(h.handleAssignmentSettings))
	mux.Handle("POST /admin/settings/assignment/{id}/activate", protected(h.handleActivateAlgorithm))
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
	// Temporary: redirect to settings campaigns page
	// TODO: Implement proper campaigns dashboard page
	campaigns := []admin.CampaignItem{
		{ID: "1", Name: "Promo Early Bird", Type: "promo", Channel: "instagram", StartDate: "2026-01-01", EndDate: "2026-02-28", FeeOverrideStr: "Gratis", IsActive: true},
		{ID: "2", Name: "Education Expo Jakarta", Type: "event", Channel: "expo", StartDate: "2026-01-15", EndDate: "2026-01-17", IsActive: true},
		{ID: "3", Name: "Instagram Ads Q1", Type: "ads", Channel: "instagram", StartDate: "2026-01-01", EndDate: "2026-03-31", IsActive: true},
		{ID: "4", Name: "Kunjungan Sekolah Q1", Type: "event", Channel: "school_visit", StartDate: "2026-01-06", EndDate: "2026-03-31", IsActive: true},
		{ID: "5", Name: "Google Ads Q4 2025", Type: "ads", Channel: "google", StartDate: "2025-10-01", EndDate: "2025-12-31", IsActive: false},
	}
	admin.SettingsCampaigns(data, campaigns).Render(r.Context(), w)
}

func (h *AdminHandler) handleReferrers(w http.ResponseWriter, r *http.Request) {
	// Redirect to settings page for referrers management
	h.handleReferrersSettings(w, r)
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
		if f.ProdiID != nil {
			prodiID = *f.ProdiID
		}
		if f.ProdiName != nil {
			prodiName = *f.ProdiName
		}
		if f.ProdiCode != nil {
			prodiCode = *f.ProdiCode
		}
		feeStructures[i] = admin.FeeStructureItem{
			ID:           f.ID,
			FeeTypeID:    f.FeeTypeID,
			FeeTypeName:  f.FeeTypeName,
			FeeTypeCode:  f.FeeTypeCode,
			ProdiID:      prodiID,
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
	feeType, err := model.FindFeeTypeByID(r.Context(), fee.FeeTypeID)
	if err != nil {
		slog.Error("failed to find fee type", "error", err)
		http.Error(w, "Failed to load fee type", http.StatusInternalServerError)
		return
	}

	// Fetch prodi if exists
	var prodiName, prodiCode string
	var prodiID string
	if fee.ProdiID != nil {
		prodiID = *fee.ProdiID
		prodi, err := model.FindProdiByID(r.Context(), *fee.ProdiID)
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
		FeeTypeID:    fee.FeeTypeID,
		FeeTypeName:  feeType.Name,
		FeeTypeCode:  feeType.Code,
		ProdiID:      prodiID,
		ProdiName:    prodiName,
		ProdiCode:    prodiCode,
		AcademicYear: fee.AcademicYear,
		Amount:       fee.Amount,
		AmountStr:    formatRupiah(fee.Amount),
		IsActive:     fee.IsActive,
	}

	admin.FeeRow(item, feeTypes, prodis).Render(r.Context(), w)
}

func (h *AdminHandler) handleRewardsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Konfigurasi Reward")

	// Fetch reward configs from database
	dbRewards, err := model.ListRewardConfigs(r.Context())
	if err != nil {
		slog.Error("failed to list reward configs", "error", err)
		http.Error(w, "Failed to load reward configs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	rewards := make([]admin.RewardConfigItem, len(dbRewards))
	for i, r := range dbRewards {
		amountStr := formatRupiah(r.Amount)
		if r.IsPercentage {
			amountStr = fmt.Sprintf("%d%%", r.Amount)
		}
		description := ""
		if r.Description != nil {
			description = *r.Description
		}
		rewards[i] = admin.RewardConfigItem{
			ID:           r.ID,
			ReferrerType: r.ReferrerType,
			RewardType:   r.RewardType,
			Amount:       r.Amount,
			AmountStr:    amountStr,
			IsPercentage: r.IsPercentage,
			TriggerEvent: r.TriggerEvent,
			Description:  description,
			IsActive:     r.IsActive,
		}
	}

	// Fetch MGM reward configs from database
	dbMGMRewards, err := model.ListMGMRewardConfigs(r.Context())
	if err != nil {
		slog.Error("failed to list MGM reward configs", "error", err)
		http.Error(w, "Failed to load MGM reward configs", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	mgmRewards := make([]admin.MGMRewardConfigItem, len(dbMGMRewards))
	for i, m := range dbMGMRewards {
		referrerStr := formatRupiah(m.ReferrerAmount)
		refereeStr := ""
		if m.RefereeAmount != nil {
			refereeStr = formatRupiah(*m.RefereeAmount)
		}
		description := ""
		if m.Description != nil {
			description = *m.Description
		}
		mgmRewards[i] = admin.MGMRewardConfigItem{
			ID:             m.ID,
			AcademicYear:   m.AcademicYear,
			RewardType:     m.RewardType,
			ReferrerAmount: m.ReferrerAmount,
			ReferrerStr:    referrerStr,
			RefereeAmount:  m.RefereeAmount,
			RefereeStr:     refereeStr,
			TriggerEvent:   m.TriggerEvent,
			Description:    description,
			IsActive:       m.IsActive,
		}
	}

	admin.SettingsRewards(data, rewards, mgmRewards).Render(r.Context(), w)
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

// Campaign Settings Handlers

func (h *AdminHandler) handleCampaignsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kampanye")

	// Fetch campaigns from database
	dbCampaigns, err := model.ListCampaigns(r.Context(), false)
	if err != nil {
		slog.Error("failed to list campaigns", "error", err)
		http.Error(w, "Failed to load campaigns", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	campaigns := make([]admin.CampaignItem, len(dbCampaigns))
	for i, c := range dbCampaigns {
		var channel, description string
		if c.Channel != nil {
			channel = *c.Channel
		}
		if c.Description != nil {
			description = *c.Description
		}

		startDate := ""
		endDate := ""
		if c.StartDate != nil {
			startDate = c.StartDate.Format("2006-01-02")
		}
		if c.EndDate != nil {
			endDate = c.EndDate.Format("2006-01-02")
		}

		feeOverrideStr := ""
		if c.RegistrationFeeOverride != nil {
			feeOverrideStr = formatRupiah(*c.RegistrationFeeOverride)
		}

		campaigns[i] = admin.CampaignItem{
			ID:                      c.ID,
			Name:                    c.Name,
			Type:                    c.Type,
			Channel:                 channel,
			Description:             description,
			StartDate:               startDate,
			EndDate:                 endDate,
			RegistrationFeeOverride: c.RegistrationFeeOverride,
			FeeOverrideStr:          feeOverrideStr,
			IsActive:                c.IsActive,
		}
	}

	admin.SettingsCampaigns(data, campaigns).Render(r.Context(), w)
}

func (h *AdminHandler) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	campaignType := r.FormValue("type")
	channel := r.FormValue("channel")
	description := r.FormValue("description")
	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")
	feeOverrideStr := r.FormValue("registration_fee_override")

	if name == "" || campaignType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	var channelPtr, descPtr *string
	if channel != "" {
		channelPtr = &channel
	}
	if description != "" {
		descPtr = &description
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var feeOverride *int64
	if feeOverrideStr != "" {
		var fee int64
		if _, err := fmt.Sscanf(feeOverrideStr, "%d", &fee); err == nil {
			feeOverride = &fee
		}
	}

	campaign, err := model.CreateCampaign(r.Context(), name, campaignType, channelPtr, descPtr, startDate, endDate, feeOverride)
	if err != nil {
		slog.Error("failed to create campaign", "error", err)
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign created", "campaign_id", campaign.ID)

	// Return new campaign row
	h.renderCampaignRow(w, r, campaign.ID)
}

func (h *AdminHandler) handleUpdateCampaign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	campaignType := r.FormValue("type")
	channel := r.FormValue("channel")
	description := r.FormValue("description")
	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")
	feeOverrideStr := r.FormValue("registration_fee_override")

	if name == "" || campaignType == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	var channelPtr, descPtr *string
	if channel != "" {
		channelPtr = &channel
	}
	if description != "" {
		descPtr = &description
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var feeOverride *int64
	if feeOverrideStr != "" {
		var fee int64
		if _, err := fmt.Sscanf(feeOverrideStr, "%d", &fee); err == nil {
			feeOverride = &fee
		}
	}

	if err := model.UpdateCampaign(r.Context(), id, name, campaignType, channelPtr, descPtr, startDate, endDate, feeOverride); err != nil {
		slog.Error("failed to update campaign", "error", err, "campaign_id", id)
		http.Error(w, "Failed to update campaign", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign updated", "campaign_id", id)

	// Return updated campaign row
	h.renderCampaignRow(w, r, id)
}

func (h *AdminHandler) handleToggleCampaignActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleCampaignActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle campaign active", "error", err, "campaign_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("campaign active status toggled", "campaign_id", id)

	// Return updated campaign row
	h.renderCampaignRow(w, r, id)
}

func (h *AdminHandler) renderCampaignRow(w http.ResponseWriter, r *http.Request, campaignID string) {
	campaign, err := model.FindCampaignByID(r.Context(), campaignID)
	if err != nil {
		slog.Error("failed to find campaign", "error", err)
		http.Error(w, "Failed to load campaign", http.StatusInternalServerError)
		return
	}

	if campaign == nil {
		http.NotFound(w, r)
		return
	}

	var channel, description string
	if campaign.Channel != nil {
		channel = *campaign.Channel
	}
	if campaign.Description != nil {
		description = *campaign.Description
	}

	startDate := ""
	endDate := ""
	if campaign.StartDate != nil {
		startDate = campaign.StartDate.Format("2006-01-02")
	}
	if campaign.EndDate != nil {
		endDate = campaign.EndDate.Format("2006-01-02")
	}

	feeOverrideStr := ""
	if campaign.RegistrationFeeOverride != nil {
		feeOverrideStr = formatRupiah(*campaign.RegistrationFeeOverride)
	}

	item := admin.CampaignItem{
		ID:                      campaign.ID,
		Name:                    campaign.Name,
		Type:                    campaign.Type,
		Channel:                 channel,
		Description:             description,
		StartDate:               startDate,
		EndDate:                 endDate,
		RegistrationFeeOverride: campaign.RegistrationFeeOverride,
		FeeOverrideStr:          feeOverrideStr,
		IsActive:                campaign.IsActive,
	}

	admin.CampaignRow(item).Render(r.Context(), w)
}

// Reward Config Handlers

func (h *AdminHandler) handleCreateRewardConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	referrerType := r.FormValue("referrer_type")
	rewardType := r.FormValue("reward_type")
	amountStr := r.FormValue("amount")
	isPercentage := r.FormValue("is_percentage") == "on"
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	reward, err := model.CreateRewardConfig(r.Context(), referrerType, rewardType, amount, isPercentage, triggerEvent, description)
	if err != nil {
		slog.Error("failed to create reward config", "error", err)
		http.Error(w, "Failed to create reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config created", "reward_id", reward.ID)
	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) handleUpdateRewardConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing reward ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	referrerType := r.FormValue("referrer_type")
	rewardType := r.FormValue("reward_type")
	amountStr := r.FormValue("amount")
	isPercentage := r.FormValue("is_percentage") == "on"
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	err = model.UpdateRewardConfig(r.Context(), id, referrerType, rewardType, amount, isPercentage, triggerEvent, description)
	if err != nil {
		slog.Error("failed to update reward config", "error", err)
		http.Error(w, "Failed to update reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config updated", "reward_id", id)

	// Fetch updated reward and render
	reward, err := model.FindRewardConfigByID(r.Context(), id)
	if err != nil || reward == nil {
		http.Error(w, "Failed to fetch updated reward", http.StatusInternalServerError)
		return
	}

	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) handleToggleRewardConfigActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing reward ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleRewardConfigActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle reward config active", "error", err)
		http.Error(w, "Failed to toggle reward status", http.StatusInternalServerError)
		return
	}

	slog.Info("reward config active status toggled", "reward_id", id)

	// Fetch updated reward and render
	reward, err := model.FindRewardConfigByID(r.Context(), id)
	if err != nil || reward == nil {
		http.Error(w, "Failed to fetch updated reward", http.StatusInternalServerError)
		return
	}

	h.renderRewardCard(w, r, reward)
}

func (h *AdminHandler) renderRewardCard(w http.ResponseWriter, r *http.Request, reward *model.RewardConfig) {
	amountStr := formatRupiah(reward.Amount)
	if reward.IsPercentage {
		amountStr = fmt.Sprintf("%d%%", reward.Amount)
	}
	description := ""
	if reward.Description != nil {
		description = *reward.Description
	}

	item := admin.RewardConfigItem{
		ID:           reward.ID,
		ReferrerType: reward.ReferrerType,
		RewardType:   reward.RewardType,
		Amount:       reward.Amount,
		AmountStr:    amountStr,
		IsPercentage: reward.IsPercentage,
		TriggerEvent: reward.TriggerEvent,
		Description:  description,
		IsActive:     reward.IsActive,
	}

	admin.RewardCard(item).Render(r.Context(), w)
}

// MGM Reward Config Handlers

func (h *AdminHandler) handleCreateMGMRewardConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	academicYear := r.FormValue("academic_year")
	rewardType := r.FormValue("reward_type")
	referrerAmountStr := r.FormValue("referrer_amount")
	refereeAmountStr := r.FormValue("referee_amount")
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	referrerAmount, err := strconv.ParseInt(referrerAmountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid referrer amount", http.StatusBadRequest)
		return
	}

	var refereeAmount *int64
	if refereeAmountStr != "" {
		amt, err := strconv.ParseInt(refereeAmountStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid referee amount", http.StatusBadRequest)
			return
		}
		refereeAmount = &amt
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	mgmReward, err := model.CreateMGMRewardConfig(r.Context(), academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description)
	if err != nil {
		slog.Error("failed to create MGM reward config", "error", err)
		http.Error(w, "Failed to create MGM reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config created", "mgm_reward_id", mgmReward.ID)
	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) handleUpdateMGMRewardConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing MGM reward ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	academicYear := r.FormValue("academic_year")
	rewardType := r.FormValue("reward_type")
	referrerAmountStr := r.FormValue("referrer_amount")
	refereeAmountStr := r.FormValue("referee_amount")
	triggerEvent := r.FormValue("trigger_event")
	descriptionStr := r.FormValue("description")

	referrerAmount, err := strconv.ParseInt(referrerAmountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid referrer amount", http.StatusBadRequest)
		return
	}

	var refereeAmount *int64
	if refereeAmountStr != "" {
		amt, err := strconv.ParseInt(refereeAmountStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid referee amount", http.StatusBadRequest)
			return
		}
		refereeAmount = &amt
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	err = model.UpdateMGMRewardConfig(r.Context(), id, academicYear, rewardType, referrerAmount, refereeAmount, triggerEvent, description)
	if err != nil {
		slog.Error("failed to update MGM reward config", "error", err)
		http.Error(w, "Failed to update MGM reward config", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config updated", "mgm_reward_id", id)

	// Fetch updated MGM reward and render
	mgmReward, err := model.FindMGMRewardConfigByID(r.Context(), id)
	if err != nil || mgmReward == nil {
		http.Error(w, "Failed to fetch updated MGM reward", http.StatusInternalServerError)
		return
	}

	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) handleToggleMGMRewardConfigActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing MGM reward ID", http.StatusBadRequest)
		return
	}

	err := model.ToggleMGMRewardConfigActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to toggle MGM reward config active", "error", err)
		http.Error(w, "Failed to toggle MGM reward status", http.StatusInternalServerError)
		return
	}

	slog.Info("MGM reward config active status toggled", "mgm_reward_id", id)

	// Fetch updated MGM reward and render
	mgmReward, err := model.FindMGMRewardConfigByID(r.Context(), id)
	if err != nil || mgmReward == nil {
		http.Error(w, "Failed to fetch updated MGM reward", http.StatusInternalServerError)
		return
	}

	h.renderMGMRewardCard(w, r, mgmReward)
}

func (h *AdminHandler) renderMGMRewardCard(w http.ResponseWriter, r *http.Request, mgmReward *model.MGMRewardConfig) {
	referrerStr := formatRupiah(mgmReward.ReferrerAmount)
	refereeStr := ""
	if mgmReward.RefereeAmount != nil {
		refereeStr = formatRupiah(*mgmReward.RefereeAmount)
	}
	description := ""
	if mgmReward.Description != nil {
		description = *mgmReward.Description
	}

	item := admin.MGMRewardConfigItem{
		ID:             mgmReward.ID,
		AcademicYear:   mgmReward.AcademicYear,
		RewardType:     mgmReward.RewardType,
		ReferrerAmount: mgmReward.ReferrerAmount,
		ReferrerStr:    referrerStr,
		RefereeAmount:  mgmReward.RefereeAmount,
		RefereeStr:     refereeStr,
		TriggerEvent:   mgmReward.TriggerEvent,
		Description:    description,
		IsActive:       mgmReward.IsActive,
	}

	admin.MGMRewardCard(item).Render(r.Context(), w)
}

// Referrer Settings Handlers

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

// handleAssignmentSettings handles GET /admin/settings/assignment
func (h *AdminHandler) handleAssignmentSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Algoritma Penugasan")

	dbAlgorithms, err := model.ListAssignmentAlgorithms(r.Context())
	if err != nil {
		slog.Error("failed to list assignment algorithms", "error", err)
		http.Error(w, "Failed to load assignment algorithms", http.StatusInternalServerError)
		return
	}

	algorithms := make([]admin.AssignmentAlgorithmItem, len(dbAlgorithms))
	for i, alg := range dbAlgorithms {
		description := ""
		if alg.Description != nil {
			description = *alg.Description
		}
		algorithms[i] = admin.AssignmentAlgorithmItem{
			ID:          alg.ID,
			Name:        alg.Name,
			Code:        alg.Code,
			Description: description,
			IsActive:    alg.IsActive,
		}
	}

	admin.SettingsAssignment(data, algorithms).Render(r.Context(), w)
}

// handleActivateAlgorithm handles POST /admin/settings/assignment/{id}/activate
func (h *AdminHandler) handleActivateAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing algorithm ID", http.StatusBadRequest)
		return
	}

	err := model.SetAssignmentAlgorithmActive(r.Context(), id)
	if err != nil {
		slog.Error("failed to activate algorithm", "error", err)
		http.Error(w, "Failed to activate algorithm", http.StatusInternalServerError)
		return
	}

	slog.Info("assignment algorithm activated", "algorithm_id", id)

	// Return updated list
	dbAlgorithms, err := model.ListAssignmentAlgorithms(r.Context())
	if err != nil {
		slog.Error("failed to list assignment algorithms", "error", err)
		http.Error(w, "Failed to load assignment algorithms", http.StatusInternalServerError)
		return
	}

	algorithms := make([]admin.AssignmentAlgorithmItem, len(dbAlgorithms))
	for i, alg := range dbAlgorithms {
		description := ""
		if alg.Description != nil {
			description = *alg.Description
		}
		algorithms[i] = admin.AssignmentAlgorithmItem{
			ID:          alg.ID,
			Name:        alg.Name,
			Code:        alg.Code,
			Description: description,
			IsActive:    alg.IsActive,
		}
	}

	admin.AlgorithmList(algorithms).Render(r.Context(), w)
}
