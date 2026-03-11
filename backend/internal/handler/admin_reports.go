package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/web/templates/admin"
)

func (h *AdminHandler) handleFunnelReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Laporan Funnel")

	funnelStats, err := model.GetFunnelStats(ctx)
	if err != nil {
		slog.Error("Failed to get funnel stats", "error", err)
		http.Error(w, "Failed to load funnel data", http.StatusInternalServerError)
		return
	}

	stages := []admin.FunnelStage{
		{
			Name:       "Registered",
			Count:      fmt.Sprintf("%d", funnelStats.Registered),
			Percentage: "100",
			Color:      "bg-gray-500",
		},
		{
			Name:       "Prospecting",
			Count:      fmt.Sprintf("%d", funnelStats.Prospecting),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Prospecting, funnelStats.Registered)),
			Color:      "bg-blue-500",
		},
		{
			Name:       "Committed",
			Count:      fmt.Sprintf("%d", funnelStats.Committed),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Committed, funnelStats.Registered)),
			Color:      "bg-yellow-500",
		},
		{
			Name:       "Enrolled",
			Count:      fmt.Sprintf("%d", funnelStats.Enrolled),
			Percentage: fmt.Sprintf("%d", calcPercentage(funnelStats.Enrolled, funnelStats.Registered)),
			Color:      "bg-green-500",
		},
	}

	regToProsp := calcConversionRate(funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled, funnelStats.Registered+funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled)
	prospToCommit := calcConversionRate(funnelStats.Committed+funnelStats.Enrolled, funnelStats.Prospecting+funnelStats.Committed+funnelStats.Enrolled)
	commitToEnroll := calcConversionRate(funnelStats.Enrolled, funnelStats.Committed+funnelStats.Enrolled)

	conversions := []admin.FunnelConversion{
		{From: "Registered", To: "Prospecting", Rate: fmt.Sprintf("%.1f%%", regToProsp), Change: "-", IsPositive: true},
		{From: "Prospecting", To: "Committed", Rate: fmt.Sprintf("%.1f%%", prospToCommit), Change: "-", IsPositive: true},
		{From: "Committed", To: "Enrolled", Rate: fmt.Sprintf("%.1f%%", commitToEnroll), Change: "-", IsPositive: true},
	}

	admin.ReportFunnel(data, stages, conversions).Render(ctx, w)
}

func calcPercentage(part, total int) int {
	if total == 0 {
		return 0
	}
	return int(float64(part) / float64(total) * 100)
}

func calcConversionRate(converted, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(converted) / float64(total) * 100
}

func (h *AdminHandler) handleConsultantsReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Kinerja Konsultan")

	performanceData, err := model.GetConsultantPerformance(ctx, nil, nil)
	if err != nil {
		slog.Error("Failed to get consultant performance", "error", err)
		http.Error(w, "Failed to load consultant data", http.StatusInternalServerError)
		return
	}

	consultants := make([]admin.ConsultantPerformance, len(performanceData))
	var totalLeads, totalInteractions, totalCommits, totalEnrollments int

	for i, cp := range performanceData {
		conversionRate := float64(0)
		if cp.TotalLeads > 0 {
			conversionRate = float64(cp.Enrollments) / float64(cp.TotalLeads) * 100
		}

		interactionsPerCandidate := float64(0)
		if cp.TotalLeads > 0 {
			interactionsPerCandidate = float64(cp.Interactions) / float64(cp.TotalLeads)
		}

		consultants[i] = admin.ConsultantPerformance{
			Rank:                     fmt.Sprintf("%d", i+1),
			Name:                     cp.ConsultantName,
			Email:                    cp.ConsultantEmail,
			SupervisorName:           cp.SupervisorName,
			Leads:                    fmt.Sprintf("%d", cp.TotalLeads),
			Interactions:             fmt.Sprintf("%d", cp.Interactions),
			Commits:                  fmt.Sprintf("%d", cp.Commits),
			Enrollments:              fmt.Sprintf("%d", cp.Enrollments),
			ConversionRate:           conversionRate,
			ConversionRateStr:        fmt.Sprintf("%.1f", conversionRate),
			AvgDaysToCommit:          fmt.Sprintf("%.0f", cp.AvgDaysToCommit),
			InteractionsPerCandidate: fmt.Sprintf("%.1f", interactionsPerCandidate),
		}

		totalLeads += cp.TotalLeads
		totalInteractions += cp.Interactions
		totalCommits += cp.Commits
		totalEnrollments += cp.Enrollments
	}

	filter := admin.ReportFilter{
		Period:    "all_time",
		StartDate: "",
		EndDate:   "",
	}

	summary := admin.ReportSummary{
		TotalLeads:        fmt.Sprintf("%d", totalLeads),
		TotalInteractions: fmt.Sprintf("%d", totalInteractions),
		TotalCommits:      fmt.Sprintf("%d", totalCommits),
		TotalEnrollments:  fmt.Sprintf("%d", totalEnrollments),
	}

	admin.ConsultantReport(data, filter, consultants, summary).Render(ctx, w)
}

func (h *AdminHandler) handleCampaignsReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "ROI Kampanye")

	campaignStats, err := model.GetCampaignStats(ctx)
	if err != nil {
		slog.Error("Failed to get campaign stats", "error", err)
		http.Error(w, "Failed to load campaign data", http.StatusInternalServerError)
		return
	}

	campaigns := make([]admin.CampaignReportItem, len(campaignStats))
	var totalLeads, totalEnrolled int
	var bestCampaign string
	var bestConversion float64

	for i, cs := range campaignStats {
		totalCandidates := cs.Registered + cs.Prospecting + cs.Committed + cs.Enrolled
		conversion := float64(0)
		if totalCandidates > 0 {
			conversion = float64(cs.Enrolled) / float64(totalCandidates) * 100
		}

		campaigns[i] = admin.CampaignReportItem{
			Name:        cs.CampaignName,
			Type:        cs.CampaignType,
			Channel:     cs.Channel,
			Leads:       fmt.Sprintf("%d", totalCandidates),
			Prospecting: fmt.Sprintf("%d", cs.Prospecting),
			Committed:   fmt.Sprintf("%d", cs.Committed),
			Enrolled:    fmt.Sprintf("%d", cs.Enrolled),
			Conversion:  fmt.Sprintf("%.1f%%", conversion),
			Cost:        "-", // Cost data not tracked in campaigns table
			CostPerLead: "-", // Cost data not tracked
		}

		totalLeads += totalCandidates
		totalEnrolled += cs.Enrolled

		if conversion > bestConversion && totalCandidates > 0 {
			bestConversion = conversion
			bestCampaign = cs.CampaignName
		}
	}

	avgConversion := float64(0)
	if totalLeads > 0 {
		avgConversion = float64(totalEnrolled) / float64(totalLeads) * 100
	}

	summary := admin.CampaignReportSummary{
		TotalLeads:     fmt.Sprintf("%d", totalLeads),
		TotalEnrolled:  fmt.Sprintf("%d", totalEnrolled),
		AvgConversion:  fmt.Sprintf("%.1f%%", avgConversion),
		TotalCost:      "-",
		AvgCostPerLead: "-",
		BestCampaign:   bestCampaign,
	}

	admin.ReportCampaigns(data, campaigns, summary).Render(ctx, w)
}

func (h *AdminHandler) handleReferrersReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := NewPageDataWithUser(ctx, "Leaderboard Referrer")

	referrerStats, err := model.GetReferrerStats(ctx)
	if err != nil {
		slog.Error("Failed to get referrer stats", "error", err)
		http.Error(w, "Failed to load referrer data", http.StatusInternalServerError)
		return
	}

	referrers := make([]admin.ReferrerReportItem, len(referrerStats))
	var totalReferrals, totalEnrolled int
	var totalCommission int64
	var bestReferrer string
	var bestEnrolled int

	for i, rs := range referrerStats {
		conversion := float64(0)
		if rs.TotalReferrals > 0 {
			conversion = float64(rs.Enrolled) / float64(rs.TotalReferrals) * 100
		}

		referrers[i] = admin.ReferrerReportItem{
			Rank:           fmt.Sprintf("%d", i+1),
			Name:           rs.ReferrerName,
			Type:           rs.ReferrerType,
			Institution:    rs.Institution,
			TotalReferrals: fmt.Sprintf("%d", rs.TotalReferrals),
			Enrolled:       fmt.Sprintf("%d", rs.Enrolled),
			Pending:        fmt.Sprintf("%d", rs.Pending),
			Conversion:     fmt.Sprintf("%.1f%%", conversion),
			CommissionPaid: formatRupiah(rs.CommissionPaid),
		}

		totalReferrals += rs.TotalReferrals
		totalEnrolled += rs.Enrolled
		totalCommission += rs.CommissionPaid

		if rs.Enrolled > bestEnrolled {
			bestEnrolled = rs.Enrolled
			bestReferrer = rs.ReferrerName
		}
	}

	summary := admin.ReferrerReportSummary{
		TotalReferrers:  fmt.Sprintf("%d", len(referrerStats)),
		TotalReferrals:  fmt.Sprintf("%d", totalReferrals),
		TotalEnrolled:   fmt.Sprintf("%d", totalEnrolled),
		TotalCommission: formatRupiah(totalCommission),
		BestReferrer:    bestReferrer,
	}

	admin.ReportReferrers(data, referrers, summary).Render(ctx, w)
}
