package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleUsersSettings(w http.ResponseWriter, r *http.Request) {
	data := NewPageDataWithUser(r.Context(), "Users")

	// Fetch users from database
	dbUsers, err := model.ListUsers(r.Context(), "", false)
	if err != nil {
		slog.Error("failed to list users", "error", err)
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	// Fetch supervisors for dropdown
	dbSupervisors, err := model.ListSupervisors(r.Context())
	if err != nil {
		slog.Error("failed to list supervisors", "error", err)
		http.Error(w, "Failed to load supervisors", http.StatusInternalServerError)
		return
	}

	// Convert supervisors to template type
	supervisors := make([]admin.SupervisorOption, len(dbSupervisors))
	for i, s := range dbSupervisors {
		supervisors[i] = admin.SupervisorOption{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	// Convert to template type
	users := make([]admin.UserItem, len(dbUsers))
	for i, u := range dbUsers {
		status := "inactive"
		if u.IsActive {
			status = "active"
		}

		supervisor := "-"
		if u.SupervisorName != nil {
			supervisor = *u.SupervisorName
		}

		supervisorID := ""
		if u.IDSupervisor != nil {
			supervisorID = *u.IDSupervisor
		}

		lastLogin := "Belum pernah"
		if u.LastLoginAt != nil {
			lastLogin = formatRelativeTime(*u.LastLoginAt)
		}

		users[i] = admin.UserItem{
			ID:           u.ID,
			Name:         u.Name,
			Email:        u.Email,
			Role:         u.Role,
			Supervisor:   supervisor,
			IDSupervisor: supervisorID,
			Status:       status,
			LastLogin:    lastLogin,
		}
	}

	// Fetch stats
	counts, err := model.CountUsersByRole(r.Context())
	if err != nil {
		slog.Error("failed to count users", "error", err)
	}
	stats := admin.UserStats{
		Total:      counts["admin"] + counts["supervisor"] + counts["consultant"] + counts["finance"] + counts["academic"],
		Admin:      counts["admin"],
		Supervisor: counts["supervisor"],
		Consultant: counts["consultant"],
		Finance:    counts["finance"],
		Academic:   counts["academic"],
	}

	admin.SettingsUsers(data, users, stats, supervisors).Render(r.Context(), w)
}

// handleUpdateUserRole handles POST /admin/settings/users/{id}/role
func (h *AdminHandler) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	role := r.FormValue("role")
	if role == "" {
		http.Error(w, "Role is required", http.StatusBadRequest)
		return
	}

	if err := model.UpdateUserRole(r.Context(), id, role); err != nil {
		slog.Error("failed to update user role", "error", err, "user_id", id, "role", role)
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	slog.Info("user role updated", "user_id", id, "role", role)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// handleUpdateUserSupervisor handles POST /admin/settings/users/{id}/supervisor
func (h *AdminHandler) handleUpdateUserSupervisor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	supervisorID := r.FormValue("supervisor_id")
	var supIDPtr *string
	if supervisorID != "" {
		supIDPtr = &supervisorID
	}

	if err := model.UpdateUserSupervisor(r.Context(), id, supIDPtr); err != nil {
		slog.Error("failed to update user supervisor", "error", err, "user_id", id, "supervisor_id", supervisorID)
		http.Error(w, "Failed to update supervisor", http.StatusInternalServerError)
		return
	}

	slog.Info("user supervisor updated", "user_id", id, "supervisor_id", supervisorID)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// handleToggleUserActive handles POST /admin/settings/users/{id}/toggle-active
func (h *AdminHandler) handleToggleUserActive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := model.ToggleUserActive(r.Context(), id); err != nil {
		slog.Error("failed to toggle user active", "error", err, "user_id", id)
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	slog.Info("user active status toggled", "user_id", id)

	// Return updated row
	h.renderUserRow(w, r, id)
}

// renderUserRow renders a single user row for HTMX updates
func (h *AdminHandler) renderUserRow(w http.ResponseWriter, r *http.Request, userID string) {
	// Fetch user by ID with supervisor info
	dbUsers, err := model.ListUsers(r.Context(), "", false)
	if err != nil {
		slog.Error("failed to list users", "error", err)
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Fetch supervisors for dropdown
	dbSupervisors, err := model.ListSupervisors(r.Context())
	if err != nil {
		slog.Error("failed to list supervisors", "error", err)
		http.Error(w, "Failed to load supervisors", http.StatusInternalServerError)
		return
	}

	// Convert supervisors to template type
	supervisors := make([]admin.SupervisorOption, len(dbSupervisors))
	for i, s := range dbSupervisors {
		supervisors[i] = admin.SupervisorOption{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	// Find the user
	var userItem admin.UserItem
	found := false
	for _, u := range dbUsers {
		if u.ID == userID {
			status := "inactive"
			if u.IsActive {
				status = "active"
			}

			supervisor := "-"
			if u.SupervisorName != nil {
				supervisor = *u.SupervisorName
			}

			supervisorID := ""
			if u.IDSupervisor != nil {
				supervisorID = *u.IDSupervisor
			}

			lastLogin := "Belum pernah"
			if u.LastLoginAt != nil {
				lastLogin = formatRelativeTime(*u.LastLoginAt)
			}

			userItem = admin.UserItem{
				ID:           u.ID,
				Name:         u.Name,
				Email:        u.Email,
				Role:         u.Role,
				Supervisor:   supervisor,
				IDSupervisor: supervisorID,
				Status:       status,
				LastLogin:    lastLogin,
			}
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	admin.UserRow(userItem, supervisors).Render(r.Context(), w)
}

// formatRelativeTime formats a time as relative to now
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Baru saja"
	}
	if diff < time.Hour {
		return "Beberapa menit lalu"
	}
	if diff < 24*time.Hour {
		return "Hari ini"
	}
	if diff < 48*time.Hour {
		return "Kemarin"
	}
	days := int(diff.Hours() / 24)
	if days < 7 {
		return formatDays(days) + " lalu"
	}
	return t.Format("2 Jan 2006")
}

func formatDays(days int) string {
	return fmt.Sprintf("%d hari", days)
}
