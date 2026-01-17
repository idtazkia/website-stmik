package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
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
