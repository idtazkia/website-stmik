package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/idtazkia/stmik-admission-api/internal/auth"
	"github.com/idtazkia/stmik-admission-api/internal/integration"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/portal"
)

// PublicHandler handles public routes for candidate registration and login
type PublicHandler struct {
	sessionMgr *auth.SessionManager
	resend     *integration.ResendClient
	whatsapp   *integration.WhatsAppClient
}

// NewPublicHandler creates a new public handler
func NewPublicHandler(sessionMgr *auth.SessionManager, resend *integration.ResendClient, whatsapp *integration.WhatsAppClient) *PublicHandler {
	return &PublicHandler{
		sessionMgr: sessionMgr,
		resend:     resend,
		whatsapp:   whatsapp,
	}
}

// RegisterRoutes registers all public routes to the mux
func (h *PublicHandler) RegisterRoutes(mux *http.ServeMux) {
	// Registration routes
	mux.HandleFunc("GET /register", h.handleRegister)
	mux.HandleFunc("POST /register/step1", h.handleRegisterStep1)
	mux.HandleFunc("POST /register/step2", h.handleRegisterStep2)
	mux.HandleFunc("POST /register/step3", h.handleRegisterStep3)
	mux.HandleFunc("POST /register/step4", h.handleRegisterStep4)

	// Login routes
	mux.HandleFunc("GET /login", h.handleLogin)
	mux.HandleFunc("POST /login", h.handleLoginSubmit)
	mux.HandleFunc("POST /logout", h.handleLogout)

	// Optional verification routes (can be used from portal later)
	mux.HandleFunc("POST /portal/verify-email", h.handleRequestEmailOTP)
	mux.HandleFunc("POST /portal/confirm-email", h.handleConfirmEmailOTP)
	mux.HandleFunc("POST /portal/verify-phone", h.handleRequestPhoneOTP)
	mux.HandleFunc("POST /portal/confirm-phone", h.handleConfirmPhoneOTP)

	// HTMX-compatible verification routes
	mux.HandleFunc("POST /portal/verify-email/send", h.handleSendEmailOTP)
	mux.HandleFunc("POST /portal/verify-email/confirm", h.handleVerifyEmailOTP)
}

// Source type options
var sourceTypes = []portal.SourceTypeOption{
	{Value: "instagram", Label: "Instagram"},
	{Value: "google", Label: "Google"},
	{Value: "tiktok", Label: "TikTok"},
	{Value: "youtube", Label: "YouTube"},
	{Value: "expo", Label: "Pameran Pendidikan"},
	{Value: "school_visit", Label: "Kunjungan Sekolah"},
	{Value: "friend_family", Label: "Teman/Keluarga"},
	{Value: "teacher_alumni", Label: "Guru/Alumni"},
	{Value: "walkin", Label: "Datang Langsung"},
	{Value: "other", Label: "Lainnya"},
}

// handleRegister shows the registration form
func (h *PublicHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Get step from query param
	step := r.URL.Query().Get("step")
	if step == "" {
		step = "account"
	}

	// Check if already logged in as candidate
	claims, _ := h.sessionMgr.GetClaimsFromRequest(r)
	if claims != nil && claims.IsCandidate {
		// If on step 1 (account), redirect to portal since account already created
		// For other steps, allow continuing registration
		if step == "account" {
			// Check if candidate has completed registration (has name set)
			candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
			if err == nil && candidate != nil && candidate.Name != nil && *candidate.Name != "" {
				http.Redirect(w, r, "/portal", http.StatusFound)
				return
			}
			// Candidate hasn't completed registration, redirect to step 2
			http.Redirect(w, r, "/register?step=personal", http.StatusFound)
			return
		}
	} else if step != "account" {
		// Not logged in but trying to access step 2+, redirect to step 1
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Get ref and campaign from URL
	refCode := r.URL.Query().Get("ref")
	campaignCode := r.URL.Query().Get("utm_campaign")

	regData := portal.RegistrationData{
		RefCode:      refCode,
		CampaignCode: campaignCode,
	}

	h.renderRegistration(w, r, step, regData, "")
}

// handleRegisterStep1 handles account creation (email/phone + password)
func (h *PublicHandler) handleRegisterStep1(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderRegistration(w, r, "account", portal.RegistrationData{}, "Gagal memproses form")
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")
	refCode := r.FormValue("ref")
	campaignCode := r.FormValue("utm_campaign")

	regData := portal.RegistrationData{
		Email:        email,
		Phone:        phone,
		RefCode:      refCode,
		CampaignCode: campaignCode,
	}

	// Validate at least one of email or phone is provided
	if email == "" && phone == "" {
		h.renderRegistration(w, r, "account", regData, "Harap isi email atau nomor HP")
		return
	}

	// Validate password
	if len(password) < 8 {
		h.renderRegistration(w, r, "account", regData, "Password minimal 8 karakter")
		return
	}

	if password != passwordConfirm {
		h.renderRegistration(w, r, "account", regData, "Konfirmasi password tidak cocok")
		return
	}

	// Check if email already exists
	if email != "" {
		existing, err := model.FindCandidateByEmail(r.Context(), email)
		if err != nil {
			slog.Error("failed to check existing email", "error", err)
			h.renderRegistration(w, r, "account", regData, "Terjadi kesalahan sistem")
			return
		}
		if existing != nil {
			h.renderRegistration(w, r, "account", regData, "Email sudah terdaftar")
			return
		}
	}

	// Check if phone already exists
	if phone != "" {
		existing, err := model.FindCandidateByPhone(r.Context(), phone)
		if err != nil {
			slog.Error("failed to check existing phone", "error", err)
			h.renderRegistration(w, r, "account", regData, "Terjadi kesalahan sistem")
			return
		}
		if existing != nil {
			h.renderRegistration(w, r, "account", regData, "Nomor HP sudah terdaftar")
			return
		}
	}

	// Hash password
	passwordHash, err := model.HashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		h.renderRegistration(w, r, "account", regData, "Terjadi kesalahan sistem")
		return
	}

	// Create candidate
	candidate, err := model.CreateCandidate(r.Context(), email, phone, passwordHash)
	if err != nil {
		slog.Error("failed to create candidate", "error", err)
		h.renderRegistration(w, r, "account", regData, "Gagal membuat akun")
		return
	}

	slog.Info("candidate account created", "id", candidate.ID, "email", email, "phone", phone)

	// Create session
	candidateEmail := ""
	if candidate.Email != nil {
		candidateEmail = *candidate.Email
	}
	candidateName := ""
	if candidate.Name != nil {
		candidateName = *candidate.Name
	}

	token, err := h.sessionMgr.CreateCandidateToken(candidate.ID, candidateEmail, candidateName)
	if err != nil {
		slog.Error("failed to create session token", "error", err)
		h.renderRegistration(w, r, "account", regData, "Gagal membuat sesi")
		return
	}

	h.sessionMgr.SetCookie(w, token)

	// Redirect to step 2
	http.Redirect(w, r, "/register?step=personal", http.StatusFound)
}

// handleRegisterStep2 handles personal info submission
func (h *PublicHandler) handleRegisterStep2(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderRegistration(w, r, "personal", portal.RegistrationData{}, "Gagal memproses form")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	city := strings.TrimSpace(r.FormValue("city"))
	province := strings.TrimSpace(r.FormValue("province"))

	regData := portal.RegistrationData{
		Name:     name,
		Address:  address,
		City:     city,
		Province: province,
	}

	// Validate required fields
	if name == "" || address == "" || city == "" || province == "" {
		h.renderRegistration(w, r, "personal", regData, "Harap lengkapi semua field")
		return
	}

	// Update candidate
	if err := model.UpdateCandidatePersonalInfo(r.Context(), claims.CandidateID, name, address, city, province); err != nil {
		slog.Error("failed to update personal info", "error", err)
		h.renderRegistration(w, r, "personal", regData, "Gagal menyimpan data")
		return
	}

	slog.Info("candidate personal info updated", "id", claims.CandidateID)

	// Redirect to step 3
	http.Redirect(w, r, "/register?step=education", http.StatusFound)
}

// handleRegisterStep3 handles education info submission
func (h *PublicHandler) handleRegisterStep3(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderRegistration(w, r, "education", portal.RegistrationData{}, "Gagal memproses form")
		return
	}

	highSchool := strings.TrimSpace(r.FormValue("high_school"))
	graduationYearStr := r.FormValue("graduation_year")
	prodiID := r.FormValue("prodi_id")

	graduationYear, _ := strconv.Atoi(graduationYearStr)

	regData := portal.RegistrationData{
		HighSchool:     highSchool,
		GraduationYear: graduationYear,
		ProdiID:        prodiID,
	}

	// Validate required fields
	if highSchool == "" || graduationYear == 0 || prodiID == "" {
		h.renderRegistration(w, r, "education", regData, "Harap lengkapi semua field")
		return
	}

	// Update candidate
	if err := model.UpdateCandidateEducation(r.Context(), claims.CandidateID, highSchool, graduationYear, prodiID); err != nil {
		slog.Error("failed to update education", "error", err)
		h.renderRegistration(w, r, "education", regData, "Gagal menyimpan data")
		return
	}

	slog.Info("candidate education updated", "id", claims.CandidateID)

	// Redirect to step 4
	http.Redirect(w, r, "/register?step=source", http.StatusFound)
}

// handleRegisterStep4 handles source tracking and completes registration
func (h *PublicHandler) handleRegisterStep4(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderRegistration(w, r, "source", portal.RegistrationData{}, "Gagal memproses form")
		return
	}

	sourceType := r.FormValue("source_type")
	sourceDetail := strings.TrimSpace(r.FormValue("source_detail"))

	regData := portal.RegistrationData{
		SourceType:   sourceType,
		SourceDetail: sourceDetail,
	}

	// Validate required fields
	if sourceType == "" {
		h.renderRegistration(w, r, "source", regData, "Harap pilih sumber informasi")
		return
	}

	// Get referrer ID from ref code if provided
	referrerID := ""
	referredByCandidateID := ""

	// Check if there's a referral code in session or URL
	// For simplicity, we'll use URL params stored in step 1

	// Get campaign ID from campaign code if provided
	campaignID := ""

	// Update source tracking
	if err := model.UpdateCandidateSourceTracking(r.Context(), claims.CandidateID, sourceType, sourceDetail, campaignID, referrerID, referredByCandidateID); err != nil {
		slog.Error("failed to update source tracking", "error", err)
		h.renderRegistration(w, r, "source", regData, "Gagal menyimpan data")
		return
	}

	// Auto-assign consultant
	consultantID, err := model.GetNextConsultantForAssignment(r.Context())
	if err != nil {
		slog.Warn("failed to get next consultant", "error", err)
		// Non-fatal, continue without assignment
	} else if consultantID != nil {
		if err := model.AssignCandidateConsultant(r.Context(), claims.CandidateID, *consultantID); err != nil {
			slog.Warn("failed to assign consultant", "error", err)
		} else {
			slog.Info("consultant assigned", "candidate_id", claims.CandidateID, "consultant_id", *consultantID)
		}
	}

	// Update status to prospecting
	if err := model.UpdateCandidateStatus(r.Context(), claims.CandidateID, "prospecting"); err != nil {
		slog.Warn("failed to update status to prospecting", "error", err)
	}

	slog.Info("candidate registration completed", "id", claims.CandidateID)

	// Redirect to portal dashboard
	http.Redirect(w, r, "/portal", http.StatusFound)
}

// handleLogin shows the login form
func (h *PublicHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in as candidate
	if claims, err := h.sessionMgr.GetClaimsFromRequest(r); err == nil && claims != nil && claims.IsCandidate {
		http.Redirect(w, r, "/portal", http.StatusFound)
		return
	}

	data := NewPageData("Login")
	portal.Login(data, "").Render(r.Context(), w)
}

// handleLoginSubmit handles login form submission
func (h *PublicHandler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := NewPageData("Login")
		portal.Login(data, "Gagal memproses form").Render(r.Context(), w)
		return
	}

	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := r.FormValue("password")

	if identifier == "" || password == "" {
		data := NewPageData("Login")
		portal.Login(data, "Harap isi email/HP dan password").Render(r.Context(), w)
		return
	}

	// Authenticate
	candidate, err := model.AuthenticateCandidate(r.Context(), identifier, password)
	if err != nil {
		slog.Error("failed to authenticate candidate", "error", err)
		data := NewPageData("Login")
		portal.Login(data, "Terjadi kesalahan sistem").Render(r.Context(), w)
		return
	}

	if candidate == nil {
		data := NewPageData("Login")
		portal.Login(data, "Email/HP atau password salah").Render(r.Context(), w)
		return
	}

	// Create session
	candidateEmail := ""
	if candidate.Email != nil {
		candidateEmail = *candidate.Email
	}
	candidateName := ""
	if candidate.Name != nil {
		candidateName = *candidate.Name
	}

	token, err := h.sessionMgr.CreateCandidateToken(candidate.ID, candidateEmail, candidateName)
	if err != nil {
		slog.Error("failed to create session token", "error", err)
		data := NewPageData("Login")
		portal.Login(data, "Gagal membuat sesi").Render(r.Context(), w)
		return
	}

	h.sessionMgr.SetCookie(w, token)

	slog.Info("candidate logged in", "id", candidate.ID, "email", candidateEmail)

	// Redirect to portal
	http.Redirect(w, r, "/portal", http.StatusFound)
}

// handleLogout handles logout
func (h *PublicHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	h.sessionMgr.ClearCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// handleRequestEmailOTP requests email OTP
func (h *PublicHandler) handleRequestEmailOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get candidate
	candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	if candidate.Email == nil || *candidate.Email == "" {
		http.Error(w, "Email not set", http.StatusBadRequest)
		return
	}

	if candidate.EmailVerified {
		http.Error(w, "Email already verified", http.StatusBadRequest)
		return
	}

	if h.resend == nil {
		http.Error(w, "Email verification not available", http.StatusServiceUnavailable)
		return
	}

	// Generate OTP
	otp, err := model.CreateVerificationToken(r.Context(), candidate.ID, model.TokenTypeEmail)
	if err != nil {
		slog.Error("failed to create email OTP", "error", err)
		http.Error(w, "Failed to create OTP", http.StatusInternalServerError)
		return
	}

	// Send OTP
	if err := h.resend.SendOTP(*candidate.Email, otp); err != nil {
		slog.Error("failed to send email OTP", "error", err)
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OTP sent")
}

// handleConfirmEmailOTP confirms email OTP
func (h *PublicHandler) handleConfirmEmailOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	otp := r.FormValue("otp")
	if otp == "" {
		http.Error(w, "OTP required", http.StatusBadRequest)
		return
	}

	// Verify OTP
	if err := model.VerifyToken(r.Context(), claims.CandidateID, model.TokenTypeEmail, otp); err != nil {
		slog.Warn("email OTP verification failed", "error", err)
		http.Error(w, "Invalid or expired OTP", http.StatusBadRequest)
		return
	}

	// Mark email as verified
	if err := model.SetCandidateEmailVerified(r.Context(), claims.CandidateID); err != nil {
		slog.Error("failed to set email verified", "error", err)
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	slog.Info("candidate email verified", "id", claims.CandidateID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Email verified")
}

// handleRequestPhoneOTP requests phone OTP via WhatsApp
func (h *PublicHandler) handleRequestPhoneOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get candidate
	candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	if candidate.Phone == nil || *candidate.Phone == "" {
		http.Error(w, "Phone not set", http.StatusBadRequest)
		return
	}

	if candidate.PhoneVerified {
		http.Error(w, "Phone already verified", http.StatusBadRequest)
		return
	}

	if h.whatsapp == nil {
		http.Error(w, "Phone verification not available", http.StatusServiceUnavailable)
		return
	}

	// Generate OTP
	otp, err := model.CreateVerificationToken(r.Context(), candidate.ID, model.TokenTypePhone)
	if err != nil {
		slog.Error("failed to create phone OTP", "error", err)
		http.Error(w, "Failed to create OTP", http.StatusInternalServerError)
		return
	}

	// Send OTP via WhatsApp
	if err := h.whatsapp.SendOTP(*candidate.Phone, otp); err != nil {
		slog.Error("failed to send phone OTP", "error", err)
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OTP sent")
}

// handleConfirmPhoneOTP confirms phone OTP
func (h *PublicHandler) handleConfirmPhoneOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	otp := r.FormValue("otp")
	if otp == "" {
		http.Error(w, "OTP required", http.StatusBadRequest)
		return
	}

	// Verify OTP
	if err := model.VerifyToken(r.Context(), claims.CandidateID, model.TokenTypePhone, otp); err != nil {
		slog.Warn("phone OTP verification failed", "error", err)
		http.Error(w, "Invalid or expired OTP", http.StatusBadRequest)
		return
	}

	// Mark phone as verified
	if err := model.SetCandidatePhoneVerified(r.Context(), claims.CandidateID); err != nil {
		slog.Error("failed to set phone verified", "error", err)
		http.Error(w, "Failed to verify phone", http.StatusInternalServerError)
		return
	}

	slog.Info("candidate phone verified", "id", claims.CandidateID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Phone verified")
}

// handleSendEmailOTP sends OTP and returns HTMX fragment
func (h *PublicHandler) handleSendEmailOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	if candidate.Email == nil || *candidate.Email == "" {
		portal.VerifyEmailOTPForm("", "Email tidak tersedia").Render(r.Context(), w)
		return
	}

	if candidate.EmailVerified {
		portal.VerifyEmailSuccess().Render(r.Context(), w)
		return
	}

	if h.resend == nil {
		portal.VerifyEmailOTPForm(*candidate.Email, "Layanan email tidak tersedia").Render(r.Context(), w)
		return
	}

	// Generate OTP
	otp, err := model.CreateVerificationToken(r.Context(), candidate.ID, model.TokenTypeEmail)
	if err != nil {
		slog.Error("failed to create email OTP", "error", err)
		portal.VerifyEmailOTPForm(*candidate.Email, "Gagal membuat kode verifikasi").Render(r.Context(), w)
		return
	}

	// Send OTP
	if err := h.resend.SendOTP(*candidate.Email, otp); err != nil {
		slog.Error("failed to send email OTP", "error", err)
		portal.VerifyEmailOTPForm(*candidate.Email, "Gagal mengirim email").Render(r.Context(), w)
		return
	}

	slog.Info("email OTP sent", "candidate_id", candidate.ID, "email", *candidate.Email)
	portal.VerifyEmailOTPForm(*candidate.Email, "").Render(r.Context(), w)
}

// handleVerifyEmailOTP verifies OTP and returns HTMX fragment
func (h *PublicHandler) handleVerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.sessionMgr.GetClaimsFromRequest(r)
	if err != nil || claims == nil || !claims.IsCandidate {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	candidate, err := model.FindCandidateByID(r.Context(), claims.CandidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	email := ""
	if candidate.Email != nil {
		email = *candidate.Email
	}

	if err := r.ParseForm(); err != nil {
		portal.VerifyEmailOTPForm(email, "Form tidak valid").Render(r.Context(), w)
		return
	}

	otp := r.FormValue("otp")
	if otp == "" {
		portal.VerifyEmailOTPForm(email, "Kode OTP harus diisi").Render(r.Context(), w)
		return
	}

	// Verify OTP
	if err := model.VerifyToken(r.Context(), claims.CandidateID, model.TokenTypeEmail, otp); err != nil {
		slog.Warn("email OTP verification failed", "error", err)
		portal.VerifyEmailOTPForm(email, "Kode OTP salah atau sudah kadaluarsa").Render(r.Context(), w)
		return
	}

	// Mark email as verified
	if err := model.SetCandidateEmailVerified(r.Context(), claims.CandidateID); err != nil {
		slog.Error("failed to set email verified", "error", err)
		portal.VerifyEmailOTPForm(email, "Gagal memverifikasi email").Render(r.Context(), w)
		return
	}

	slog.Info("candidate email verified", "id", claims.CandidateID)
	portal.VerifyEmailSuccess().Render(r.Context(), w)
}

// renderRegistration renders the registration form
func (h *PublicHandler) renderRegistration(w http.ResponseWriter, r *http.Request, currentStep string, regData portal.RegistrationData, errorMsg string) {
	data := NewPageData("Pendaftaran")

	// Build steps based on current step
	steps := []portal.RegistrationStep{
		{Number: "1", Title: "Akun", Status: "pending"},
		{Number: "2", Title: "Data Diri", Status: "pending"},
		{Number: "3", Title: "Pendidikan", Status: "pending"},
		{Number: "4", Title: "Selesai", Status: "pending"},
	}

	stepOrder := []string{"account", "personal", "education", "source"}
	currentIdx := 0
	for i, s := range stepOrder {
		if s == currentStep {
			currentIdx = i
			break
		}
	}

	for i := range steps {
		if i < currentIdx {
			steps[i].Status = "completed"
		} else if i == currentIdx {
			steps[i].Status = "current"
		}
	}

	// Get programs for step 3
	var programs []portal.ProgramOption
	if currentStep == "education" {
		prodiList, err := model.GetActiveProdiWithFees(r.Context())
		if err != nil {
			slog.Error("failed to get programs", "error", err)
		} else {
			for _, p := range prodiList {
				fee := "Biaya tersedia"
				if p.TotalFee > 0 {
					fee = fmt.Sprintf("Rp %d/semester", p.TotalFee)
				}
				programs = append(programs, portal.ProgramOption{
					ID:   p.ID,
					Code: p.Code,
					Name: p.Name,
					Fee:  fee,
				})
			}
		}
	}

	portal.Registration(data, steps, currentStep, regData, programs, sourceTypes, errorMsg).Render(r.Context(), w)
}
