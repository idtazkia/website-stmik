package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/mockdata"
	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// Mockup handlers - these use mock data and will be replaced with real implementations

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

func (h *AdminHandler) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kampanye")
	// Temporary: redirect to settings campaigns page
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
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Klaim Referral")

	// Get unverified claims
	claims, err := model.ListUnverifiedReferralClaims(ctx)
	if err != nil {
		slog.Error("Failed to list referral claims", "error", err)
		http.Error(w, "Failed to load referral claims", http.StatusInternalServerError)
		return
	}

	// Convert to template items
	claimItems := make([]admin.ReferralClaimItem, len(claims))
	for i, c := range claims {
		claimItems[i] = admin.ReferralClaimItem{
			CandidateID:   c.CandidateID,
			CandidateName: c.CandidateName,
			ProdiName:     c.ProdiName,
			SourceType:    c.SourceType,
			SourceDetail:  c.SourceDetail,
			Status:        c.Status,
			CreatedAt:     c.CreatedAt.Format("2 Jan 2006"),
		}
	}

	// Get all referrers for dropdown
	referrers, err := model.ListReferrers(ctx, "")
	if err != nil {
		slog.Error("Failed to list referrers", "error", err)
		http.Error(w, "Failed to load referrers", http.StatusInternalServerError)
		return
	}

	// Convert to template options
	referrerOptions := make([]admin.ReferrerOption, len(referrers))
	for i, r := range referrers {
		referrerOptions[i] = admin.ReferrerOption{
			ID:   r.ID,
			Name: r.Name,
			Type: r.Type,
		}
		if r.Institution != nil {
			referrerOptions[i].Institution = *r.Institution
		}
		if r.Code != nil {
			referrerOptions[i].Code = *r.Code
		}
	}

	admin.ReferralClaims(data, claimItems, referrerOptions).Render(ctx, w)
}

func (h *AdminHandler) handleLinkReferralClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseForm()

	candidateID := r.FormValue("candidate_id")
	referrerID := r.FormValue("referrer_id")
	mgmCode := r.FormValue("mgm_code")

	if candidateID == "" {
		http.Error(w, "Missing candidate_id", http.StatusBadRequest)
		return
	}

	if referrerID == "" && mgmCode == "" {
		http.Error(w, "Either referrer_id or mgm_code is required", http.StatusBadRequest)
		return
	}

	// Link to external referrer
	if referrerID != "" {
		err := model.LinkCandidateToReferrer(ctx, candidateID, referrerID)
		if err != nil {
			slog.Error("Failed to link candidate to referrer", "error", err)
			http.Error(w, "Failed to link referrer", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", "referralClaimLinked")
		w.Header().Set("HX-Refresh", "true")
		return
	}

	// Link to MGM referrer (enrolled student)
	if mgmCode != "" {
		// Find the candidate with this referral code
		referrerCandidate, err := model.FindCandidateByReferralCode(ctx, mgmCode)
		if err != nil {
			slog.Error("Failed to find MGM referrer", "error", err)
			http.Error(w, "Failed to find MGM referrer", http.StatusInternalServerError)
			return
		}
		if referrerCandidate == nil {
			http.Error(w, "MGM code not found", http.StatusBadRequest)
			return
		}

		err = model.LinkCandidateToMGMReferrer(ctx, candidateID, referrerCandidate.ID)
		if err != nil {
			slog.Error("Failed to link candidate to MGM referrer", "error", err)
			http.Error(w, "Failed to link MGM referrer", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", "referralClaimLinked")
		w.Header().Set("HX-Refresh", "true")
		return
	}
}

func (h *AdminHandler) handleInvalidReferralClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	err := model.ClearReferralClaim(ctx, candidateID)
	if err != nil {
		slog.Error("Failed to clear referral claim", "error", err)
		http.Error(w, "Failed to clear claim", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "referralClaimInvalid")
	w.Header().Set("HX-Refresh", "true")
}

func (h *AdminHandler) handleCommissions(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Komisi")
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
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Laporan Funnel")

	// Get funnel stats from database
	funnelStats, err := model.GetFunnelStats(ctx)
	if err != nil {
		slog.Error("Failed to get funnel stats", "error", err)
		http.Error(w, "Failed to load funnel data", http.StatusInternalServerError)
		return
	}

	// Calculate total for percentages (registered + prospecting + committed + enrolled)
	// We don't count lost as part of the funnel
	total := funnelStats.Registered + funnelStats.Prospecting + funnelStats.Committed + funnelStats.Enrolled
	if total == 0 {
		total = 1 // Avoid division by zero
	}

	// Build funnel stages
	stages := []admin.FunnelStage{
		{
			Name:       "Registered",
			Count:      fmt.Sprintf("%d", funnelStats.Registered),
			Percentage: "100",
			Color:      "bg-gray-500",
		},
		{
			Name:       "Prospecting",
			Count:      fmt.Sprintf("%d", funnelStats.Prospecting),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Prospecting, funnelStats.Registered)),
			Color:      "bg-blue-500",
		},
		{
			Name:       "Committed",
			Count:      fmt.Sprintf("%d", funnelStats.Committed),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Committed, funnelStats.Registered)),
			Color:      "bg-yellow-500",
		},
		{
			Name:       "Enrolled",
			Count:      fmt.Sprintf("%d", funnelStats.Enrolled),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Enrolled, funnelStats.Registered)),
			Color:      "bg-green-500",
		},
	}

	// Build conversion rates
	regToProsp := calcConversionRate(funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled, funnelStats.Registered+funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled)
	prospToCommit := calcConversionRate(funnelStats.Committed+funnelStats.Enrolled, funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled)
	commitToEnroll := calcConversionRate(funnelStats.Enrolled, funnelStats.Committed+funnelStats.Enrolled)

	conversions := []admin.FunnelConversion{
		{From: "Registered", To: "Prospecting", Rate: fmt.Sprintf("%.1f%%", regToProsp), Change: "-", IsPositive: true},
		{From: "Prospecting", To: "Committed", Rate: fmt.Sprintf("%.1f%%", prospToCommit), Change: "-", IsPositive: true},
		{From: "Committed", To: "Enrolled", Rate: fmt.Sprintf("%.1f%%", commitToEnroll), Change: "-", IsPositive: true},
	}

	admin.ReportFunnel(data, stages, conversions).Render(ctx, w)
}

func calcPercentage(part, total int) int {
	if total == 0 {
		return 0
	}
	return int(float64(part) / float64(total) * 100)
}

func calcConversionRate(converted, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(converted) / float64(total) * 100
}

func (h *AdminHandler) handleConsultantsReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Kinerja Konsultan")

	// Get consultant performance from database
	performanceData, err := model.GetConsultantPerformance(ctx, nil, nil)
	if err != nil {
		slog.Error("Failed to get consultant performance", "error", err)
		http.Error(w, "Failed to load consultant data", http.StatusInternalServerError)
		return
	}

	// Convert to template items
	consultants := make([]admin.ConsultantPerformance, len(performanceData))
	var totalLeads, totalInteractions, totalCommits, totalEnrollments int

	for i, cp := range performanceData {
		conversionRate := float64(0)
		if cp.TotalLeads > 0 {
			conversionRate = float64(cp.Enrollments) / float64(cp.TotalLeads) * 100
		}

		interactionsPerCandidate := float64(0)
		if cp.TotalLeads > 0 {
			interactionsPerCandidate = float64(cp.Interactions) / float64(cp.TotalLeads)
		}

		consultants[i] = admin.ConsultantPerformance{
			Rank:                     fmt.Sprintf("%d", i+1),
			Name:                     cp.ConsultantName,
			Email:                    cp.ConsultantEmail,
			SupervisorName:           cp.SupervisorName,
			Leads:                    fmt.Sprintf("%d", cp.TotalLeads),
			Interactions:             fmt.Sprintf("%d", cp.Interactions),
			Commits:                  fmt.Sprintf("%d", cp.Commits),
			Enrollments:              fmt.Sprintf("%d", cp.Enrollments),
			ConversionRate:           conversionRate,
			ConversionRateStr:        fmt.Sprintf("%.1f", conversionRate),
			AvgDaysToCommit:          fmt.Sprintf("%.0f", cp.AvgDaysToCommit),
			InteractionsPerCandidate: fmt.Sprintf("%.1f", interactionsPerCandidate),
		}

		totalLeads += cp.TotalLeads
		totalInteractions += cp.Interactions
		totalCommits += cp.Commits
		totalEnrollments += cp.Enrollments
	}

	filter := admin.ReportFilter{
		Period:    "all_time",
		StartDate: "",
		EndDate:   "",
	}

	summary := admin.ReportSummary{
		TotalLeads:        fmt.Sprintf("%d", totalLeads),
		TotalInteractions: fmt.Sprintf("%d", totalInteractions),
		TotalCommits:      fmt.Sprintf("%d", totalCommits),
		TotalEnrollments:  fmt.Sprintf("%d", totalEnrollments),
	}

	admin.ConsultantReport(data, filter, consultants, summary).Render(ctx, w)
}

func (h *AdminHandler) handleCampaignsReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "ROI Kampanye")

	// Get campaign stats from database
	campaignStats, err := model.GetCampaignStats(ctx)
	if err != nil {
		slog.Error("Failed to get campaign stats", "error", err)
		http.Error(w, "Failed to load campaign data", http.StatusInternalServerError)
		return
	}

	// Convert to template items
	campaigns := make([]admin.CampaignReportItem, len(campaignStats))
	var totalLeads, totalEnrolled int
	var bestCampaign string
	var bestConversion float64

	for i, cs := range campaignStats {
		totalCandidates := cs.Registered + cs.Prospecting + cs.Committed + cs.Enrolled
		conversion := float64(0)
		if totalCandidates > 0 {
			conversion = float64(cs.Enrolled) / float64(totalCandidates) * 100
		}

		campaigns[i] = admin.CampaignReportItem{
			Name:        cs.CampaignName,
			Type:        cs.CampaignType,
			Channel:     cs.Channel,
			Leads:       fmt.Sprintf("%d", totalCandidates),
			Prospecting: fmt.Sprintf("%d", cs.Prospecting),
			Committed:   fmt.Sprintf("%d", cs.Committed),
			Enrolled:    fmt.Sprintf("%d", cs.Enrolled),
			Conversion:  fmt.Sprintf("%.1f%%", conversion),
			Cost:        "-",        // Cost data not tracked in campaigns table
			CostPerLead: "-",        // Cost data not tracked
		}

		totalLeads += totalCandidates
		totalEnrolled += cs.Enrolled

		if conversion > bestConversion && totalCandidates > 0 {
			bestConversion = conversion
			bestCampaign = cs.CampaignName
		}
	}

	avgConversion := float64(0)
	if totalLeads > 0 {
		avgConversion = float64(totalEnrolled) / float64(totalLeads) * 100
	}

	summary := admin.CampaignReportSummary{
		TotalLeads:     fmt.Sprintf("%d", totalLeads),
		TotalEnrolled:  fmt.Sprintf("%d", totalEnrolled),
		AvgConversion:  fmt.Sprintf("%.1f%%", avgConversion),
		TotalCost:      "-",
		AvgCostPerLead: "-",
		BestCampaign:   bestCampaign,
	}

	admin.ReportCampaigns(data, campaigns, summary).Render(ctx, w)
}

func (h *AdminHandler) handleReferrersReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Leaderboard Referrer")

	// Get referrer stats from database
	referrerStats, err := model.GetReferrerStats(ctx)
	if err != nil {
		slog.Error("Failed to get referrer stats", "error", err)
		http.Error(w, "Failed to load referrer data", http.StatusInternalServerError)
		return
	}

	// Convert to template items
	referrers := make([]admin.ReferrerReportItem, len(referrerStats))
	var totalReferrals, totalEnrolled int
	var totalCommission int64
	var bestReferrer string
	var bestEnrolled int

	for i, rs := range referrerStats {
		conversion := float64(0)
		if rs.TotalReferrals > 0 {
			conversion = float64(rs.Enrolled) / float64(rs.TotalReferrals) * 100
		}

		referrers[i] = admin.ReferrerReportItem{
			Rank:           fmt.Sprintf("%d", i+1),
			Name:           rs.ReferrerName,
			Type:           rs.ReferrerType,
			Institution:    rs.Institution,
			TotalReferrals: fmt.Sprintf("%d", rs.TotalReferrals),
			Enrolled:       fmt.Sprintf("%d", rs.Enrolled),
			Pending:        fmt.Sprintf("%d", rs.Pending),
			Conversion:     fmt.Sprintf("%.1f%%", conversion),
			CommissionPaid: formatRupiah(rs.CommissionPaid),
		}

		totalReferrals += rs.TotalReferrals
		totalEnrolled += rs.Enrolled
		totalCommission += rs.CommissionPaid

		if rs.Enrolled > bestEnrolled {
			bestEnrolled = rs.Enrolled
			bestReferrer = rs.ReferrerName
		}
	}

	summary := admin.ReferrerReportSummary{
		TotalReferrers:  fmt.Sprintf("%d", len(referrerStats)),
		TotalReferrals:  fmt.Sprintf("%d", totalReferrals),
		TotalEnrolled:   fmt.Sprintf("%d", totalEnrolled),
		TotalCommission: formatRupiah(totalCommission),
		BestReferrer:    bestReferrer,
	}

	admin.ReportReferrers(data, referrers, summary).Render(ctx, w)
}

func (h *AdminHandler) handleConsultantDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Dashboard Saya")

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
		MonthlyCommits:      "2",
		MonthlyEnrollments:  "1",
	}

	overdueList := []admin.CandidateSummary{
		{ID: "1", Name: "Ahmad Pratama", ProdiName: "S1 Sistem Informasi", WhatsApp: "6281234567890", Status: "prospecting", LastContact: "10 Jan 2026"},
	}

	todayTasks := []admin.CandidateSummary{
		{ID: "2", Name: "Siti Rahayu", ProdiName: "S1 Teknik Informatika", WhatsApp: "6281234567891", Status: "prospecting", LastContact: "14 Jan 2026"},
		{ID: "3", Name: "Budi Santoso", ProdiName: "S1 Sistem Informasi", WhatsApp: "6281234567892", Status: "committed", LastContact: "13 Jan 2026"},
	}

	suggestions := []admin.SupervisorSuggestion{
		{ID: "1", CandidateID: "1", CandidateName: "Ahmad Pratama", Suggestion: "Coba tawarkan program beasiswa untuk menarik minatnya.", SupervisorName: "Dr. Ahmad Fauzi", CreatedAt: "14 Jan 2026", IsRead: false},
	}

	admin.ConsultantDashboard(data, stats, overdueList, todayTasks, suggestions).Render(r.Context(), w)
}

func (h *AdminHandler) handleSupervisorDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := NewPageDataWithUser(ctx, "Dashboard Supervisor")

	// Get supervisor dashboard stats
	dashStats, err := model.GetSupervisorDashboardStats(ctx, claims.UserID)
	if err != nil {
		slog.Error("Failed to get supervisor dashboard stats", "error", err)
		http.Error(w, "Failed to load dashboard data", http.StatusInternalServerError)
		return
	}

	// Get team consultants
	teamConsultantsData, err := model.GetTeamConsultants(ctx, claims.UserID)
	if err != nil {
		slog.Error("Failed to get team consultants", "error", err)
		http.Error(w, "Failed to load team data", http.StatusInternalServerError)
		return
	}

	// Get stuck candidates (limit to 10)
	stuckData, err := model.GetStuckCandidatesForTeam(ctx, claims.UserID, 10)
	if err != nil {
		slog.Error("Failed to get stuck candidates", "error", err)
		http.Error(w, "Failed to load stuck candidates", http.StatusInternalServerError)
		return
	}

	// Convert to template types
	stats := admin.SupervisorDashboardStats{
		SupervisorName:   claims.Name,
		TodayDate:        formatDateIndonesian(time.Now()),
		TeamMemberCount:  fmt.Sprintf("%d", len(teamConsultantsData)),
		TeamRegistered:   fmt.Sprintf("%d", dashStats.TeamRegistered),
		TeamProspecting:  fmt.Sprintf("%d", dashStats.TeamProspecting),
		TeamCommitted:    fmt.Sprintf("%d", dashStats.TeamCommitted),
		TeamEnrolled:     fmt.Sprintf("%d", dashStats.TeamEnrolled),
		TeamLost:         fmt.Sprintf("%d", dashStats.TeamLost),
		StuckCandidates:  fmt.Sprintf("%d", dashStats.StuckCandidates),
		TodayFollowups:   fmt.Sprintf("%d", dashStats.TodayFollowups),
		MonthlyNewLeads:  fmt.Sprintf("%d", dashStats.MonthlyNewLeads),
		MonthlyEnrolled:  fmt.Sprintf("%d", dashStats.MonthlyEnrolled),
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

func formatDateIndonesian(t time.Time) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()], t.Year())
}

// handleDocumentReview moved to admin_documents.go with real implementation
