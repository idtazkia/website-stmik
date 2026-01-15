// Package mockdata provides hardcoded dummy data for UI mockup
package mockdata

import (
	"fmt"
	"strings"
	"time"
)

// User represents a staff user
type User struct {
	ID           int
	Email        string
	Name         string
	Role         string // admin, supervisor, consultant
	SupervisorID *int
	IsActive     bool
}

// Prodi represents a study program
type Prodi struct {
	ID       int
	Name     string
	Code     string
	Degree   string // S1, D3
	IsActive bool
}

// Candidate represents a prospective student
type Candidate struct {
	ID                   int
	Name                 string
	Email                string
	Phone                string
	WhatsApp             string
	Address              string
	City                 string
	Province             string
	HighSchool           string
	GraduationYear       int
	ProdiID              int
	ProdiName            string
	SourceType           string
	SourceDetail         string
	CampaignName         string
	ReferrerName         string
	Status               string // registered, prospecting, committed, enrolled, lost
	AssignedConsultantID int
	ConsultantName       string
	RegistrationFeePaid  bool
	CreatedAt            time.Time
	LastInteractionAt    *time.Time
	NextFollowupAt       *time.Time
}

// Interaction represents a contact log
type Interaction struct {
	ID                   int
	CandidateID          int
	ConsultantID         int
	ConsultantName       string
	Channel              string // call, whatsapp, email, campus_visit, home_visit
	Category             string // interested, considering, hesitant, cold, unreachable
	CategorySentiment    string // positive, neutral, negative
	Obstacle             string
	Remarks              string
	NextFollowupDate     *time.Time
	SupervisorSuggestion string
	SuggestionReadAt     *time.Time
	CreatedAt            time.Time
}

// Campaign represents a marketing campaign
type Campaign struct {
	ID                      int
	Name                    string
	Type                    string // promo, event, ads
	SourceChannel           string
	StartDate               time.Time
	EndDate                 time.Time
	RegistrationFeeOverride *int
	IsActive                bool
	LeadsCount              int
	EnrollmentsCount        int
}

// Referrer represents someone who refers candidates
type Referrer struct {
	ID                  int
	Name                string
	Type                string // alumni, teacher, student, partner, staff
	Institution         string
	Phone               string
	Email               string
	Code                string
	BankName            string
	BankAccount         string
	CommissionOverride  *int
	IsActive            bool
	ReferralsCount      int
	EnrollmentsCount    int
	TotalCommission     int
	PaidCommission      int
}

// FeeStructure represents fee configuration
type FeeStructure struct {
	ID           int
	FeeType      string // registration, tuition, dormitory
	ProdiName    string
	AcademicYear string
	Amount       int
	IsActive     bool
}

// RewardConfig represents reward configuration by referrer type
type RewardConfig struct {
	ID           int
	ReferrerType string
	RewardType   string // cash, tuition_discount, merchandise
	Amount       int
	IsPercentage bool
	TriggerEvent string // enrollment, commitment
	IsActive     bool
}

// Billing represents a bill for a candidate
type Billing struct {
	ID               int
	CandidateID      int
	FeeType          string
	AcademicYear     string
	Semester         int
	TotalAmount      int
	PaidAmount       int
	InstallmentCount int
	Status           string // pending, partial, paid, cancelled
}

// Document represents an uploaded document
type Document struct {
	ID              int
	CandidateID     int
	DocumentType    string // ktp, photo, ijazah, transcript
	FileName        string
	Status          string // pending, approved, rejected
	RejectionReason string
	CreatedAt       time.Time
}

// Announcement represents an announcement for candidates
type Announcement struct {
	ID           int
	Title        string
	Content      string
	TargetStatus string
	TargetProdi  string
	PublishedAt  time.Time
	IsActive     bool
}

// InteractionCategory represents interaction category
type InteractionCategory struct {
	ID        int
	Name      string
	Sentiment string // positive, neutral, negative
	IsActive  bool
}

// Obstacle represents common obstacles
type Obstacle struct {
	ID                int
	Name              string
	SuggestedResponse string
	IsActive          bool
}

// Stats for dashboard
type DashboardStats struct {
	TotalCandidates    string
	RegisteredCount    string
	ProspectingCount   string
	CommittedCount     string
	EnrolledCount      string
	LostCount          string
	OverdueFollowups   string
	TodayFollowups     string
	ConversionRate     float64
	ThisMonthLeads     string
	LastMonthLeads     string
}

// ConsultantStats for consultant dashboard
type ConsultantStats struct {
	MyCandidates      int
	MyProspecting     int
	MyCommitted       int
	MyEnrolled        int
	MyLost            int
	OverdueFollowups  int
	TodayFollowups    int
	UnreadSuggestions int
}

// Dummy data instances
var (
	Users = []User{
		{ID: 1, Email: "admin@tazkia.ac.id", Name: "Admin Utama", Role: "admin", IsActive: true},
		{ID: 2, Email: "supervisor1@tazkia.ac.id", Name: "Budi Santoso", Role: "supervisor", IsActive: true},
		{ID: 3, Email: "konsultan1@tazkia.ac.id", Name: "Siti Rahayu", Role: "consultant", SupervisorID: intPtr(2), IsActive: true},
		{ID: 4, Email: "konsultan2@tazkia.ac.id", Name: "Ahmad Hidayat", Role: "consultant", SupervisorID: intPtr(2), IsActive: true},
		{ID: 5, Email: "konsultan3@tazkia.ac.id", Name: "Dewi Lestari", Role: "consultant", SupervisorID: intPtr(2), IsActive: false},
	}

	Prodis = []Prodi{
		{ID: 1, Name: "Sistem Informasi", Code: "SI", Degree: "S1", IsActive: true},
		{ID: 2, Name: "Teknik Informatika", Code: "TI", Degree: "S1", IsActive: true},
	}

	Candidates = []Candidate{
		{ID: 1, Name: "Muhammad Rizky", Email: "rizky@gmail.com", Phone: "081234567890", WhatsApp: "081234567890", Address: "Jl. Merdeka No. 1", City: "Bogor", Province: "Jawa Barat", HighSchool: "SMAN 1 Bogor", GraduationYear: 2024, ProdiID: 1, ProdiName: "Sistem Informasi", SourceType: "instagram", Status: "registered", AssignedConsultantID: 3, ConsultantName: "Siti Rahayu", RegistrationFeePaid: false, CreatedAt: time.Now().AddDate(0, 0, -2), NextFollowupAt: timePtr(time.Now().AddDate(0, 0, 1))},
		{ID: 2, Name: "Putri Amelia", Email: "putri@gmail.com", Phone: "081234567891", WhatsApp: "081234567891", Address: "Jl. Sudirman No. 5", City: "Jakarta", Province: "DKI Jakarta", HighSchool: "SMAN 8 Jakarta", GraduationYear: 2024, ProdiID: 2, ProdiName: "Teknik Informatika", SourceType: "expo", SourceDetail: "Education Expo Jakarta", CampaignName: "Expo Jakarta 2025", Status: "prospecting", AssignedConsultantID: 3, ConsultantName: "Siti Rahayu", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, 0, -5), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -1)), NextFollowupAt: timePtr(time.Now())},
		{ID: 3, Name: "Dimas Pratama", Email: "dimas@gmail.com", Phone: "081234567892", WhatsApp: "081234567892", Address: "Jl. Asia Afrika No. 10", City: "Bandung", Province: "Jawa Barat", HighSchool: "SMAN 3 Bandung", GraduationYear: 2024, ProdiID: 1, ProdiName: "Sistem Informasi", SourceType: "friend_family", SourceDetail: "Direferensikan oleh Pak Asep", ReferrerName: "Asep Supriatna", Status: "committed", AssignedConsultantID: 4, ConsultantName: "Ahmad Hidayat", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, 0, -10), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -2))},
		{ID: 4, Name: "Anisa Fitri", Email: "anisa@gmail.com", Phone: "081234567893", WhatsApp: "081234567893", Address: "Jl. Diponegoro No. 15", City: "Depok", Province: "Jawa Barat", HighSchool: "SMAN 1 Depok", GraduationYear: 2024, ProdiID: 2, ProdiName: "Teknik Informatika", SourceType: "google", Status: "enrolled", AssignedConsultantID: 4, ConsultantName: "Ahmad Hidayat", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, -1, 0), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -7))},
		{ID: 5, Name: "Fajar Nugroho", Email: "fajar@gmail.com", Phone: "081234567894", WhatsApp: "081234567894", Address: "Jl. Gatot Subroto No. 20", City: "Bekasi", Province: "Jawa Barat", HighSchool: "SMAN 2 Bekasi", GraduationYear: 2024, ProdiID: 1, ProdiName: "Sistem Informasi", SourceType: "tiktok", Status: "lost", AssignedConsultantID: 3, ConsultantName: "Siti Rahayu", RegistrationFeePaid: false, CreatedAt: time.Now().AddDate(0, 0, -15), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -10))},
		{ID: 6, Name: "Rina Wulandari", Email: "rina@gmail.com", Phone: "081234567895", WhatsApp: "081234567895", Address: "Jl. Veteran No. 8", City: "Tangerang", Province: "Banten", HighSchool: "SMAN 1 Tangerang", GraduationYear: 2025, ProdiID: 2, ProdiName: "Teknik Informatika", SourceType: "school_visit", CampaignName: "Kunjungan Sekolah Q1", Status: "prospecting", AssignedConsultantID: 3, ConsultantName: "Siti Rahayu", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, 0, -3), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -1)), NextFollowupAt: timePtr(time.Now().AddDate(0, 0, 2))},
		{ID: 7, Name: "Bayu Setiawan", Email: "bayu@gmail.com", Phone: "081234567896", WhatsApp: "081234567896", Address: "Jl. Pemuda No. 12", City: "Bogor", Province: "Jawa Barat", HighSchool: "SMK 1 Bogor", GraduationYear: 2024, ProdiID: 1, ProdiName: "Sistem Informasi", SourceType: "teacher_alumni", SourceDetail: "Direferensikan Bu Yanti (Guru BK)", Status: "registered", AssignedConsultantID: 4, ConsultantName: "Ahmad Hidayat", RegistrationFeePaid: false, CreatedAt: time.Now().AddDate(0, 0, -1), NextFollowupAt: timePtr(time.Now())},
		{ID: 8, Name: "Citra Dewi", Email: "citra@gmail.com", Phone: "081234567897", WhatsApp: "081234567897", Address: "Jl. Raya Bogor No. 100", City: "Jakarta", Province: "DKI Jakarta", HighSchool: "SMAN 70 Jakarta", GraduationYear: 2025, ProdiID: 2, ProdiName: "Teknik Informatika", SourceType: "instagram", CampaignName: "Instagram Ads Q1", Status: "prospecting", AssignedConsultantID: 4, ConsultantName: "Ahmad Hidayat", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, 0, -4), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -3)), NextFollowupAt: timePtr(time.Now().AddDate(0, 0, -1))},
		{ID: 9, Name: "Eko Prasetyo", Email: "eko@gmail.com", Phone: "081234567898", WhatsApp: "081234567898", Address: "Jl. Margonda No. 50", City: "Depok", Province: "Jawa Barat", HighSchool: "SMAN 2 Depok", GraduationYear: 2024, ProdiID: 1, ProdiName: "Sistem Informasi", SourceType: "walkin", Status: "committed", AssignedConsultantID: 3, ConsultantName: "Siti Rahayu", RegistrationFeePaid: true, CreatedAt: time.Now().AddDate(0, 0, -8), LastInteractionAt: timePtr(time.Now().AddDate(0, 0, -2)), NextFollowupAt: timePtr(time.Now().AddDate(0, 0, 3))},
		{ID: 10, Name: "Fina Rahmawati", Email: "fina@gmail.com", Phone: "081234567899", WhatsApp: "081234567899", Address: "Jl. Pahlawan No. 25", City: "Bandung", Province: "Jawa Barat", HighSchool: "SMAN 5 Bandung", GraduationYear: 2025, ProdiID: 2, ProdiName: "Teknik Informatika", SourceType: "youtube", Status: "registered", AssignedConsultantID: 4, ConsultantName: "Ahmad Hidayat", RegistrationFeePaid: false, CreatedAt: time.Now(), NextFollowupAt: timePtr(time.Now().AddDate(0, 0, 1))},
	}

	Interactions = []Interaction{
		{ID: 1, CandidateID: 2, ConsultantID: 3, ConsultantName: "Siti Rahayu", Channel: "whatsapp", Category: "interested", CategorySentiment: "positive", Remarks: "Kandidat sangat tertarik dengan prodi TI. Sudah lihat website dan fasilitas.", NextFollowupDate: timePtr(time.Now()), CreatedAt: time.Now().AddDate(0, 0, -1)},
		{ID: 2, CandidateID: 2, ConsultantID: 3, ConsultantName: "Siti Rahayu", Channel: "call", Category: "considering", CategorySentiment: "neutral", Obstacle: "Biaya", Remarks: "Masih mempertimbangkan biaya. Sudah dijelaskan skema cicilan.", SupervisorSuggestion: "Tawarkan promo early bird untuk diskon biaya pendaftaran", CreatedAt: time.Now().AddDate(0, 0, -3)},
		{ID: 3, CandidateID: 3, ConsultantID: 4, ConsultantName: "Ahmad Hidayat", Channel: "campus_visit", Category: "interested", CategorySentiment: "positive", Remarks: "Kandidat dan orang tua visit ke kampus. Sangat impressed dengan fasilitas lab.", CreatedAt: time.Now().AddDate(0, 0, -2)},
		{ID: 4, CandidateID: 3, ConsultantID: 4, ConsultantName: "Ahmad Hidayat", Channel: "whatsapp", Category: "interested", CategorySentiment: "positive", Remarks: "Konfirmasi akan mendaftar minggu depan. Sudah siapkan dokumen.", CreatedAt: time.Now().AddDate(0, 0, -5)},
		{ID: 5, CandidateID: 5, ConsultantID: 3, ConsultantName: "Siti Rahayu", Channel: "call", Category: "cold", CategorySentiment: "negative", Obstacle: "Kompetitor", Remarks: "Kandidat memilih kampus lain yang lebih dekat rumah.", CreatedAt: time.Now().AddDate(0, 0, -10)},
	}

	Campaigns = []Campaign{
		{ID: 1, Name: "Promo Early Bird 2025", Type: "promo", SourceChannel: "all", StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now().AddDate(0, 1, 0), RegistrationFeeOverride: intPtr(0), IsActive: true, LeadsCount: 25, EnrollmentsCount: 5},
		{ID: 2, Name: "Expo Jakarta 2025", Type: "event", SourceChannel: "expo", StartDate: time.Now().AddDate(0, 0, -15), EndDate: time.Now().AddDate(0, 0, -14), IsActive: false, LeadsCount: 50, EnrollmentsCount: 8},
		{ID: 3, Name: "Instagram Ads Q1", Type: "ads", SourceChannel: "instagram", StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now().AddDate(0, 2, 0), IsActive: true, LeadsCount: 30, EnrollmentsCount: 3},
		{ID: 4, Name: "Kunjungan Sekolah Q1", Type: "event", SourceChannel: "school_visit", StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now().AddDate(0, 2, 0), IsActive: true, LeadsCount: 40, EnrollmentsCount: 10},
	}

	Referrers = []Referrer{
		{ID: 1, Name: "Asep Supriatna", Type: "teacher", Institution: "SMAN 1 Bogor", Phone: "081111111111", Email: "asep@gmail.com", Code: "REF-ASEP", BankName: "BCA", BankAccount: "1234567890", IsActive: true, ReferralsCount: 5, EnrollmentsCount: 2, TotalCommission: 1500000, PaidCommission: 750000},
		{ID: 2, Name: "Yanti Susanti", Type: "teacher", Institution: "SMK 1 Bogor", Phone: "081222222222", Email: "yanti@gmail.com", Code: "REF-YANTI", BankName: "Mandiri", BankAccount: "0987654321", IsActive: true, ReferralsCount: 3, EnrollmentsCount: 1, TotalCommission: 750000, PaidCommission: 0},
		{ID: 3, Name: "Andi Wijaya", Type: "alumni", Institution: "STMIK Tazkia", Phone: "081333333333", Email: "andi@gmail.com", Code: "REF-ANDI", BankName: "BNI", BankAccount: "1122334455", IsActive: true, ReferralsCount: 8, EnrollmentsCount: 4, TotalCommission: 2000000, PaidCommission: 2000000},
		{ID: 4, Name: "Bimbel Cerdas", Type: "partner", Institution: "Bimbel Cerdas Bogor", Phone: "081444444444", Email: "bimbel@gmail.com", Code: "REF-BIMBEL", BankName: "BCA", BankAccount: "5566778899", CommissionOverride: intPtr(1000000), IsActive: true, ReferralsCount: 15, EnrollmentsCount: 7, TotalCommission: 7000000, PaidCommission: 5000000},
	}

	FeeStructures = []FeeStructure{
		{ID: 1, FeeType: "registration", ProdiName: "Semua Prodi", AcademicYear: "2025/2026", Amount: 500000, IsActive: true},
		{ID: 2, FeeType: "tuition", ProdiName: "Sistem Informasi", AcademicYear: "2025/2026", Amount: 7500000, IsActive: true},
		{ID: 3, FeeType: "tuition", ProdiName: "Teknik Informatika", AcademicYear: "2025/2026", Amount: 8000000, IsActive: true},
		{ID: 4, FeeType: "dormitory", ProdiName: "Semua Prodi", AcademicYear: "2025/2026", Amount: 12000000, IsActive: true},
	}

	RewardConfigs = []RewardConfig{
		{ID: 1, ReferrerType: "alumni", RewardType: "cash", Amount: 500000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 2, ReferrerType: "teacher", RewardType: "cash", Amount: 750000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 3, ReferrerType: "student", RewardType: "cash", Amount: 300000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 4, ReferrerType: "partner", RewardType: "cash", Amount: 1000000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 5, ReferrerType: "staff", RewardType: "cash", Amount: 250000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 6, ReferrerType: "mgm_referrer", RewardType: "cash", Amount: 200000, IsPercentage: false, TriggerEvent: "enrollment", IsActive: true},
		{ID: 7, ReferrerType: "mgm_referee", RewardType: "tuition_discount", Amount: 10, IsPercentage: true, TriggerEvent: "enrollment", IsActive: true},
	}

	Categories = []InteractionCategory{
		{ID: 1, Name: "Tertarik", Sentiment: "positive", IsActive: true},
		{ID: 2, Name: "Mempertimbangkan", Sentiment: "neutral", IsActive: true},
		{ID: 3, Name: "Ragu-ragu", Sentiment: "neutral", IsActive: true},
		{ID: 4, Name: "Dingin", Sentiment: "negative", IsActive: true},
		{ID: 5, Name: "Tidak bisa dihubungi", Sentiment: "negative", IsActive: true},
	}

	Obstacles = []Obstacle{
		{ID: 1, Name: "Biaya terlalu mahal", SuggestedResponse: "Jelaskan skema cicilan dan beasiswa yang tersedia", IsActive: true},
		{ID: 2, Name: "Lokasi jauh", SuggestedResponse: "Tawarkan fasilitas asrama dan transportasi", IsActive: true},
		{ID: 3, Name: "Orang tua belum setuju", SuggestedResponse: "Undang orang tua untuk campus visit", IsActive: true},
		{ID: 4, Name: "Waktu belum tepat", SuggestedResponse: "Follow up kembali 1-2 minggu lagi", IsActive: true},
		{ID: 5, Name: "Memilih kampus lain", SuggestedResponse: "Tanyakan alasan dan bandingkan keunggulan kita", IsActive: true},
	}

	Announcements = []Announcement{
		{ID: 1, Title: "Selamat Datang di PMB 2025!", Content: "Selamat bergabung di proses penerimaan mahasiswa baru STMIK Tazkia. Silakan lengkapi dokumen Anda.", TargetStatus: "registered", PublishedAt: time.Now().AddDate(0, 0, -7), IsActive: true},
		{ID: 2, Title: "Batas Waktu Pembayaran", Content: "Batas waktu pembayaran biaya pendaftaran adalah 7 hari setelah registrasi.", TargetStatus: "registered", PublishedAt: time.Now().AddDate(0, 0, -5), IsActive: true},
		{ID: 3, Title: "Jadwal Ospek 2025", Content: "Ospek akan dilaksanakan pada tanggal 1-3 September 2025. Pastikan Anda sudah terdaftar.", TargetStatus: "committed", PublishedAt: time.Now().AddDate(0, 0, -2), IsActive: true},
	}

	// Current logged in user (for mockup)
	CurrentUser = Users[2] // Siti Rahayu (consultant)
)

// Helper functions
func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// GetAdminStats returns dashboard statistics
func GetAdminStats() DashboardStats {
	return DashboardStats{
		TotalCandidates:  "10",
		RegisteredCount:  "3",
		ProspectingCount: "3",
		CommittedCount:   "2",
		EnrolledCount:    "1",
		LostCount:        "1",
		OverdueFollowups: "1",
		TodayFollowups:   "2",
		ConversionRate:   10.0,
		ThisMonthLeads:   "8",
		LastMonthLeads:   "5",
	}
}

// GetConsultantStats returns stats for current consultant
func GetConsultantStats() ConsultantStats {
	return ConsultantStats{
		MyCandidates:      5,
		MyProspecting:     2,
		MyCommitted:       1,
		MyEnrolled:        0,
		MyLost:            1,
		OverdueFollowups:  1,
		TodayFollowups:    2,
		UnreadSuggestions: 1,
	}
}

// CandidateView is a view model with string IDs for templates
type CandidateView struct {
	ID                  string
	Name                string
	Email               string
	Phone               string
	WhatsApp            string
	Address             string
	City                string
	Province            string
	HighSchool          string
	GraduationYear      string
	ProdiName           string
	SourceType          string
	SourceDetail        string
	CampaignName        string
	ReferrerName        string
	Status              string
	ConsultantName      string
	RegistrationFeePaid bool
	CreatedAt           string
	NextFollowup        string
	IsOverdue           bool
}

// InteractionView is a view model for interactions
type InteractionView struct {
	ID                   string
	Channel              string
	Category             string
	CategorySentiment    string
	Obstacle             string
	Remarks              string
	NextFollowupDate     string
	SupervisorSuggestion string
	SuggestionRead       bool
	ConsultantName       string
	CreatedAt            string
}

// GetCandidateByID returns a candidate by ID (string)
func GetCandidateByID(id string) *CandidateView {
	for _, c := range Candidates {
		if fmt.Sprintf("%d", c.ID) == id {
			var nextFollowup string
			var isOverdue bool
			if c.NextFollowupAt != nil {
				nextFollowup = c.NextFollowupAt.Format("02 Jan 2006")
				isOverdue = c.NextFollowupAt.Before(time.Now())
			}
			return &CandidateView{
				ID:                  fmt.Sprintf("%d", c.ID),
				Name:                c.Name,
				Email:               c.Email,
				Phone:               c.Phone,
				WhatsApp:            c.WhatsApp,
				Address:             c.Address,
				City:                c.City,
				Province:            c.Province,
				HighSchool:          c.HighSchool,
				GraduationYear:      fmt.Sprintf("%d", c.GraduationYear),
				ProdiName:           c.ProdiName,
				SourceType:          c.SourceType,
				SourceDetail:        c.SourceDetail,
				CampaignName:        c.CampaignName,
				ReferrerName:        c.ReferrerName,
				Status:              c.Status,
				ConsultantName:      c.ConsultantName,
				RegistrationFeePaid: c.RegistrationFeePaid,
				CreatedAt:           c.CreatedAt.Format("02 Jan 2006"),
				NextFollowup:        nextFollowup,
				IsOverdue:           isOverdue,
			}
		}
	}
	return nil
}

// GetInteractionsByCandidateID returns interactions for a candidate (string ID)
func GetInteractionsByCandidateID(candidateID string) []InteractionView {
	var result []InteractionView
	for _, i := range Interactions {
		if fmt.Sprintf("%d", i.CandidateID) == candidateID {
			var nextFollowup string
			if i.NextFollowupDate != nil {
				nextFollowup = i.NextFollowupDate.Format("02 Jan 2006")
			}
			result = append(result, InteractionView{
				ID:                   fmt.Sprintf("%d", i.ID),
				Channel:              i.Channel,
				Category:             i.Category,
				CategorySentiment:    i.CategorySentiment,
				Obstacle:             i.Obstacle,
				Remarks:              i.Remarks,
				NextFollowupDate:     nextFollowup,
				SupervisorSuggestion: i.SupervisorSuggestion,
				SuggestionRead:       i.SuggestionReadAt != nil,
				ConsultantName:       i.ConsultantName,
				CreatedAt:            i.CreatedAt.Format("02 Jan 2006 15:04"),
			})
		}
	}
	return result
}

// FilterCandidates filters candidates by multiple criteria
func FilterCandidates(status, prodi, consultant, search string) []CandidateView {
	var result []CandidateView
	for _, c := range Candidates {
		// Filter by status
		if status != "" && status != "all" && c.Status != status {
			continue
		}
		// Filter by prodi
		if prodi != "" && fmt.Sprintf("%d", c.ProdiID) != prodi {
			continue
		}
		// Filter by consultant
		if consultant != "" && fmt.Sprintf("%d", c.AssignedConsultantID) != consultant {
			continue
		}
		// Filter by search (name, email, phone)
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(c.Name), searchLower) &&
				!strings.Contains(strings.ToLower(c.Email), searchLower) &&
				!strings.Contains(c.Phone, search) {
				continue
			}
		}

		var nextFollowup string
		var isOverdue bool
		if c.NextFollowupAt != nil {
			nextFollowup = c.NextFollowupAt.Format("02 Jan 2006")
			isOverdue = c.NextFollowupAt.Before(time.Now())
		}

		result = append(result, CandidateView{
			ID:             fmt.Sprintf("%d", c.ID),
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
			NextFollowup:   nextFollowup,
			IsOverdue:      isOverdue,
		})
	}
	return result
}
