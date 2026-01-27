package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleConsultantDashboardReal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Dashboard Saya")

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Get consultant stats
	stats, err := model.GetConsultantStats(ctx, claims.UserID)
	if err != nil {
		log.Printf("Error getting consultant stats: %v", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	// Get overdue candidates
	overdueCandidates, err := model.GetOverdueCandidates(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting overdue candidates: %v", err)
		overdueCandidates = []model.OverdueCandidate{}
	}

	// Get today's tasks
	todayTasks, err := model.GetTodayTasks(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting today tasks: %v", err)
		todayTasks = []model.TodayTask{}
	}

	// Get unread suggestions
	suggestions, err := model.GetUnreadSuggestions(ctx, claims.UserID, 10)
	if err != nil {
		log.Printf("Error getting suggestions: %v", err)
		suggestions = []model.UnreadSuggestion{}
	}

	// Get consultant name
	consultantName := claims.Name
	if consultantName == "" {
		consultantName = claims.Email
	}

	// Convert to template types
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
			IsRead:         false, // These are all unread
		}
	}

	admin.ConsultantDashboard(data, templateStats, overdueList, taskList, suggestionList).Render(ctx, w)
}

// formatWhatsApp formats phone number for WhatsApp link
func formatWhatsApp(phone string) string {
	if phone == "" {
		return ""
	}
	// Remove non-digit characters
	result := ""
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			result += string(c)
		}
	}
	// Add country code if missing
	if len(result) > 0 && result[0] == '0' {
		result = "62" + result[1:]
	}
	return result
}
