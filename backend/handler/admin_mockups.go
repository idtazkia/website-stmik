package handler

import (
	"log/slog"
	"net/http"

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
	data := NewPageDataWithUser(r.Context(), "Laporan Funnel")
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

func (h *AdminHandler) handleConsultantsReport(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kinerja Konsultan")
	filter := admin.ReportFilter{
		Period:    "this_month",
		StartDate: "2026-01-01",
		EndDate:   "2026-01-31",
	}
	consultants := []admin.ConsultantPerformance{
		{Rank: "1", Name: "Siti Rahayu", Email: "siti@kampus.edu", SupervisorName: "Dr. Ahmad", Leads: "45", Interactions: "120", Commits: "25", Enrollments: "12", ConversionRate: 26.7, ConversionRateStr: "26.7", AvgDaysToCommit: "14", InteractionsPerCandidate: "2.7"},
		{Rank: "2", Name: "Ahmad Fauzi", Email: "ahmad@kampus.edu", SupervisorName: "Dr. Ahmad", Leads: "38", Interactions: "95", Commits: "18", Enrollments: "10", ConversionRate: 26.3, ConversionRateStr: "26.3", AvgDaysToCommit: "12", InteractionsPerCandidate: "2.5"},
		{Rank: "3", Name: "Dewi Lestari", Email: "dewi@kampus.edu", SupervisorName: "Dr. Budi", Leads: "52", Interactions: "140", Commits: "22", Enrollments: "15", ConversionRate: 28.8, ConversionRateStr: "28.8", AvgDaysToCommit: "10", InteractionsPerCandidate: "2.7"},
	}
	summary := admin.ReportSummary{
		TotalLeads:        "135",
		TotalInteractions: "355",
		TotalCommits:      "65",
		TotalEnrollments:  "37",
	}
	admin.ConsultantReport(data, filter, consultants, summary).Render(r.Context(), w)
}

func (h *AdminHandler) handleCampaignsReport(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "ROI Kampanye")
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

// handleDocumentReview moved to admin_documents.go with real implementation
