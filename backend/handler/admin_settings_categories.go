package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/model"
	"github.com/idtazkia/stmik-admission-api/templates/admin"
)

func (h *AdminHandler) handleCategoriesSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Kategori")

	// Fetch categories from database
	dbCategories, err := model.ListInteractionCategories(r.Context(), false)
	if err != nil {
		slog.Error("failed to list interaction categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	categories := make([]admin.CategoryItem, len(dbCategories))
	for i, c := range dbCategories {
		categories[i] = admin.CategoryItem{
			ID:        c.ID,
			Name:      c.Name,
			Sentiment: c.Sentiment,
			Count:     "0", // TODO: count interactions using this category
			IsActive:  c.IsActive,
		}
	}

	// Fetch obstacles from database
	dbObstacles, err := model.ListObstacles(r.Context(), false)
	if err != nil {
		slog.Error("failed to list obstacles", "error", err)
		http.Error(w, "Failed to load obstacles", http.StatusInternalServerError)
		return
	}

	// Convert to template type
	obstacles := make([]admin.ObstacleItem, len(dbObstacles))
	for i, o := range dbObstacles {
		obstacles[i] = admin.ObstacleItem{
			ID:       o.ID,
			Name:     o.Name,
			Count:    "0", // TODO: count interactions with this obstacle
			IsActive: o.IsActive,
		}
	}

	admin.SettingsCategories(data, categories, obstacles).Render(r.Context(), w)
}

// handleCreateCategory handles POST /admin/settings/categories
func (h *AdminHandler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	sentiment := r.FormValue("sentiment")

	if name == "" || sentiment == "" {
		http.Error(w, "Name and sentiment are required", http.StatusBadRequest)
		return
	}

	cat, err := model.CreateInteractionCategory(r.Context(), name, sentiment, 0)
	if err != nil {
		slog.Error("failed to create category", "error", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category created", "category_id", cat.ID)

	// Return new category card
	h.renderCategoryCard(w, r, cat.ID)
}

// handleUpdateCategory handles POST /admin/settings/categories/{id}
func (h *AdminHandler) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	sentiment := r.FormValue("sentiment")

	if name == "" || sentiment == "" {
		http.Error(w, "Name and sentiment are required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateInteractionCategory(r.Context(), id, name, sentiment, 0); err != nil {
		slog.Error("failed to update category", "error", err, "category_id", id)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category updated", "category_id", id)

	// Return updated category card
	h.renderCategoryCard(w, r, id)
}

// handleToggleCategoryActive handles POST /admin/settings/categories/{id}/toggle-active
func (h *AdminHandler) handleToggleCategoryActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleInteractionCategoryActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle category active", "error", err, "category_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("interaction category active status toggled", "category_id", id)

	// Return updated category card
	h.renderCategoryCard(w, r, id)
}

// renderCategoryCard renders a single category card for HTMX updates
func (h *AdminHandler) renderCategoryCard(w http.ResponseWriter, r *http.Request, categoryID string) {
	cat, err := model.FindInteractionCategoryByID(r.Context(), categoryID)
	if err != nil {
		slog.Error("failed to find category", "error", err)
		http.Error(w, "Failed to load category", http.StatusInternalServerError)
		return
	}

	item := admin.CategoryItem{
		ID:        cat.ID,
		Name:      cat.Name,
		Sentiment: cat.Sentiment,
		Count:     "0", // TODO: count interactions using this category
		IsActive:  cat.IsActive,
	}

	admin.CategoryCard(item).Render(r.Context(), w)
}

// handleCreateObstacle handles POST /admin/settings/obstacles
func (h *AdminHandler) handleCreateObstacle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	obs, err := model.CreateObstacle(r.Context(), name, nil, 0)
	if err != nil {
		slog.Error("failed to create obstacle", "error", err)
		http.Error(w, "Failed to create obstacle", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle created", "obstacle_id", obs.ID)

	// Return new obstacle card
	h.renderObstacleCard(w, r, obs.ID)
}

// handleUpdateObstacle handles POST /admin/settings/obstacles/{id}
func (h *AdminHandler) handleUpdateObstacle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateObstacle(r.Context(), id, name, nil, 0); err != nil {
		slog.Error("failed to update obstacle", "error", err, "obstacle_id", id)
		http.Error(w, "Failed to update obstacle", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle updated", "obstacle_id", id)

	// Return updated obstacle card
	h.renderObstacleCard(w, r, id)
}

// handleToggleObstacleActive handles POST /admin/settings/obstacles/{id}/toggle-active
func (h *AdminHandler) handleToggleObstacleActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleObstacleActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle obstacle active", "error", err, "obstacle_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("obstacle active status toggled", "obstacle_id", id)

	// Return updated obstacle card
	h.renderObstacleCard(w, r, id)
}

// renderObstacleCard renders a single obstacle card for HTMX updates
func (h *AdminHandler) renderObstacleCard(w http.ResponseWriter, r *http.Request, obstacleID string) {
	obs, err := model.FindObstacleByID(r.Context(), obstacleID)
	if err != nil {
		slog.Error("failed to find obstacle", "error", err)
		http.Error(w, "Failed to load obstacle", http.StatusInternalServerError)
		return
	}

	item := admin.ObstacleItem{
		ID:       obs.ID,
		Name:     obs.Name,
		Count:    "0", // TODO: count interactions with this obstacle
		IsActive: obs.IsActive,
	}

	admin.ObstacleCard(item).Render(r.Context(), w)
}
