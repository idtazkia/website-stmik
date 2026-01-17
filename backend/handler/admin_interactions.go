package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

func (h *AdminHandler) handleInteractionForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Fetch candidate
	candidate, err := model.GetCandidateDetailData(ctx, candidateID)
	if err != nil {
		log.Printf("Error fetching candidate: %v", err)
		http.Error(w, "Failed to load candidate", http.StatusInternalServerError)
		return
	}
	if candidate == nil {
		http.NotFound(w, r)
		return
	}

	// Build page title
	title := "Log Interaksi"
	if candidate.Name != nil && *candidate.Name != "" {
		title = "Log Interaksi - " + *candidate.Name
	}
	data := NewPageDataWithUser(ctx, title)

	// Build candidate summary for template
	c := admin.CandidateSummary{
		ID:        candidate.ID,
		Name:      ptrToString(candidate.Name),
		ProdiName: ptrToString(candidate.ProdiName),
		WhatsApp:  ptrToString(candidate.Phone),
		Status:    candidate.Status,
	}

	// Fetch categories from database
	dbCategories, err := model.ListInteractionCategories(ctx, true)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		dbCategories = []model.InteractionCategory{}
	}

	// Convert to template types with icons based on sentiment
	categories := make([]admin.InteractionCategoryOption, len(dbCategories))
	for i, cat := range dbCategories {
		icon := "🔵" // default
		switch cat.Sentiment {
		case "positive":
			icon = "😊"
		case "neutral":
			icon = "🤔"
		case "negative":
			icon = "😟"
		}
		categories[i] = admin.InteractionCategoryOption{
			Value:     cat.ID,
			Label:     cat.Name,
			Icon:      icon,
			Sentiment: cat.Sentiment,
		}
	}

	// Fetch obstacles from database
	dbObstacles, err := model.ListObstacles(ctx, true)
	if err != nil {
		log.Printf("Error fetching obstacles: %v", err)
		dbObstacles = []model.Obstacle{}
	}

	// Convert to template types
	obstacles := make([]admin.ObstacleOption, len(dbObstacles))
	for i, obs := range dbObstacles {
		obstacles[i] = admin.ObstacleOption{
			Value: obs.ID,
			Label: obs.Name,
		}
	}

	admin.InteractionForm(data, c, categories, obstacles).Render(ctx, w)
}

func (h *AdminHandler) handleCreateInteraction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	// Parse form
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get current user (consultant)
	claims := GetUserClaims(ctx)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract form values
	channel := r.FormValue("channel")
	categoryID := r.FormValue("category")
	obstacleID := r.FormValue("obstacle")
	remarks := r.FormValue("remarks")
	nextFollowupStr := r.FormValue("next_followup_date")
	nextAction := r.FormValue("next_action")

	// Validate required fields
	if channel == "" {
		http.Error(w, "Channel is required", http.StatusBadRequest)
		return
	}
	if remarks == "" {
		http.Error(w, "Remarks is required", http.StatusBadRequest)
		return
	}

	// Convert to pointers
	var categoryPtr, obstaclePtr, nextActionPtr *string
	if categoryID != "" {
		categoryPtr = &categoryID
	}
	if obstacleID != "" {
		obstaclePtr = &obstacleID
	}
	if nextAction != "" {
		nextActionPtr = &nextAction
	}

	// Parse followup date
	var nextFollowupDate *time.Time
	if nextFollowupStr != "" {
		t, err := time.Parse("2006-01-02", nextFollowupStr)
		if err == nil {
			nextFollowupDate = &t
		}
	}

	// Create interaction
	_, err := model.CreateInteraction(ctx, candidateID, claims.UserID, channel, categoryPtr, obstaclePtr, remarks, nextFollowupDate, nextActionPtr)
	if err != nil {
		log.Printf("Error creating interaction: %v", err)
		http.Error(w, "Failed to create interaction", http.StatusInternalServerError)
		return
	}

	// Check if this is save_and_next action
	action := r.FormValue("action")
	if action == "save_and_next" {
		// TODO: Redirect to next candidate in queue
		// For now, redirect back to candidate detail
		http.Redirect(w, r, "/admin/candidates/"+candidateID, http.StatusSeeOther)
		return
	}

	// Redirect back to candidate detail
	http.Redirect(w, r, "/admin/candidates/"+candidateID, http.StatusSeeOther)
}
