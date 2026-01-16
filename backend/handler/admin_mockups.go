package handler

import (
	"net/http"

	"github.com/idtazkia/stmik-admission-api/mockdata"
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

func (h *AdminHandler) handleCandidates(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kandidat")

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

	data := NewPageDataWithUser(r.Context(), "Detail Kandidat - "+candidate.Name)

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
	data := NewPageDataWithUser(r.Context(), "Klaim Referral")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Klaim Referral - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
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

func (h *AdminHandler) handleInteractionForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	candidate := mockdata.GetCandidateByID(id)
	if candidate == nil {
		http.NotFound(w, r)
		return
	}

	data := NewPageDataWithUser(r.Context(), "Catat Interaksi - "+candidate.Name)

	c := admin.CandidateSummary{
		ID:        candidate.ID,
		Name:      candidate.Name,
		ProdiName: candidate.ProdiName,
		WhatsApp:  candidate.WhatsApp,
		Status:    candidate.Status,
	}

	categories := []admin.InteractionCategoryOption{
		{Value: "interested", Label: "Tertarik", Icon: "😊", Sentiment: "positive"},
		{Value: "considering", Label: "Mempertimbangkan", Icon: "🤔", Sentiment: "neutral"},
		{Value: "hesitant", Label: "Ragu-ragu", Icon: "😕", Sentiment: "neutral"},
		{Value: "cold", Label: "Dingin", Icon: "😐", Sentiment: "negative"},
		{Value: "unreachable", Label: "Tidak bisa dihubungi", Icon: "📵", Sentiment: "negative"},
	}

	obstacles := []admin.ObstacleOption{
		{Value: "expensive", Label: "Biaya mahal"},
		{Value: "far", Label: "Lokasi jauh"},
		{Value: "parent_not_agreed", Label: "Orang tua belum setuju"},
		{Value: "bad_timing", Label: "Waktu belum tepat"},
		{Value: "other_campus", Label: "Memilih kampus lain"},
	}

	admin.InteractionForm(data, c, categories, obstacles).Render(r.Context(), w)
}

func (h *AdminHandler) handleDocumentReview(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Review Dokumen")

	filter := admin.DocumentFilter{
		Status: r.URL.Query().Get("status"),
		Type:   r.URL.Query().Get("type"),
		Search: r.URL.Query().Get("search"),
	}

	documents := []admin.DocumentReviewItem{
		{ID: "1", CandidateID: "1", CandidateName: "Ahmad Pratama", ProdiName: "S1 Sistem Informasi", Type: "ktp", TypeName: "KTP", FileName: "ktp_ahmad.jpg", FileSize: "1.2 MB", FileURL: "#", ThumbnailURL: "#", IsImage: true, Status: "pending", UploadedAt: "14 Jan 2026"},
		{ID: "2", CandidateID: "2", CandidateName: "Siti Rahayu", ProdiName: "S1 Teknik Informatika", Type: "ijazah", TypeName: "Ijazah", FileName: "ijazah_siti.pdf", FileSize: "2.5 MB", FileURL: "#", ThumbnailURL: "", IsImage: false, Status: "approved", UploadedAt: "13 Jan 2026"},
		{ID: "3", CandidateID: "3", CandidateName: "Budi Santoso", ProdiName: "S1 Sistem Informasi", Type: "photo", TypeName: "Foto", FileName: "foto_budi.jpg", FileSize: "800 KB", FileURL: "#", ThumbnailURL: "#", IsImage: true, Status: "rejected", RejectionReason: "Foto tidak jelas", UploadedAt: "12 Jan 2026"},
	}

	stats := admin.DocumentStats{
		Pending:       "5",
		ApprovedToday: "3",
		RejectedToday: "1",
		Total:         "53",
	}

	admin.DocumentReviewList(data, filter, documents, stats).Render(r.Context(), w)
}
