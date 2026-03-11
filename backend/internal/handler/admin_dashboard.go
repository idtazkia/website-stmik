package handler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Dashboard")

	stats, err := model.GetCandidateStatusStats(ctx, nil, nil)
	if err != nil {
		slog.Error("Failed to get dashboard stats", "error", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	extraStats, err := model.GetAdminDashboardExtraStats(ctx)
	if err != nil {
		slog.Error("Failed to get extra dashboard stats", "error", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	// Calculate funnel bar percentages relative to total
	total := stats.Total
	if total == 0 {
		total = 1
	}
	prospPct := stats.Prospecting * 100 / total
	commitPct := stats.Committed * 100 / total
	enrollPct := stats.Enrolled * 100 / total

	dashboardStats := admin.DashboardStats{
		TotalCandidates:  fmt.Sprintf("%d", stats.Total),
		RegisteredCount:  fmt.Sprintf("%d", stats.Registered),
		ProspectingCount: fmt.Sprintf("%d", stats.Prospecting),
		CommittedCount:   fmt.Sprintf("%d", stats.Committed),
		EnrolledCount:    fmt.Sprintf("%d", stats.Enrolled),
		LostCount:        fmt.Sprintf("%d", stats.Lost),
		OverdueFollowups: fmt.Sprintf("%d", extraStats.OverdueFollowups),
		TodayFollowups:   fmt.Sprintf("%d", extraStats.TodayFollowups),
		ThisMonthLeads:   fmt.Sprintf("%d", extraStats.ThisMonthLeads),
		ProspectingPct:   fmt.Sprintf("%d", prospPct),
		CommittedPct:     fmt.Sprintf("%d", commitPct),
		EnrolledPct:      fmt.Sprintf("%d", enrollPct),
	}

	// Overdue candidates
	overdueData, err := model.GetAdminOverdueCandidates(ctx, 5)
	if err != nil {
		slog.Error("Failed to get overdue candidates", "error", err)
		overdueData = []model.OverdueCandidate{}
	}
	overdueCandidates := make([]admin.DashboardCandidate, len(overdueData))
	for i, c := range overdueData {
		overdueCandidates[i] = admin.DashboardCandidate{
			ID:     c.ID,
			Name:   c.Name,
			Detail: fmt.Sprintf("%d hari lalu", c.DaysOverdue),
		}
	}

	// Today's tasks
	todayData, err := model.GetAdminTodayTasks(ctx, 5)
	if err != nil {
		slog.Error("Failed to get today tasks", "error", err)
		todayData = []model.TodayTask{}
	}
	todayTasks := make([]admin.DashboardCandidate, len(todayData))
	for i, t := range todayData {
		todayTasks[i] = admin.DashboardCandidate{
			ID:     t.ID,
			Name:   t.Name,
			Detail: t.ProdiName,
		}
	}

	// Recent candidates
	recentData, err := model.GetRecentCandidates(ctx, 5)
	if err != nil {
		slog.Error("Failed to get recent candidates", "error", err)
		recentData = []model.RecentCandidate{}
	}
	recentCandidates := make([]admin.DashboardRecentCandidate, len(recentData))
	for i, c := range recentData {
		recentCandidates[i] = admin.DashboardRecentCandidate{
			ID:             c.ID,
			Name:           c.Name,
			ProdiName:      c.ProdiName,
			Status:         c.Status,
			ConsultantName: c.ConsultantName,
			RelativeDate:   formatRelativeTime(c.CreatedAt),
		}
	}

	admin.Dashboard(data, dashboardStats, overdueCandidates, todayTasks, recentCandidates).Render(ctx, w)
}

func (h *AdminHandler) handleConsultantDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Dashboard Saya")

	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	stats, err := model.GetConsultantStats(ctx, claims.UserID)
	if err != nil {
		log.Printf("Error getting consultant stats: %v", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	overdueCandidates, err := model.GetOverdueCandidates(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting overdue candidates: %v", err)
		overdueCandidates = []model.OverdueCandidate{}
	}

	todayTasks, err := model.GetTodayTasks(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting today tasks: %v", err)
		todayTasks = []model.TodayTask{}
	}

	suggestions, err := model.GetUnreadSuggestions(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting suggestions: %v", err)
		suggestions = []model.UnreadSuggestion{}
	}

	consultantName := claims.Name
	if consultantName == "" {
		consultantName = claims.Email
	}

	templateStats := admin.ConsultantDashboardStats{
		ConsultantName:      consultantName,
		TodayDate:           time.Now().Format("2 January 2006"),
		MyCandidates:        fmt.Sprintf("%d", stats.TotalCandidates),
		MyProspecting:       fmt.Sprintf("%d", stats.Prospecting),
		MyCommitted:         fmt.Sprintf("%d", stats.Committed),
		MyEnrolled:          fmt.Sprintf("%d", stats.Enrolled),
		OverdueCount:        fmt.Sprintf("%d", stats.OverdueCount),
		TodayTasks:          fmt.Sprintf("%d", stats.TodayTasks),
		UnreadSuggestions:   fmt.Sprintf("%d", stats.UnreadSuggestions),
		MonthlyNewLeads:     fmt.Sprintf("%d", stats.MonthlyNewLeads),
		MonthlyInteractions: fmt.Sprintf("%d", stats.MonthlyInteractions),
		MonthlyCommits:      fmt.Sprintf("%d", stats.MonthlyCommits),
		MonthlyEnrollments:  fmt.Sprintf("%d", stats.MonthlyEnrollments),
	}

	overdueList := make([]admin.CandidateSummary, len(overdueCandidates))
	for i, c := range overdueCandidates {
		lastContact := c.LastContact.Format("2 Jan 2006")
		overdueList[i] = admin.CandidateSummary{
			ID:          c.ID,
			Name:        c.Name,
			ProdiName:   c.ProdiName,
			WhatsApp:    formatWhatsApp(c.Phone),
			Status:      c.Status,
			LastContact: lastContact,
		}
	}

	taskList := make([]admin.CandidateSummary, len(todayTasks))
	for i, t := range todayTasks {
		taskList[i] = admin.CandidateSummary{
			ID:        t.ID,
			Name:      t.Name,
			ProdiName: t.ProdiName,
			WhatsApp:  formatWhatsApp(t.Phone),
			Status:    t.Status,
		}
	}

	suggestionList := make([]admin.SupervisorSuggestion, len(suggestions))
	for i, s := range suggestions {
		suggestionList[i] = admin.SupervisorSuggestion{
			ID:             s.ID,
			CandidateID:    s.CandidateID,
			CandidateName:  s.CandidateName,
			Suggestion:     s.Suggestion,
			SupervisorName: s.SupervisorName,
			CreatedAt:      s.CreatedAt.Format("2 Jan 2006"),
			IsRead:         false,
		}
	}

	admin.ConsultantDashboard(data, templateStats, overdueList, taskList, suggestionList).Render(ctx, w)
}

func (h *AdminHandler) handleSupervisorDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := NewPageDataWithUser(ctx, "Dashboard Supervisor")

	dashStats, err := model.GetSupervisorDashboardStats(ctx, claims.UserID)
	if err != nil {
		slog.Error("Failed to get supervisor dashboard stats", "error", err)
		http.Error(w, "Failed to load dashboard data", http.StatusInternalServerError)
		return
	}

	teamConsultantsData, err := model.GetTeamConsultants(ctx, claims.UserID)
	if err != nil {
		slog.Error("Failed to get team consultants", "error", err)
		http.Error(w, "Failed to load team data", http.StatusInternalServerError)
		return
	}

	stuckData, err := model.GetStuckCandidatesForTeam(ctx, claims.UserID, 10)
	if err != nil {
		slog.Error("Failed to get stuck candidates", "error", err)
		http.Error(w, "Failed to load stuck candidates", http.StatusInternalServerError)
		return
	}

	stats := admin.SupervisorDashboardStats{
		SupervisorName:  claims.Name,
		TodayDate:       formatDateIndonesian(time.Now()),
		TeamMemberCount: fmt.Sprintf("%d", len(teamConsultantsData)),
		TeamRegistered:  fmt.Sprintf("%d", dashStats.TeamRegistered),
		TeamProspecting: fmt.Sprintf("%d", dashStats.TeamProspecting),
		TeamCommitted:   fmt.Sprintf("%d", dashStats.TeamCommitted),
		TeamEnrolled:    fmt.Sprintf("%d", dashStats.TeamEnrolled),
		TeamLost:        fmt.Sprintf("%d", dashStats.TeamLost),
		StuckCandidates: fmt.Sprintf("%d", dashStats.StuckCandidates),
		TodayFollowups:  fmt.Sprintf("%d", dashStats.TodayFollowups),
		MonthlyNewLeads: fmt.Sprintf("%d", dashStats.MonthlyNewLeads),
		MonthlyEnrolled: fmt.Sprintf("%d", dashStats.MonthlyEnrolled),
	}

	teamConsultants := make([]admin.TeamConsultantItem, len(teamConsultantsData))
	for i, tc := range teamConsultantsData {
		teamConsultants[i] = admin.TeamConsultantItem{
			ID:             tc.ConsultantID,
			Name:           tc.ConsultantName,
			ActiveLeads:    fmt.Sprintf("%d", tc.ActiveLeads),
			TodayTasks:     fmt.Sprintf("%d", tc.TodayTasks),
			Overdue:        fmt.Sprintf("%d", tc.Overdue),
			MonthlyCommits: fmt.Sprintf("%d", tc.MonthlyCommits),
		}
	}

	stuckCandidates := make([]admin.StuckCandidateItem, len(stuckData))
	for i, sc := range stuckData {
		stuckCandidates[i] = admin.StuckCandidateItem{
			CandidateID:    sc.CandidateID,
			CandidateName:  sc.CandidateName,
			ProdiName:      sc.ProdiName,
			Status:         sc.Status,
			ConsultantName: sc.ConsultantName,
			DaysStuck:      fmt.Sprintf("%d", sc.DaysStuck),
		}
	}

	admin.SupervisorDashboard(data, stats, teamConsultants, stuckCandidates).Render(ctx, w)
}

func (h *AdminHandler) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	h.handleCampaignsSettings(w, r)
}

func (h *AdminHandler) handleReferrers(w http.ResponseWriter, r *http.Request) {
	h.handleReferrersSettings(w, r)
}

func formatDateIndonesian(t time.Time) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()], t.Year())
}

// formatWhatsApp formats phone number for WhatsApp link
func formatWhatsApp(phone string) string {
	if phone == "" {
		return ""
	}
	result := ""
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			result += string(c)
		}
	}
	if len(result) > 0 && result[0] == '0' {
		result = "62" + result[1:]
	}
	return result
}
