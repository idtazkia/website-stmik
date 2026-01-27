package handler

import (
	"net/http"

	"github.com/idtazkia/stmik-admission-api/internal/auth"
	"github.com/idtazkia/stmik-admission-api/internal/integration"
)

// AdminHandler handles all admin routes
type AdminHandler struct {
	sessionMgr *auth.SessionManager
	resend     *integration.ResendClient
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(sessionMgr *auth.SessionManager, resend *integration.ResendClient) *AdminHandler {
	return &AdminHandler{sessionMgr: sessionMgr, resend: resend}
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
	mux.Handle("GET /admin/my-dashboard", protected(h.handleConsultantDashboardReal))

	// Supervisor dashboard
	mux.Handle("GET /admin/supervisor-dashboard", protected(h.handleSupervisorDashboard))

	// Candidates
	mux.Handle("GET /admin/candidates", protected(h.handleCandidates))
	mux.Handle("GET /admin/candidates/{id}", protected(h.handleCandidateDetail))
	mux.Handle("GET /admin/candidates/{id}/interaction", protected(h.handleInteractionForm))
	mux.Handle("POST /admin/candidates/{id}/interaction", protected(h.handleCreateInteraction))
	mux.Handle("POST /admin/candidates/{id}/reassign", protected(h.handleReassignCandidate))
	mux.Handle("GET /admin/reassign-modal", protected(h.handleGetConsultantsForReassign))
	mux.Handle("GET /admin/candidates/{id}/lost-modal", protected(h.handleGetLostModal))
	mux.Handle("POST /admin/candidates/{id}/lost", protected(h.handleMarkLost))
	mux.Handle("GET /admin/candidates/{id}/commitment-modal", protected(h.handleGetCommitmentModal))
	mux.Handle("POST /admin/candidates/{id}/commit", protected(h.handleCommitCandidate))
	mux.Handle("GET /admin/candidates/{id}/enrollment-modal", protected(h.handleGetEnrollmentModal))
	mux.Handle("POST /admin/candidates/{id}/enroll", protected(h.handleEnrollCandidate))

	// Interactions - Supervisor Suggestions
	mux.Handle("POST /admin/interactions/{id}/suggestion", protected(h.handleAddSuggestion))
	mux.Handle("POST /admin/interactions/{id}/mark-read", protected(h.handleMarkSuggestionRead))

	// Documents
	mux.Handle("GET /admin/documents", protected(h.handleDocumentReview))
	mux.Handle("POST /admin/documents/approve", protected(h.handleApproveDocument))
	mux.Handle("POST /admin/documents/reject", protected(h.handleRejectDocument))

	// Marketing
	mux.Handle("GET /admin/campaigns", protected(h.handleCampaigns))
	mux.Handle("GET /admin/referrers", protected(h.handleReferrers))
	mux.Handle("GET /admin/referral-claims", protected(h.handleReferralClaims))
	mux.Handle("POST /admin/referral-claims/link", protected(h.handleLinkReferralClaim))
	mux.Handle("POST /admin/referral-claims/{id}/invalid", protected(h.handleInvalidReferralClaim))
	mux.Handle("GET /admin/commissions", protected(h.handleCommissionsReal))
	mux.Handle("GET /admin/commissions/export", protected(h.handleExportCommissions))
	mux.Handle("POST /admin/commissions/{id}/approve", protected(h.handleApproveCommission))
	mux.Handle("POST /admin/commissions/{id}/paid", protected(h.handleMarkCommissionPaid))
	mux.Handle("POST /admin/commissions/batch-approve", protected(h.handleBatchApproveCommissions))
	mux.Handle("POST /admin/commissions/batch-paid", protected(h.handleBatchMarkCommissionsPaid))

	// Reports
	mux.Handle("GET /admin/reports/funnel", protected(h.handleFunnelReport))
	mux.Handle("GET /admin/reports/consultants", protected(h.handleConsultantsReport))
	mux.Handle("GET /admin/reports/campaigns", protected(h.handleCampaignsReport))
	mux.Handle("GET /admin/reports/referrers", protected(h.handleReferrersReport))

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

	// Settings - Lost Reasons
	mux.Handle("GET /admin/settings/lost-reasons", protected(h.handleLostReasonsSettings))
	mux.Handle("POST /admin/settings/lost-reasons", protected(h.handleCreateLostReason))
	mux.Handle("POST /admin/settings/lost-reasons/{id}", protected(h.handleUpdateLostReason))
	mux.Handle("POST /admin/settings/lost-reasons/{id}/toggle-active", protected(h.handleToggleLostReasonActive))

	// Announcements
	mux.Handle("GET /admin/announcements", protected(h.handleAnnouncementsSettings))
	mux.Handle("POST /admin/announcements", protected(h.handleCreateAnnouncement))
	mux.Handle("GET /admin/announcements/{id}/edit", protected(h.handleEditAnnouncementForm))
	mux.Handle("PUT /admin/announcements/{id}", protected(h.handleUpdateAnnouncement))
	mux.Handle("POST /admin/announcements/{id}/publish", protected(h.handlePublishAnnouncement))
	mux.Handle("POST /admin/announcements/{id}/unpublish", protected(h.handleUnpublishAnnouncement))
	mux.Handle("DELETE /admin/announcements/{id}", protected(h.handleDeleteAnnouncement))
}
