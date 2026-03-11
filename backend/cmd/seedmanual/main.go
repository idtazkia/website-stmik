// Command seedmanual seeds comprehensive test data for user manual screenshots.
// Creates realistic dataset: staff, candidates at every status, interactions,
// documents, billings, payments, campaigns, referrers, announcements, etc.
//
// This should only be run in test/development environments.
package main

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/idtazkia/stmik-admission-api/internal/config"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/internal/pkg/crypto"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := crypto.Init(cfg.Encryption.Key); err != nil {
		log.Fatalf("failed to initialize encryption: %v", err)
	}

	ctx := context.Background()
	if err := model.Connect(ctx, cfg.Database.DSN()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer model.Close()

	log.Println("Seeding comprehensive manual data...")

	s := &seeder{ctx: ctx}
	s.seedAll()

	fmt.Println("Manual seed data created successfully")
}

type seeder struct {
	ctx context.Context

	// IDs collected during seeding for cross-references
	adminID      string
	supervisorID string
	consultant1  string
	consultant2  string
	financeID    string
	academicID   string

	prodi1 string
	prodi2 string

	campaign1 string
	campaign2 string
	campaign3 string

	referrer1 string
	referrer2 string

	category1 string // positif
	category2 string // negatif
	category3 string // netral

	obstacle1 string
	obstacle2 string
	obstacle3 string

	lostReason1 string
	lostReason2 string

	docType1 string
	docType2 string
	docType3 string
	docType4 string

	feeType1 string // from migration seed
	feeType2 string
}

func (s *seeder) seedAll() {
	s.seedUsers()
	s.seedProdi()
	s.seedInteractionCategories()
	s.seedObstacles()
	s.seedLostReasons()
	s.seedDocumentTypes()
	s.seedCampaigns()
	s.seedReferrers()
	s.seedRewardConfigs()
	s.seedFeeStructures()
	s.seedCandidates()
	s.seedAnnouncements()
}

// --- Staff Users ---

func (s *seeder) seedUsers() {
	log.Println("  Seeding staff users...")

	users := []struct {
		email, name, role string
		idField           *string
	}{
		{"admin@tazkia.ac.id", "Dr. Ahmad Fauzan, M.Kom", "admin", &s.adminID},
		{"supervisor@tazkia.ac.id", "Dra. Siti Nurhaliza, M.M.", "supervisor", &s.supervisorID},
		{"konsultan1@tazkia.ac.id", "Rizky Pratama, S.Kom", "consultant", &s.consultant1},
		{"konsultan2@tazkia.ac.id", "Dewi Anggraini, S.E.", "consultant", &s.consultant2},
		{"finance@tazkia.ac.id", "Hendra Wijaya, S.Ak", "finance", &s.financeID},
		{"akademik@tazkia.ac.id", "Dr. Budi Santoso, M.T.", "academic", &s.academicID},
	}

	for _, u := range users {
		existing, _ := model.FindUserByEmail(s.ctx, u.email)
		if existing != nil {
			*u.idField = existing.ID
			log.Printf("    [skip] %s", u.email)
			continue
		}
		user, err := model.CreateUser(s.ctx, u.email, u.name, "", u.role)
		if err != nil {
			log.Fatalf("    failed to create user %s: %v", u.email, err)
		}
		*u.idField = user.ID
		log.Printf("    [created] %s (%s)", u.name, u.role)
	}

	// Assign supervisor to both consultants
	model.UpdateUserSupervisor(s.ctx, s.consultant1, &s.supervisorID)
	model.UpdateUserSupervisor(s.ctx, s.consultant2, &s.supervisorID)
}

// --- Program Studi ---

func (s *seeder) seedProdi() {
	log.Println("  Seeding program studi...")

	prodis := []struct {
		name, code, degree string
		idField            *string
	}{
		{"Sistem Informasi", "SI", "S1", &s.prodi1},
		{"Teknik Informatika", "TI", "S1", &s.prodi2},
	}

	for _, p := range prodis {
		prodi, err := model.CreateProdi(s.ctx, p.name, p.code, p.degree)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", p.name)
			// Try to find existing
			list, _ := model.ListProdis(s.ctx, false)
			for _, existing := range list {
				if existing.Code == p.code {
					*p.idField = existing.ID
				}
			}
			continue
		}
		*p.idField = prodi.ID
		log.Printf("    [created] %s (%s)", p.name, p.code)
	}
}

// --- Interaction Categories ---

func (s *seeder) seedInteractionCategories() {
	log.Println("  Seeding interaction categories...")

	cats := []struct {
		name, sentiment string
		order           int
		idField         *string
	}{
		{"Tertarik & Antusias", "positive", 1, &s.category1},
		{"Ragu / Keberatan", "negative", 2, &s.category2},
		{"Belum Bisa Dihubungi", "neutral", 3, &s.category3},
	}

	for _, c := range cats {
		cat, err := model.CreateInteractionCategory(s.ctx, c.name, c.sentiment, c.order)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", c.name)
			list, _ := model.ListInteractionCategories(s.ctx, false)
			for _, existing := range list {
				if existing.Name == c.name {
					*c.idField = existing.ID
				}
			}
			continue
		}
		*c.idField = cat.ID
		log.Printf("    [created] %s", c.name)
	}
}

// --- Obstacles ---

func (s *seeder) seedObstacles() {
	log.Println("  Seeding obstacles...")

	resp1 := "Jelaskan program beasiswa dan cicilan yang tersedia"
	resp2 := "Tawarkan program kuliah kelas karyawan"
	resp3 := "Jelaskan keunggulan dan akreditasi program studi"

	obstacles := []struct {
		name    string
		resp    *string
		order   int
		idField *string
	}{
		{"Kendala Biaya", &resp1, 1, &s.obstacle1},
		{"Bentrok Jadwal / Kerja", &resp2, 2, &s.obstacle2},
		{"Mempertimbangkan Kampus Lain", &resp3, 3, &s.obstacle3},
	}

	for _, o := range obstacles {
		obs, err := model.CreateObstacle(s.ctx, o.name, o.resp, o.order)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", o.name)
			list, _ := model.ListObstacles(s.ctx, false)
			for _, existing := range list {
				if existing.Name == o.name {
					*o.idField = existing.ID
				}
			}
			continue
		}
		*o.idField = obs.ID
		log.Printf("    [created] %s", o.name)
	}
}

// --- Lost Reasons ---

func (s *seeder) seedLostReasons() {
	log.Println("  Seeding lost reasons...")

	desc1 := "Kandidat memilih mendaftar ke institusi lain"
	desc2 := "Kandidat tidak merespon setelah beberapa kali dihubungi"

	reasons := []struct {
		name    string
		desc    *string
		order   int
		idField *string
	}{
		{"Memilih Kampus Lain", &desc1, 1, &s.lostReason1},
		{"Tidak Merespon", &desc2, 2, &s.lostReason2},
	}

	for _, r := range reasons {
		lr, err := model.CreateLostReason(s.ctx, r.name, r.desc, r.order)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", r.name)
			list, _ := model.ListLostReasons(s.ctx, false)
			for _, existing := range list {
				if existing.Name == r.name {
					*r.idField = existing.ID
				}
			}
			continue
		}
		*r.idField = lr.ID
		log.Printf("    [created] %s", r.name)
	}
}

// --- Document Types ---

func (s *seeder) seedDocumentTypes() {
	log.Println("  Seeding document types...")

	desc1 := "Scan KTP atau Kartu Pelajar"
	desc2 := "Scan ijazah SMA/SMK/MA"
	desc3 := "Scan transkrip nilai atau rapor semester terakhir"
	desc4 := "Foto 3x4 latar belakang merah (format JPG/PNG)"

	types := []struct {
		name, code string
		desc       *string
		required   bool
		order      int
		idField    *string
	}{
		{"KTP / Kartu Identitas", "KTP", &desc1, true, 1, &s.docType1},
		{"Ijazah", "IJAZAH", &desc2, true, 2, &s.docType2},
		{"Transkrip Nilai", "TRANSKRIP", &desc3, true, 3, &s.docType3},
		{"Pas Foto 3x4", "FOTO", &desc4, true, 4, &s.docType4},
	}

	for _, dt := range types {
		docType, err := model.CreateDocumentType(s.ctx, dt.name, dt.code, dt.desc, dt.required, false, 5, dt.order)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", dt.name)
			list, _ := model.ListDocumentTypes(s.ctx, false)
			for _, existing := range list {
				if existing.Code == dt.code {
					*dt.idField = existing.ID
				}
			}
			continue
		}
		*dt.idField = docType.ID
		log.Printf("    [created] %s", dt.name)
	}
}

// --- Campaigns ---

func (s *seeder) seedCampaigns() {
	log.Println("  Seeding campaigns...")

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.Local)
	ch1 := "Instagram"
	ch2 := "Google Ads"
	ch3 := "Offline"
	d1 := "Kampanye Instagram Ads untuk PMB 2026"
	d2 := "Google Search Ads untuk keyword pendaftaran kuliah"
	d3 := "Expo pendidikan di Jakarta Convention Center"

	campaigns := []struct {
		name, typ string
		channel   *string
		desc      *string
		idField   *string
	}{
		{"PMB 2026 - Instagram", "digital", &ch1, &d1, &s.campaign1},
		{"PMB 2026 - Google Ads", "digital", &ch2, &d2, &s.campaign2},
		{"Expo Pendidikan JCC 2026", "event", &ch3, &d3, &s.campaign3},
	}

	for _, c := range campaigns {
		camp, err := model.CreateCampaign(s.ctx, c.name, c.typ, c.channel, c.desc, &start, &end, nil)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", c.name)
			list, _ := model.ListCampaigns(s.ctx, true)
			for _, existing := range list {
				if existing.Name == c.name {
					*c.idField = existing.ID
				}
			}
			continue
		}
		*c.idField = camp.ID
		log.Printf("    [created] %s", c.name)
	}
}

// --- Referrers ---

func (s *seeder) seedReferrers() {
	log.Println("  Seeding referrers...")

	inst1 := "STMIK Tazkia"
	ph1 := "081234567890"
	em1 := "mahasiswa1@student.tazkia.ac.id"
	code1 := "REF-MHS001"
	bank1 := "BCA"
	acc1 := "1234567890"
	holder1 := "Andi Saputra"

	inst2 := "SMA Negeri 1 Bogor"
	ph2 := "081298765432"
	em2 := "bk@sman1bogor.sch.id"
	code2 := "REF-GURU01"
	bank2 := "BRI"
	acc2 := "0987654321"
	holder2 := "Ibu Ratna"

	referrers := []struct {
		name, typ string
		inst      *string
		phone     *string
		email     *string
		code      *string
		bank      *string
		acc       *string
		holder    *string
		idField   *string
	}{
		{"Andi Saputra", "student", &inst1, &ph1, &em1, &code1, &bank1, &acc1, &holder1, &s.referrer1},
		{"Ibu Ratna, S.Pd", "teacher", &inst2, &ph2, &em2, &code2, &bank2, &acc2, &holder2, &s.referrer2},
	}

	for _, r := range referrers {
		ref, err := model.CreateReferrer(s.ctx, r.name, r.typ, r.inst, r.phone, r.email, r.code, r.bank, r.acc, r.holder, nil, "transfer")
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", r.name)
			list, _ := model.ListReferrers(s.ctx, "")
			for _, existing := range list {
				if existing.Name == r.name {
					*r.idField = existing.ID
				}
			}
			continue
		}
		*r.idField = ref.ID
		log.Printf("    [created] %s (%s)", r.name, r.typ)
	}
}

// --- Reward Configs ---

func (s *seeder) seedRewardConfigs() {
	log.Println("  Seeding reward configs...")

	d1 := "Komisi untuk mahasiswa yang mereferensikan teman"
	d2 := "Komisi untuk guru BK yang mereferensikan siswa"

	configs := []struct {
		referrerType, rewardType, trigger string
		amount                            int64
		desc                              *string
	}{
		{"student", "cash", "enrolled", 500000, &d1},
		{"teacher", "cash", "enrolled", 1000000, &d2},
	}

	for _, c := range configs {
		_, err := model.CreateRewardConfig(s.ctx, c.referrerType, c.rewardType, c.amount, false, c.trigger, c.desc)
		if err != nil {
			log.Printf("    [skip] reward %s/%s (likely exists)", c.referrerType, c.trigger)
			continue
		}
		log.Printf("    [created] reward %s: Rp %d", c.referrerType, c.amount)
	}
}

// --- Fee Structures ---

func (s *seeder) seedFeeStructures() {
	log.Println("  Seeding fee structures...")

	// Get fee types from migration seed
	feeTypes, _ := model.ListFeeTypes(s.ctx)
	for _, ft := range feeTypes {
		if ft.Code == "REGISTRATION" || ft.Code == "registration" {
			s.feeType1 = ft.ID
		}
		if ft.Code == "TUITION" || ft.Code == "tuition" {
			s.feeType2 = ft.ID
		}
	}

	if s.feeType1 == "" || s.feeType2 == "" {
		log.Println("    [skip] fee types not found in migration seed, skipping fee structures")
		return
	}

	structures := []struct {
		feeTypeID string
		prodiID   *string
		year      string
		amount    int64
	}{
		{s.feeType1, &s.prodi1, "2026/2027", 500000},
		{s.feeType1, &s.prodi2, "2026/2027", 500000},
		{s.feeType2, &s.prodi1, "2026/2027", 7500000},
		{s.feeType2, &s.prodi2, "2026/2027", 8000000},
	}

	for _, fs := range structures {
		_, err := model.CreateFeeStructure(s.ctx, fs.feeTypeID, fs.prodiID, fs.year, fs.amount)
		if err != nil {
			log.Printf("    [skip] fee structure (likely exists)")
			continue
		}
		log.Printf("    [created] fee structure Rp %d", fs.amount)
	}
}

// --- Candidates (the bulk of the data) ---

func (s *seeder) seedCandidates() {
	log.Println("  Seeding candidates...")

	now := time.Now()

	// Candidate dataset — covers every status in the funnel
	candidates := []struct {
		email, phone, name     string
		address, city, prov    string
		highSchool             string
		gradYear               int
		prodiID                string
		consultantID           string
		campaignID             string
		referrerID             string
		status                 string // target status
		numInteractions        int
		hasDocuments           bool
		hasBilling             bool
		hasPayment             bool
		lostReasonID           string
	}{
		// === ENROLLED (3) — complete journey ===
		{
			"farhan.akbar@gmail.com", "081311111111", "Muhammad Farhan Akbar",
			"Jl. Raya Bogor No. 123", "Bogor", "Jawa Barat",
			"SMA Negeri 1 Bogor", 2025, "", "", "", "",
			"enrolled", 5, true, true, true, "",
		},
		{
			"nadia.safitri@gmail.com", "081322222222", "Nadia Putri Safitri",
			"Jl. Sudirman No. 45", "Jakarta Selatan", "DKI Jakarta",
			"SMA Labschool Kebayoran", 2025, "", "", "", "",
			"enrolled", 4, true, true, true, "",
		},
		{
			"reza.firmansyah@gmail.com", "081333333333", "Reza Firmansyah",
			"Jl. Margonda Raya No. 78", "Depok", "Jawa Barat",
			"SMK Telkom Depok", 2025, "", "", "", "",
			"enrolled", 6, true, true, true, "",
		},

		// === COMMITTED (4) — awaiting enrollment ===
		{
			"sari.wulandari@gmail.com", "081344444444", "Sari Wulandari",
			"Jl. Pemuda No. 10", "Bekasi", "Jawa Barat",
			"SMA Negeri 3 Bekasi", 2025, "", "", "", "",
			"committed", 3, true, true, false, "",
		},
		{
			"dimas.prasetyo@gmail.com", "081355555555", "Dimas Prasetyo",
			"Jl. Veteran No. 22", "Tangerang", "Banten",
			"SMA Al-Azhar BSD", 2025, "", "", "", "",
			"committed", 4, true, false, false, "",
		},
		{
			"anisa.rahma@gmail.com", "081366666666", "Anisa Rahmawati",
			"Jl. Pahlawan No. 5", "Bogor", "Jawa Barat",
			"MA Negeri 1 Bogor", 2025, "", "", "", "",
			"committed", 2, true, false, false, "",
		},
		{
			"fajar.nugroho@gmail.com", "081377777777", "Fajar Nugroho",
			"Jl. Merdeka No. 99", "Depok", "Jawa Barat",
			"SMA Negeri 2 Depok", 2026, "", "", "", "",
			"committed", 3, false, false, false, "",
		},

		// === PROSPECTING (6) — in communication ===
		{
			"lia.kusuma@gmail.com", "081388888888", "Lia Kusumawati",
			"Jl. Gatot Subroto No. 15", "Jakarta Selatan", "DKI Jakarta",
			"SMA Negeri 8 Jakarta", 2026, "", "", "", "",
			"prospecting", 2, false, false, false, "",
		},
		{
			"aldi.ramadhan@gmail.com", "081399999999", "Aldi Ramadhan",
			"Jl. Asia Afrika No. 30", "Bandung", "Jawa Barat",
			"SMA Negeri 3 Bandung", 2026, "", "", "", "",
			"prospecting", 1, false, false, false, "",
		},
		{
			"putri.amelia@gmail.com", "081400000001", "Putri Amelia Sari",
			"Jl. Diponegoro No. 7", "Bogor", "Jawa Barat",
			"SMK Wikrama Bogor", 2026, "", "", "", "",
			"prospecting", 3, true, false, false, "",
		},
		{
			"bima.aditya@gmail.com", "081400000002", "Bima Aditya Putra",
			"Jl. Ahmad Yani No. 42", "Bekasi", "Jawa Barat",
			"SMA Negeri 1 Bekasi", 2026, "", "", "", "",
			"prospecting", 1, false, false, false, "",
		},
		{
			"intan.permata@gmail.com", "081400000003", "Intan Permatasari",
			"Jl. Thamrin No. 18", "Jakarta Pusat", "DKI Jakarta",
			"SMA Gonzaga Jakarta", 2026, "", "", "", "",
			"prospecting", 2, false, false, false, "",
		},
		{
			"yoga.pratama@gmail.com", "081400000004", "Yoga Eka Pratama",
			"Jl. Juanda No. 55", "Bogor", "Jawa Barat",
			"SMA Negeri 5 Bogor", 2026, "", "", "", "",
			"prospecting", 1, false, false, false, "",
		},

		// === REGISTERED (5) — just signed up, no interaction yet ===
		{
			"maya.anggraini@gmail.com", "081400000005", "Maya Anggraini",
			"Jl. Cikini Raya No. 9", "Jakarta Pusat", "DKI Jakarta",
			"SMA Negeri 4 Jakarta", 2026, "", "", "", "",
			"registered", 0, false, false, false, "",
		},
		{
			"hadi.wijaya@gmail.com", "081400000006", "Hadi Wijaya",
			"Jl. Surya Kencana No. 33", "Bogor", "Jawa Barat",
			"SMK Negeri 1 Bogor", 2026, "", "", "", "",
			"registered", 0, false, false, false, "",
		},
		{
			"ratna.dewi@gmail.com", "081400000007", "Ratna Dewi Lestari",
			"Jl. Raden Saleh No. 12", "Cirebon", "Jawa Barat",
			"SMA Negeri 1 Cirebon", 2026, "", "", "", "",
			"registered", 0, false, false, false, "",
		},
		{
			"fikri.hakim@gmail.com", "081400000008", "Fikri Hakim",
			"Jl. Ciputat Raya No. 20", "Tangerang Selatan", "Banten",
			"SMA Muhammadiyah 1 Ciputat", 2026, "", "", "", "",
			"registered", 0, false, false, false, "",
		},
		{
			"lina.marlina@gmail.com", "081400000009", "Lina Marlina",
			"Jl. Pajajaran No. 65", "Bogor", "Jawa Barat",
			"MA Amanatul Ummah", 2026, "", "", "", "",
			"registered", 0, false, false, false, "",
		},

		// === LOST (3) — dropped out of funnel ===
		{
			"taufik.hidayat@gmail.com", "081400000010", "Taufik Hidayat",
			"Jl. Kebon Jeruk No. 8", "Jakarta Barat", "DKI Jakarta",
			"SMA Negeri 70 Jakarta", 2025, "", "", "", "",
			"lost", 2, false, false, false, "",
		},
		{
			"sarah.azzahra@gmail.com", "081400000011", "Sarah Azzahra",
			"Jl. Pangeran Antasari No. 14", "Jakarta Selatan", "DKI Jakarta",
			"SMA Al-Izhar Pondok Labu", 2025, "", "", "", "",
			"lost", 3, false, false, false, "",
		},
		{
			"doni.setiawan@gmail.com", "081400000012", "Doni Setiawan",
			"Jl. Raya Serpong No. 40", "Tangerang Selatan", "Banten",
			"SMA Negeri 1 Serpong", 2025, "", "", "", "",
			"lost", 1, false, false, false, "",
		},
	}

	for i, c := range candidates {
		// Assign prodi alternating
		prodiID := s.prodi1
		if i%2 == 1 {
			prodiID = s.prodi2
		}
		candidates[i].prodiID = prodiID

		// Assign consultant alternating
		consultantID := s.consultant1
		if i%2 == 1 {
			consultantID = s.consultant2
		}
		candidates[i].consultantID = consultantID

		// Assign campaign round-robin
		switch i % 3 {
		case 0:
			candidates[i].campaignID = s.campaign1
		case 1:
			candidates[i].campaignID = s.campaign2
		case 2:
			candidates[i].campaignID = s.campaign3
		}

		// Some candidates have referrers
		if i%4 == 0 {
			candidates[i].referrerID = s.referrer1
		} else if i%4 == 1 {
			candidates[i].referrerID = s.referrer2
		}

		// Lost reasons
		if c.status == "lost" {
			if i%2 == 0 {
				candidates[i].lostReasonID = s.lostReason1
			} else {
				candidates[i].lostReasonID = s.lostReason2
			}
		}
	}

	for _, c := range candidates {
		s.seedOneCandidate(c.email, c.phone, c.name, c.address, c.city, c.prov,
			c.highSchool, c.gradYear, c.prodiID, c.consultantID,
			c.campaignID, c.referrerID, c.status,
			c.numInteractions, c.hasDocuments, c.hasBilling, c.hasPayment,
			c.lostReasonID, now)
	}
}

func (s *seeder) seedOneCandidate(
	email, phone, name, address, city, province string,
	highSchool string, gradYear int, prodiID, consultantID string,
	campaignID, referrerID, targetStatus string,
	numInteractions int, hasDocuments, hasBilling, hasPayment bool,
	lostReasonID string, now time.Time,
) {
	ctx := s.ctx

	// Check if already exists
	existing, _ := model.FindCandidateByEmail(ctx, email)
	if existing != nil {
		log.Printf("    [skip] %s (exists)", name)
		return
	}

	// Create candidate
	hash, _ := model.HashPassword("password123")
	candidate, err := model.CreateCandidate(ctx, email, phone, hash)
	if err != nil {
		log.Printf("    [error] %s: %v", name, err)
		return
	}

	// Update personal info
	model.UpdateCandidatePersonalInfo(ctx, candidate.ID, name, address, city, province)

	// Update education
	model.UpdateCandidateEducation(ctx, candidate.ID, highSchool, gradYear, prodiID)

	// Source tracking
	sourceType := "website"
	if campaignID != "" {
		sourceType = "campaign"
	}
	if referrerID != "" {
		sourceType = "referral"
	}
	model.UpdateCandidateSourceTracking(ctx, candidate.ID, sourceType, "", campaignID, referrerID, "")

	// Assign consultant
	model.AssignCandidateConsultant(ctx, candidate.ID, consultantID)

	// Create interactions
	categories := []string{s.category1, s.category2, s.category3}
	channels := []string{"whatsapp", "phone", "email", "campus_visit"}

	for j := 0; j < numInteractions; j++ {
		catID := categories[j%len(categories)]
		channel := channels[j%len(channels)]
		followUp := now.AddDate(0, 0, j+3)
		var obstacleID *string
		if j > 0 && catID == s.category2 {
			obs := []string{s.obstacle1, s.obstacle2, s.obstacle3}
			obstacleID = &obs[j%len(obs)]
		}

		remarks := fmt.Sprintf("Interaksi ke-%d dengan %s via %s", j+1, name, channel)
		model.CreateInteraction(ctx, candidate.ID, consultantID, channel, &catID, obstacleID, remarks, &followUp, nil)
	}

	// Upload documents
	if hasDocuments {
		docTypes := []string{s.docType1, s.docType2, s.docType3, s.docType4}
		for k, dtID := range docTypes {
			if dtID == "" {
				continue
			}
			fileName := fmt.Sprintf("doc_%s_%d.pdf", candidate.ID[:8], k+1)
			filePath := fmt.Sprintf("uploads/%s/%s", candidate.ID[:8], fileName)
			doc, err := model.CreateDocument(ctx, candidate.ID, dtID, fileName, filePath, 1024000, "application/pdf")
			if err != nil {
				continue
			}
			// Approve first 2 documents, leave rest pending
			if k < 2 {
				model.ApproveDocument(ctx, doc.ID, s.adminID)
			}
		}
	}

	// Create billing
	if hasBilling {
		dueDate := now.AddDate(0, 1, 0)
		billing, err := model.CreateBilling(ctx, candidate.ID, "registration", nil, 500000, &dueDate)
		if err == nil && hasPayment {
			// Create payment
			payment, err := model.CreatePayment(ctx, billing.ID, 500000, now, "uploads/proof.jpg", "proof.jpg", 512000, "image/jpeg")
			if err == nil {
				model.ApprovePayment(ctx, payment.ID, s.financeID)
			}
		}
	}

	// Set final status
	switch targetStatus {
	case "prospecting":
		// Already set by first interaction creation
	case "committed":
		model.UpdateCandidateStatus(ctx, candidate.ID, "committed")
	case "enrolled":
		model.UpdateCandidateStatus(ctx, candidate.ID, "committed")
		model.UpdateCandidateStatus(ctx, candidate.ID, "enrolled")
	case "lost":
		if numInteractions > 0 {
			// prospecting -> lost
			model.MarkCandidateLost(ctx, candidate.ID, lostReasonID)
		} else {
			model.UpdateCandidateStatus(ctx, candidate.ID, "lost")
		}
	}

	log.Printf("    [created] %s → %s (%d interactions)", name, targetStatus, numInteractions)
}

// --- Announcements ---

func (s *seeder) seedAnnouncements() {
	log.Println("  Seeding announcements...")

	announcements := []struct {
		title, content string
		target         *string
		publish        bool
	}{
		{
			"Jadwal Ujian Masuk PMB 2026/2027",
			"Ujian masuk gelombang pertama akan dilaksanakan pada tanggal 15 April 2026. " +
				"Peserta yang sudah melengkapi dokumen dan pembayaran dapat mengikuti ujian di kampus STMIK Tazkia. " +
				"Silakan bawa KTP asli dan kartu peserta ujian.\n\n" +
				"Waktu: 08.00 - 12.00 WIB\nTempat: Gedung A Lt. 3, STMIK Tazkia",
			nil, true,
		},
		{
			"Info Beasiswa Prestasi 2026",
			"STMIK Tazkia membuka program beasiswa prestasi untuk calon mahasiswa baru dengan ketentuan:\n\n" +
				"- Beasiswa 100%: Nilai rata-rata rapor >= 90\n" +
				"- Beasiswa 50%: Nilai rata-rata rapor >= 85\n" +
				"- Beasiswa 25%: Nilai rata-rata rapor >= 80\n\n" +
				"Dokumen tambahan: Scan rapor semester 1-5 dan surat rekomendasi dari kepala sekolah.",
			nil, true,
		},
		{
			"Pengumuman Teknis Registrasi Online",
			"Bagi calon mahasiswa yang mengalami kendala saat registrasi online, silakan hubungi:\n\n" +
				"- WhatsApp: 0812-3456-7890\n- Email: admisi@tazkia.ac.id\n\n" +
				"Jam layanan: Senin-Jumat, 08.00-16.00 WIB",
			nil, true,
		},
		{
			"[Draft] Jadwal Ospek 2026",
			"Orientasi Studi dan Pengenalan Kampus (Ospek) akan dilaksanakan pada minggu pertama September 2026. " +
				"Detail jadwal dan ketentuan akan diumumkan setelah masa registrasi ditutup.",
			nil, false, // draft
		},
	}

	for _, a := range announcements {
		ann, err := model.CreateAnnouncement(s.ctx, a.title, a.content, a.target, nil, &s.adminID)
		if err != nil {
			log.Printf("    [skip] %s (likely exists)", a.title)
			continue
		}
		if a.publish {
			model.PublishAnnouncement(s.ctx, ann.ID)
		}
		log.Printf("    [created] %s (published=%v)", a.title, a.publish)
	}
}
