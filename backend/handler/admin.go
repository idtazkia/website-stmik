package handler

import (
	"net/http"

	"github.com/idtazkia/stmik-admission-api/auth"
)

// AdminHandler handles all admin routes
type AdminHandler struct {
	sessionMgr *auth.SessionManager
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(sessionMgr *auth.SessionManager) *AdminHandler {
	return &AdminHandler{sessionMgr: sessionMgr}
}

// RegisterRoutes registers all admin routes to the mux
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	// Helper to wrap handlers with auth middleware
	protected := func(handler http.HandlerFunc) http.Handler {
		return RequireAuth(h.sessionMgr, handler)
	}

	// Dashboard
	mux.Handle("GET /admin", protected(h.handleDashboard))
	mux.Handle("GET /admin/", protected(h.handleDashboard))

	// Consultant personal dashboard
	mux.Handle("GET /admin/my-dashboard", protected(h.handleConsultantDashboard))

	// Candidates
	mux.Handle("GET /admin/candidates", protected(h.handleCandidates))
	mux.Handle("GET /admin/candidates/{id}", protected(h.handleCandidateDetail))
	mux.Handle("GET /admin/candidates/{id}/interaction", protected(h.handleInteractionForm))

	// Documents
	mux.Handle("GET /admin/documents", protected(h.handleDocumentReview))

	// Marketing
	mux.Handle("GET /admin/campaigns", protected(h.handleCampaigns))
	mux.Handle("GET /admin/referrers", protected(h.handleReferrers))
	mux.Handle("GET /admin/referral-claims", protected(h.handleReferralClaims))
	mux.Handle("GET /admin/commissions", protected(h.handleCommissions))

	// Reports
	mux.Handle("GET /admin/reports/funnel", protected(h.handleFunnelReport))
	mux.Handle("GET /admin/reports/consultants", protected(h.handleConsultantsReport))
	mux.Handle("GET /admin/reports/campaigns", protected(h.handleCampaignsReport))

	// Settings - Users
	mux.Handle("GET /admin/settings/users", protected(h.handleUsersSettings))
	mux.Handle("POST /admin/settings/users/{id}/role", protected(h.handleUpdateUserRole))
	mux.Handle("POST /admin/settings/users/{id}/supervisor", protected(h.handleUpdateUserSupervisor))
	mux.Handle("POST /admin/settings/users/{id}/toggle-active", protected(h.handleToggleUserActive))

	// Settings - Programs
	mux.Handle("GET /admin/settings/programs", protected(h.handleProgramsSettings))
	mux.Handle("POST /admin/settings/programs", protected(h.handleCreateProgram))
	mux.Handle("POST /admin/settings/programs/{id}", protected(h.handleUpdateProgram))
	mux.Handle("POST /admin/settings/programs/{id}/toggle-active", protected(h.handleToggleProgramActive))

	// Settings - Categories & Obstacles
	mux.Handle("GET /admin/settings/categories", protected(h.handleCategoriesSettings))
	mux.Handle("POST /admin/settings/categories", protected(h.handleCreateCategory))
	mux.Handle("POST /admin/settings/categories/{id}", protected(h.handleUpdateCategory))
	mux.Handle("POST /admin/settings/categories/{id}/toggle-active", protected(h.handleToggleCategoryActive))
	mux.Handle("POST /admin/settings/obstacles", protected(h.handleCreateObstacle))
	mux.Handle("POST /admin/settings/obstacles/{id}", protected(h.handleUpdateObstacle))
	mux.Handle("POST /admin/settings/obstacles/{id}/toggle-active", protected(h.handleToggleObstacleActive))

	// Settings - Fees
	mux.Handle("GET /admin/settings/fees", protected(h.handleFeesSettings))
	mux.Handle("POST /admin/settings/fees", protected(h.handleCreateFeeStructure))
	mux.Handle("POST /admin/settings/fees/{id}", protected(h.handleUpdateFeeStructure))
	mux.Handle("POST /admin/settings/fees/{id}/toggle-active", protected(h.handleToggleFeeStructureActive))

	// Settings - Campaigns
	mux.Handle("GET /admin/settings/campaigns", protected(h.handleCampaignsSettings))
	mux.Handle("POST /admin/settings/campaigns", protected(h.handleCreateCampaign))
	mux.Handle("POST /admin/settings/campaigns/{id}", protected(h.handleUpdateCampaign))
	mux.Handle("POST /admin/settings/campaigns/{id}/toggle-active", protected(h.handleToggleCampaignActive))

	// Settings - Rewards
	mux.Handle("GET /admin/settings/rewards", protected(h.handleRewardsSettings))
	mux.Handle("POST /admin/settings/rewards", protected(h.handleCreateRewardConfig))
	mux.Handle("POST /admin/settings/rewards/{id}", protected(h.handleUpdateRewardConfig))
	mux.Handle("POST /admin/settings/rewards/{id}/toggle-active", protected(h.handleToggleRewardConfigActive))
	mux.Handle("POST /admin/settings/mgm-rewards", protected(h.handleCreateMGMRewardConfig))
	mux.Handle("POST /admin/settings/mgm-rewards/{id}", protected(h.handleUpdateMGMRewardConfig))
	mux.Handle("POST /admin/settings/mgm-rewards/{id}/toggle-active", protected(h.handleToggleMGMRewardConfigActive))

	// Settings - Referrers
	mux.Handle("GET /admin/settings/referrers", protected(h.handleReferrersSettings))
	mux.Handle("POST /admin/settings/referrers", protected(h.handleCreateReferrer))
	mux.Handle("POST /admin/settings/referrers/{id}", protected(h.handleUpdateReferrer))
	mux.Handle("POST /admin/settings/referrers/{id}/toggle-active", protected(h.handleToggleReferrerActive))

	// Settings - Assignment Algorithm
	mux.Handle("GET /admin/settings/assignment", protected(h.handleAssignmentSettings))
	mux.Handle("POST /admin/settings/assignment/{id}/activate", protected(h.handleActivateAlgorithm))

	// Settings - Document Types
	mux.Handle("GET /admin/settings/document-types", protected(h.handleDocumentTypesSettings))
	mux.Handle("POST /admin/settings/document-types", protected(h.handleCreateDocumentType))
	mux.Handle("POST /admin/settings/document-types/{id}", protected(h.handleUpdateDocumentType))
	mux.Handle("POST /admin/settings/document-types/{id}/toggle-active", protected(h.handleToggleDocumentTypeActive))
}
