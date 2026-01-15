package handler

import (
	"net/http"

	"github.com/idtazkia/stmik-admission-api/mockdata"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

// AdminHandler handles all admin routes
type AdminHandler struct{}

// NewAdminHandler creates a new admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// RegisterRoutes registers all admin routes to the mux
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	// Dashboard
	mux.HandleFunc("GET /admin", h.handleDashboard)
	mux.HandleFunc("GET /admin/", h.handleDashboard)

	// Consultant personal dashboard
	mux.HandleFunc("GET /admin/my-dashboard", h.handleConsultantDashboard)

	// Candidates
	mux.HandleFunc("GET /admin/candidates", h.handleCandidates)
	mux.HandleFunc("GET /admin/candidates/{id}", h.handleCandidateDetail)
	mux.HandleFunc("GET /admin/candidates/{id}/interaction", h.handleInteractionForm)

	// Documents
	mux.HandleFunc("GET /admin/documents", h.handleDocumentReview)

	// Marketing
	mux.HandleFunc("GET /admin/campaigns", h.handleCampaigns)
	mux.HandleFunc("GET /admin/referrers", h.handleReferrers)
	mux.HandleFunc("GET /admin/referral-claims", h.handleReferralClaims)
	mux.HandleFunc("GET /admin/commissions", h.handleCommissions)

	// Reports
	mux.HandleFunc("GET /admin/reports/funnel", h.handleFunnelReport)
	mux.HandleFunc("GET /admin/reports/consultants", h.handleConsultantsReport)
	mux.HandleFunc("GET /admin/reports/campaigns", h.handleCampaignsReport)

	// Settings
	mux.HandleFunc("GET /admin/settings/users", h.handleUsersSettings)
	mux.HandleFunc("GET /admin/settings/programs", h.handleProgramsSettings)
	mux.HandleFunc("GET /admin/settings/fees", h.handleFeesSettings)
	mux.HandleFunc("GET /admin/settings/rewards", h.handleRewardsSettings)
	mux.HandleFunc("GET /admin/settings/categories", h.handleCategoriesSettings)

	// Auth
	mux.HandleFunc("GET /admin/login", h.handleLogin)
	mux.HandleFunc("POST /admin/login", h.handleLoginSubmit)
	mux.HandleFunc("GET /admin/logout", h.handleLogout)
}

func (h *AdminHandler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Dashboard")
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
	data := NewPageData("Kandidat")

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

	data := NewPageData("Detail Kandidat - " + candidate.Name)

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
	data := NewPageData("Kampanye")
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
	data := NewPageData("Referrer")
	referrers := []admin.ReferrerItem{
		{ID: "1", Name: "Pak Ahmad Fauzi", Type: "guru", Institution: "SMAN 1 Bogor", Phone: "081234567890", Code: "REF-AF01", Commission: "Rp 750.000", Referrals: "8", Enrolled: "5", TotalEarned: "Rp 3.750.000", Status: "active"},
		{ID: "2", Name: "Siti Nurhaliza", Type: "alumni", Institution: "STMIK Tazkia 2022", Phone: "081234567891", Code: "REF-SN02", Commission: "Rp 500.000", Referrals: "4", Enrolled: "2", TotalEarned: "Rp 1.000.000", Status: "active"},
		{ID: "3", Name: "PT Edutech Indonesia", Type: "partner", Institution: "Bimbel Edutech", Phone: "021-7654321", Code: "REF-EDU", Commission: "Rp 1.000.000", Referrals: "12", Enrolled: "6", TotalEarned: "Rp 6.000.000", Status: "active"},
		{ID: "4", Name: "Budi Santoso", Type: "staff", Institution: "STMIK Tazkia", Phone: "081234567893", Code: "REF-BS04", Commission: "Rp 250.000", Referrals: "3", Enrolled: "2", TotalEarned: "Rp 500.000", Status: "active"},
	}
	admin.SettingsReferrers(data, referrers).Render(r.Context(), w)
}

func (h *AdminHandler) handleReferralClaims(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Klaim Referral")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Klaim Referral - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleCommissions(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Komisi")
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
	data := NewPageData("Laporan Funnel")
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
	data := NewPageData("ROI Kampanye")
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
	data := NewPageData("Users")
	users := []admin.UserItem{
		{ID: "1", Name: "Admin PMB", Email: "admin@tazkia.ac.id", Role: "admin", Supervisor: "-", Status: "active", LastLogin: "Hari ini"},
		{ID: "2", Name: "Budi Santoso", Email: "budi@tazkia.ac.id", Role: "supervisor", Supervisor: "-", Status: "active", LastLogin: "Hari ini"},
		{ID: "3", Name: "Siti Rahayu", Email: "siti@tazkia.ac.id", Role: "konsultan", Supervisor: "Budi Santoso", Status: "active", LastLogin: "Hari ini"},
		{ID: "4", Name: "Ahmad Hidayat", Email: "ahmad@tazkia.ac.id", Role: "konsultan", Supervisor: "Budi Santoso", Status: "active", LastLogin: "Kemarin"},
		{ID: "5", Name: "Dewi Lestari", Email: "dewi@tazkia.ac.id", Role: "konsultan", Supervisor: "Budi Santoso", Status: "active", LastLogin: "2 hari lalu"},
	}
	admin.SettingsUsers(data, users).Render(r.Context(), w)
}

func (h *AdminHandler) handleProgramsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Prodi")
	programs := []admin.ProgramItem{
		{ID: "1", Name: "Sistem Informasi", Code: "SI", Level: "S1", SPPFee: "Rp 7.500.000", Status: "active", Students: "245"},
		{ID: "2", Name: "Teknik Informatika", Code: "TI", Level: "S1", SPPFee: "Rp 8.000.000", Status: "active", Students: "312"},
	}
	admin.SettingsPrograms(data, programs).Render(r.Context(), w)
}

func (h *AdminHandler) handleFeesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Biaya")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Biaya - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleRewardsSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Reward")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h1>Reward - Coming Soon</h1><a href="/admin">Back to Dashboard</a></body></html>`))
	_ = data
}

func (h *AdminHandler) handleCategoriesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Kategori")
	categories := []admin.CategoryItem{
		{ID: "1", Name: "Tertarik", Sentiment: "positive", Count: "125"},
		{ID: "2", Name: "Mempertimbangkan", Sentiment: "neutral", Count: "89"},
		{ID: "3", Name: "Ragu-ragu", Sentiment: "neutral", Count: "45"},
		{ID: "4", Name: "Dingin", Sentiment: "negative", Count: "32"},
		{ID: "5", Name: "Tidak bisa dihubungi", Sentiment: "negative", Count: "28"},
	}
	obstacles := []admin.ObstacleItem{
		{ID: "1", Name: "Biaya terlalu mahal", Count: "67"},
		{ID: "2", Name: "Orang tua belum setuju", Count: "45"},
		{ID: "3", Name: "Lokasi jauh", Count: "23"},
		{ID: "4", Name: "Memilih kampus lain", Count: "18"},
		{ID: "5", Name: "Waktu belum tepat", Count: "12"},
	}
	admin.SettingsCategories(data, categories, obstacles).Render(r.Context(), w)
}

func (h *AdminHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Login")
	admin.Login(data, "").Render(r.Context(), w)
}

func (h *AdminHandler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	// Mock login - just redirect to dashboard
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Mock logout - just redirect to login
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func (h *AdminHandler) handleConsultantDashboard(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("Dashboard Saya")

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
	data := NewPageData("Performa Konsultan")

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

	data := NewPageData("Log Interaksi - " + candidate.Name)

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
	data := NewPageData("Review Dokumen")

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
