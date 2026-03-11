package handler

import (
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleReferralClaims(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Klaim Referral")

	claims, err := model.ListUnverifiedReferralClaims(ctx)
	if err != nil {
		slog.Error("Failed to list referral claims", "error", err)
		http.Error(w, "Failed to load referral claims", http.StatusInternalServerError)
		return
	}

	claimItems := make([]admin.ReferralClaimItem, len(claims))
	for i, c := range claims {
		claimItems[i] = admin.ReferralClaimItem{
			CandidateID:   c.CandidateID,
			CandidateName: c.CandidateName,
			ProdiName:     c.ProdiName,
			SourceType:    c.SourceType,
			SourceDetail:  c.SourceDetail,
			Status:        c.Status,
			CreatedAt:     c.CreatedAt.Format("2 Jan 2006"),
		}
	}

	referrers, err := model.ListReferrers(ctx, "")
	if err != nil {
		slog.Error("Failed to list referrers", "error", err)
		http.Error(w, "Failed to load referrers", http.StatusInternalServerError)
		return
	}

	referrerOptions := make([]admin.ReferrerOption, len(referrers))
	for i, r := range referrers {
		referrerOptions[i] = admin.ReferrerOption{
			ID:   r.ID,
			Name: r.Name,
			Type: r.Type,
		}
		if r.Institution != nil {
			referrerOptions[i].Institution = *r.Institution
		}
		if r.Code != nil {
			referrerOptions[i].Code = *r.Code
		}
	}

	admin.ReferralClaims(data, claimItems, referrerOptions).Render(ctx, w)
}

func (h *AdminHandler) handleLinkReferralClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseForm()

	candidateID := r.FormValue("candidate_id")
	referrerID := r.FormValue("referrer_id")
	mgmCode := r.FormValue("mgm_code")

	if candidateID == "" {
		http.Error(w, "Missing candidate_id", http.StatusBadRequest)
		return
	}

	if referrerID == "" && mgmCode == "" {
		http.Error(w, "Either referrer_id or mgm_code is required", http.StatusBadRequest)
		return
	}

	if referrerID != "" {
		err := model.LinkCandidateToReferrer(ctx, candidateID, referrerID)
		if err != nil {
			slog.Error("Failed to link candidate to referrer", "error", err)
			http.Error(w, "Failed to link referrer", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", "referralClaimLinked")
		w.Header().Set("HX-Refresh", "true")
		return
	}

	if mgmCode != "" {
		referrerCandidate, err := model.FindCandidateByReferralCode(ctx, mgmCode)
		if err != nil {
			slog.Error("Failed to find MGM referrer", "error", err)
			http.Error(w, "Failed to find MGM referrer", http.StatusInternalServerError)
			return
		}
		if referrerCandidate == nil {
			http.Error(w, "MGM code not found", http.StatusBadRequest)
			return
		}

		err = model.LinkCandidateToMGMReferrer(ctx, candidateID, referrerCandidate.ID)
		if err != nil {
			slog.Error("Failed to link candidate to MGM referrer", "error", err)
			http.Error(w, "Failed to link MGM referrer", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", "referralClaimLinked")
		w.Header().Set("HX-Refresh", "true")
		return
	}
}

func (h *AdminHandler) handleInvalidReferralClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	candidateID := r.PathValue("id")

	err := model.ClearReferralClaim(ctx, candidateID)
	if err != nil {
		slog.Error("Failed to clear referral claim", "error", err)
		http.Error(w, "Failed to clear claim", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "referralClaimInvalid")
	w.Header().Set("HX-Refresh", "true")
}
