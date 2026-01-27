package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleCandidates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Kandidat")

	// Get current user for role-based visibility
	claims := GetUserClaims(ctx)
	var visibilityConsultantID, visibilitySupervisorID *string
	if claims != nil {
		switch claims.Role {
		case "consultant":
			visibilityConsultantID = &claims.UserID
		case "supervisor":
			visibilitySupervisorID = &claims.UserID
		// admin sees all - no filter
		}
	}

	// Parse filter parameters from query string
	filters := model.CandidateListFilters{
		Status:       r.URL.Query().Get("status"),
		ConsultantID: r.URL.Query().Get("consultant_id"),
		ProdiID:      r.URL.Query().Get("prodi_id"),
		CampaignID:   r.URL.Query().Get("campaign_id"),
		SourceType:   r.URL.Query().Get("source_type"),
		Search:       r.URL.Query().Get("search"),
		SortBy:       r.URL.Query().Get("sort_by"),
		SortOrder:    r.URL.Query().Get("sort_order"),
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// Fetch candidates
	result, err := model.ListCandidates(ctx, filters, visibilityConsultantID, visibilitySupervisorID)
	if err != nil {
		log.Printf("Error listing candidates: %v", err)
		http.Error(w, "Failed to load candidates", http.StatusInternalServerError)
		return
	}

	// Fetch stats
	stats, err := model.GetCandidateStatusStats(ctx, visibilityConsultantID, visibilitySupervisorID)
	if err != nil {
		log.Printf("Error getting candidate stats: %v", err)
		http.Error(w, "Failed to load stats", http.StatusInternalServerError)
		return
	}

	// Fetch filter options
	consultants, err := model.ListUsers(ctx, "consultant", true)
	if err != nil {
		log.Printf("Error listing consultants: %v", err)
		consultants = []model.UserWithSupervisor{}
	}

	prodis, err := model.ListProdis(ctx, true)
	if err != nil {
		log.Printf("Error listing prodis: %v", err)
		prodis = []model.Prodi{}
	}

	campaigns, err := model.ListCampaigns(ctx, true)
	if err != nil {
		log.Printf("Error listing campaigns: %v", err)
		campaigns = []model.Campaign{}
	}

	// Convert model data to template types
	candidateItems := make([]admin.CandidateListItem, len(result.Candidates))
	for i, c := range result.Candidates {
		name := ptrToString(c.Name)
		email := ptrToString(c.Email)
		phone := ptrToString(c.Phone)
		prodiName := ptrToString(c.ProdiName)
		consultantName := ptrToString(c.ConsultantName)
		campaignName := ptrToString(c.CampaignName)
		sourceType := ptrToString(c.SourceType)

		candidateItems[i] = admin.CandidateListItem{
			ID:             c.ID,
			Name:           name,
			Email:          email,
			Phone:          phone,
			Status:         c.Status,
			StatusLabel:    candidateStatusLabel(c.Status),
			ProdiName:      prodiName,
			ConsultantName: consultantName,
			CampaignName:   campaignName,
			SourceType:     sourceType,
			SourceLabel:    candidateSourceLabel(sourceType),
			CreatedAt:      c.CreatedAt.Format("2 Jan 2006"),
		}
	}

	// Convert filter options to template types
	consultantOpts := make([]admin.FilterOption, len(consultants))
	for i, c := range consultants {
		name := c.Name
		if name == "" {
			name = c.Email
		}
		consultantOpts[i] = admin.FilterOption{Value: c.ID, Label: name}
	}

	prodiOpts := make([]admin.FilterOption, len(prodis))
	for i, p := range prodis {
		prodiOpts[i] = admin.FilterOption{Value: p.ID, Label: p.Name}
	}

	campaignOpts := make([]admin.FilterOption, len(campaigns))
	for i, c := range campaigns {
		campaignOpts[i] = admin.FilterOption{Value: c.ID, Label: c.Name}
	}

	// Build template data
	listData := admin.CandidateListData{
		Candidates: candidateItems,
		Stats: admin.CandidateListStats{
			Total:       stats.Total,
			Registered:  stats.Registered,
			Prospecting: stats.Prospecting,
			Committed:   stats.Committed,
			Enrolled:    stats.Enrolled,
			Lost:        stats.Lost,
		},
		Filters: admin.CandidateFilters{
			Status:       filters.Status,
			ConsultantID: filters.ConsultantID,
			ProdiID:      filters.ProdiID,
			CampaignID:   filters.CampaignID,
			SourceType:   filters.SourceType,
			Search:       filters.Search,
		},
		Total:       result.Total,
		Limit:       result.Limit,
		Offset:      result.Offset,
		Consultants: consultantOpts,
		Prodis:      prodiOpts,
		Campaigns:   campaignOpts,
	}

	// Check if this is an HTMX request for table body only
	if r.Header.Get("HX-Request") == "true" {
		admin.CandidatesTableBody(candidateItems, result.Total, result.Limit, result.Offset).Render(ctx, w)
		return
	}

	admin.Candidates(data, listData).Render(ctx, w)
}

func (h *AdminHandler) handleCandidateDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	// Fetch candidate with all related data
	candidate, err := model.GetCandidateDetailData(ctx, id)
	if err != nil {
		log.Printf("Error fetching candidate detail: %v", err)
		http.Error(w, "Failed to load candidate", http.StatusInternalServerError)
		return
	}
	if candidate == nil {
		http.NotFound(w, r)
		return
	}

	// Build page title
	title := "Detail Kandidat"
	if candidate.Name != nil && *candidate.Name != "" {
		title = "Detail Kandidat - " + *candidate.Name
	}
	data := NewPageDataWithUser(ctx, title)

	// Convert to template type
	graduationYear := ""
	if candidate.GraduationYear != nil {
		graduationYear = strconv.Itoa(*candidate.GraduationYear)
	}

	c := admin.CandidateDetail{
		ID:                  candidate.ID,
		Name:                ptrToString(candidate.Name),
		Email:               ptrToString(candidate.Email),
		Phone:               ptrToString(candidate.Phone),
		WhatsApp:            ptrToString(candidate.Phone), // Use phone as WhatsApp for now
		Address:             ptrToString(candidate.Address),
		City:                ptrToString(candidate.City),
		Province:            ptrToString(candidate.Province),
		HighSchool:          ptrToString(candidate.HighSchool),
		GraduationYear:      graduationYear,
		ProdiName:           ptrToString(candidate.ProdiName),
		SourceType:          candidateSourceLabel(ptrToString(candidate.SourceType)),
		SourceDetail:        ptrToString(candidate.SourceDetail),
		CampaignName:        ptrToString(candidate.CampaignName),
		ReferrerName:        ptrToString(candidate.ReferrerName),
		Status:              candidate.Status,
		ConsultantName:      ptrToString(candidate.ConsultantName),
		RegistrationFeePaid: false, // TODO: check actual payment status when billing is implemented
		CreatedAt:           candidate.CreatedAt.Format("2 Jan 2006"),
	}

	// Fetch interactions from database
	dbInteractions, err := model.ListInteractionsByCandidate(ctx, id)
	if err != nil {
		log.Printf("Error fetching interactions: %v", err)
		dbInteractions = []model.InteractionWithDetails{}
	}

	// Convert to template type
	interactions := make([]admin.Interaction, len(dbInteractions))
	for i, intr := range dbInteractions {
		nextFollowup := ""
		if intr.NextFollowupDate != nil {
			nextFollowup = intr.NextFollowupDate.Format("2 Jan 2006")
		}

		interactions[i] = admin.Interaction{
			ID:                   intr.ID,
			Channel:              intr.Channel,
			Category:             ptrToString(intr.CategoryName),
			CategorySentiment:    ptrToString(intr.CategorySentiment),
			Obstacle:             ptrToString(intr.ObstacleName),
			Remarks:              intr.Remarks,
			NextFollowupDate:     nextFollowup,
			SupervisorSuggestion: ptrToString(intr.SupervisorSuggestion),
			SuggestionRead:       intr.SuggestionReadAt != nil,
			ConsultantName:       intr.ConsultantName,
			CreatedAt:            intr.CreatedAt.Format("2 Jan 2006 15:04"),
		}
	}

	admin.KandidatDetail(data, c, interactions).Render(ctx, w)
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func candidateStatusLabel(status string) string {
	switch status {
	case "registered":
		return "Terdaftar"
	case "prospecting":
		return "Dalam Proses"
	case "committed":
		return "Berkomitmen"
	case "enrolled":
		return "Diterima"
	case "lost":
		return "Tidak Lanjut"
	default:
		return status
	}
}

func candidateSourceLabel(sourceType string) string {
	switch sourceType {
	case "instagram":
		return "Instagram"
	case "google":
		return "Google"
	case "tiktok":
		return "TikTok"
	case "youtube":
		return "YouTube"
	case "expo":
		return "Expo/Pameran"
	case "school_visit":
		return "Kunjungan Sekolah"
	case "friend_family":
		return "Teman/Keluarga"
	case "teacher_alumni":
		return "Guru/Alumni"
	case "walkin":
		return "Datang Langsung"
	case "other":
		return "Lainnya"
	default:
		return sourceType
	}
}

func (h *AdminHandler) handleReassignCandidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Only admin and supervisor can reassign
	if claims.Role != "admin" && claims.Role != "supervisor" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	newConsultantID := r.FormValue("consultant_id")
	if newConsultantID == "" {
		http.Error(w, "Consultant ID is required", http.StatusBadRequest)
		return
	}

	// Reassign the candidate
	err := model.ReassignCandidate(ctx, candidateID, newConsultantID, claims.UserID)
	if err != nil {
		log.Printf("Error reassigning candidate: %v", err)
		http.Error(w, "Failed to reassign candidate", http.StatusInternalServerError)
		return
	}

	log.Printf("Candidate %s reassigned to consultant %s by %s", candidateID, newConsultantID, claims.UserID)

	// Redirect back to candidate detail
	http.Redirect(w, r, "/admin/candidates/"+candidateID, http.StatusSeeOther)
}

func (h *AdminHandler) handleGetConsultantsForReassign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only admin and supervisor can view consultant list for reassign
	if claims.Role != "admin" && claims.Role != "supervisor" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	consultants, err := model.ListConsultantsWithWorkload(ctx)
	if err != nil {
		log.Printf("Error listing consultants: %v", err)
		http.Error(w, "Failed to load consultants", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	items := make([]admin.ConsultantWorkloadItem, len(consultants))
	for i, c := range consultants {
		items[i] = admin.ConsultantWorkloadItem{
			ID:          c.ID,
			Name:        c.Name,
			Email:       c.Email,
			ActiveCount: c.ActiveCount,
			TotalCount:  c.TotalCount,
		}
	}

	// Get candidate ID from query for current selection
	candidateID := r.URL.Query().Get("candidate_id")
	currentConsultantID := ""
	if candidateID != "" {
		candidate, err := model.FindCandidateByID(ctx, candidateID)
		if err == nil && candidate != nil && candidate.AssignedConsultantID != nil {
			currentConsultantID = *candidate.AssignedConsultantID
		}
	}

	admin.ReassignModal(candidateID, currentConsultantID, items).Render(ctx, w)
}

func (h *AdminHandler) handleGetLostModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get candidate info
	candidate, err := model.FindCandidateByID(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Don't allow marking already lost or enrolled candidates
	if candidate.Status == "lost" || candidate.Status == "enrolled" {
		http.Error(w, "Cannot mark this candidate as lost", http.StatusBadRequest)
		return
	}

	// Get lost reasons
	reasons, err := model.ListLostReasons(ctx, true)
	if err != nil {
		log.Printf("Error listing lost reasons: %v", err)
		http.Error(w, "Failed to load lost reasons", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	reasonItems := make([]admin.LostReasonItem, len(reasons))
	for i, r := range reasons {
		reasonItems[i] = admin.LostReasonItem{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	candidateName := ""
	if candidate.Name != nil {
		candidateName = *candidate.Name
	}

	admin.LostCandidateModal(candidateID, candidateName, reasonItems).Render(ctx, w)
}

func (h *AdminHandler) handleMarkLost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	lostReasonID := r.FormValue("lost_reason_id")
	if lostReasonID == "" {
		http.Error(w, "Lost reason is required", http.StatusBadRequest)
		return
	}

	// Get candidate to validate
	candidate, err := model.FindCandidateByID(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Don't allow marking already lost or enrolled candidates
	if candidate.Status == "lost" || candidate.Status == "enrolled" {
		http.Error(w, "Cannot mark this candidate as lost", http.StatusBadRequest)
		return
	}

	// Mark as lost
	err = model.MarkCandidateLost(ctx, candidateID, lostReasonID)
	if err != nil {
		log.Printf("Error marking candidate as lost: %v", err)
		http.Error(w, "Failed to mark candidate as lost", http.StatusInternalServerError)
		return
	}

	// Log the action as an interaction
	_, err = model.CreateInteraction(ctx, candidateID, claims.UserID, "system", nil, nil, "Kandidat ditandai sebagai lost", nil, nil)
	if err != nil {
		log.Printf("Error logging lost interaction: %v", err)
		// Don't fail the request, just log
	}

	log.Printf("Candidate %s marked as lost by %s with reason %s", candidateID, claims.UserID, lostReasonID)

	// Send HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/candidates/"+candidateID)
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) handleGetCommitmentModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get candidate info
	candidate, err := model.GetCandidateDetailData(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Only prospecting candidates can be committed
	if candidate.Status != "prospecting" {
		http.Error(w, "Only prospecting candidates can be committed", http.StatusBadRequest)
		return
	}

	// Get current academic year (simple logic: year/year+1)
	academicYear := getCurrentAcademicYear()

	// Get tuition fee for candidate's prodi
	tuitionAmount, err := model.GetTuitionFeeForCandidate(ctx, candidateID, academicYear)
	if err != nil {
		log.Printf("Error getting tuition fee: %v", err)
	}

	// Get dormitory fee (global)
	dormAmount, err := model.GetDormitoryFee(ctx, academicYear)
	if err != nil {
		log.Printf("Error getting dormitory fee: %v", err)
	}

	candidateName := ""
	if candidate.Name != nil {
		candidateName = *candidate.Name
	}
	prodiName := ""
	if candidate.ProdiName != nil {
		prodiName = *candidate.ProdiName
	}

	data := admin.CommitmentFormData{
		CandidateID:   candidateID,
		CandidateName: candidateName,
		ProdiName:     prodiName,
		TuitionAmount: tuitionAmount,
		DormAmount:    dormAmount,
		AcademicYear:  academicYear,
	}

	admin.CommitmentModal(data).Render(ctx, w)
}

func (h *AdminHandler) handleCommitCandidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Get candidate to validate
	candidate, err := model.GetCandidateDetailData(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Only prospecting candidates can be committed
	if candidate.Status != "prospecting" {
		http.Error(w, "Only prospecting candidates can be committed", http.StatusBadRequest)
		return
	}

	academicYear := getCurrentAcademicYear()

	// Parse form values
	tuitionInstallments := 1 // Currently only 1 is supported
	if instStr := r.FormValue("tuition_installments"); instStr != "" {
		if inst, err := strconv.Atoi(instStr); err == nil {
			tuitionInstallments = inst
		}
	}

	includeDorm := r.FormValue("include_dorm") == "true"
	dormInstallments := 1
	if instStr := r.FormValue("dorm_installments"); instStr != "" {
		if inst, err := strconv.Atoi(instStr); err == nil {
			dormInstallments = inst
		}
	}

	// Get tuition fee
	tuitionAmount, err := model.GetTuitionFeeForCandidate(ctx, candidateID, academicYear)
	if err != nil {
		log.Printf("Error getting tuition fee: %v", err)
		http.Error(w, "Failed to get tuition fee", http.StatusInternalServerError)
		return
	}

	// Create tuition billing
	if tuitionAmount > 0 {
		_, err = model.CreateTuitionBilling(ctx, candidateID, tuitionAmount, academicYear, tuitionInstallments)
		if err != nil {
			log.Printf("Error creating tuition billing: %v", err)
			http.Error(w, "Failed to create tuition billing", http.StatusInternalServerError)
			return
		}
		log.Printf("Created tuition billing for candidate %s: %d", candidateID, tuitionAmount)
	}

	// Create dormitory billings if included
	if includeDorm {
		dormAmount, err := model.GetDormitoryFee(ctx, academicYear)
		if err != nil {
			log.Printf("Error getting dormitory fee: %v", err)
		}
		if dormAmount > 0 {
			_, err = model.CreateDormitoryBillings(ctx, candidateID, dormAmount, academicYear, dormInstallments)
			if err != nil {
				log.Printf("Error creating dormitory billings: %v", err)
				http.Error(w, "Failed to create dormitory billings", http.StatusInternalServerError)
				return
			}
			log.Printf("Created dormitory billings for candidate %s: %d (%d installments)", candidateID, dormAmount, dormInstallments)
		}
	}

	// Update candidate status to committed (this also triggers commission creation)
	err = model.UpdateCandidateStatus(ctx, candidateID, "committed")
	if err != nil {
		log.Printf("Error updating candidate status: %v", err)
		http.Error(w, "Failed to update candidate status", http.StatusInternalServerError)
		return
	}

	// Log the action as an interaction
	remarks := "Kandidat berkomitmen untuk mendaftar"
	_, err = model.CreateInteraction(ctx, candidateID, claims.UserID, "system", nil, nil, remarks, nil, nil)
	if err != nil {
		log.Printf("Error logging commitment interaction: %v", err)
	}

	log.Printf("Candidate %s committed by %s", candidateID, claims.UserID)

	// Send HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/candidates/"+candidateID)
	w.WriteHeader(http.StatusOK)
}

// getCurrentAcademicYear returns the current academic year in format "2024/2025"
func getCurrentAcademicYear() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	// Academic year starts in August
	if month < 8 {
		year--
	}
	return strconv.Itoa(year) + "/" + strconv.Itoa(year+1)
}

func (h *AdminHandler) handleGetEnrollmentModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get candidate info
	candidate, err := model.GetCandidateDetailData(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Only committed candidates can be enrolled
	if candidate.Status != "committed" {
		http.Error(w, "Only committed candidates can be enrolled", http.StatusBadRequest)
		return
	}

	// Validate enrollment requirements
	validation, err := model.ValidateEnrollment(ctx, candidateID)
	if err != nil {
		log.Printf("Error validating enrollment: %v", err)
		http.Error(w, "Failed to validate enrollment", http.StatusInternalServerError)
		return
	}

	candidateName := ""
	if candidate.Name != nil {
		candidateName = *candidate.Name
	}
	prodiName := ""
	if candidate.ProdiName != nil {
		prodiName = *candidate.ProdiName
	}

	data := admin.EnrollmentFormData{
		CandidateID:          candidateID,
		CandidateName:        candidateName,
		ProdiName:            prodiName,
		CanEnroll:            validation.CanEnroll,
		RegistrationFeePaid:  validation.RegistrationFeePaid,
		TuitionFirstPaid:     validation.TuitionFirstPaid,
		RequiredDocsApproved: validation.RequiredDocsApproved,
		MissingDocs:          validation.MissingDocs,
		UnpaidBillings:       validation.UnpaidBillings,
	}

	admin.EnrollmentModal(data).Render(ctx, w)
}

func (h *AdminHandler) handleEnrollCandidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Get current user
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Get candidate to validate
	candidate, err := model.GetCandidateDetailData(ctx, candidateID)
	if err != nil || candidate == nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	// Only committed candidates can be enrolled
	if candidate.Status != "committed" {
		http.Error(w, "Only committed candidates can be enrolled", http.StatusBadRequest)
		return
	}

	// Enroll the candidate (this validates requirements, generates NIM, referral code, and updates status)
	enrolledCandidate, err := model.EnrollCandidate(ctx, candidateID)
	if err != nil {
		log.Printf("Error enrolling candidate: %v", err)
		http.Error(w, "Failed to enroll candidate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the action as an interaction
	var nim string
	if enrolledCandidate.NIM != nil {
		nim = *enrolledCandidate.NIM
	}
	remarks := "Kandidat berhasil di-enroll dengan NIM: " + nim
	_, err = model.CreateInteraction(ctx, candidateID, claims.UserID, "system", nil, nil, remarks, nil, nil)
	if err != nil {
		log.Printf("Error logging enrollment interaction: %v", err)
	}

	log.Printf("Candidate %s enrolled by %s with NIM %s", candidateID, claims.UserID, nim)

	// Send HTMX redirect
	w.Header().Set("HX-Redirect", "/admin/candidates/"+candidateID)
	w.WriteHeader(http.StatusOK)
}
