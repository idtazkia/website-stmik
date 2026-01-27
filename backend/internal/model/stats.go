package model

import (
	"context"
	"fmt"
	"time"
)

// ConsultantStats contains statistics for a consultant's dashboard
type ConsultantStats struct {
	TotalCandidates     int
	Prospecting         int
	Committed           int
	Enrolled            int
	Lost                int
	OverdueCount        int
	TodayTasks          int
	UnreadSuggestions   int
	MonthlyNewLeads     int
	MonthlyInteractions int
	MonthlyCommits      int
	MonthlyEnrollments  int
}

// OverdueCandidate represents a candidate with overdue follow-up
type OverdueCandidate struct {
	ID          string
	Name        string
	Phone       string
	ProdiName   string
	Status      string
	LastContact time.Time
	DaysOverdue int
}

// TodayTask represents a follow-up task scheduled for today
type TodayTask struct {
	ID           string
	Name         string
	Phone        string
	ProdiName    string
	Status       string
	FollowupDate time.Time
}

// UnreadSuggestion represents a supervisor suggestion not yet read
type UnreadSuggestion struct {
	ID             string
	InteractionID  string
	CandidateID    string
	CandidateName  string
	Suggestion     string
	SupervisorID   string
	SupervisorName string
	CreatedAt      time.Time
}

// GetConsultantStats returns dashboard statistics for a consultant
func GetConsultantStats(ctx context.Context, consultantID string) (*ConsultantStats, error) {
	var stats ConsultantStats

	// Get candidate counts by status
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'prospecting') as prospecting,
			COUNT(*) FILTER (WHERE status = 'committed') as committed,
			COUNT(*) FILTER (WHERE status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE status = 'lost') as lost
		FROM candidates
		WHERE id_assigned_consultant = $1
	`, consultantID).Scan(
		&stats.TotalCandidates, &stats.Prospecting, &stats.Committed,
		&stats.Enrolled, &stats.Lost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate counts: %w", err)
	}

	// Get overdue count (candidates with last interaction > 7 days and not enrolled/lost)
	err = pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT c.id)
		FROM candidates c
		LEFT JOIN LATERAL (
			SELECT MAX(created_at) as last_contact
			FROM interactions
			WHERE id_candidate = c.id
		) li ON true
		WHERE c.id_assigned_consultant = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND (li.last_contact IS NULL OR li.last_contact < NOW() - INTERVAL '7 days')
	`, consultantID).Scan(&stats.OverdueCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue count: %w", err)
	}

	// Get today's tasks (follow-ups scheduled for today)
	err = pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT i.id_candidate)
		FROM interactions i
		JOIN candidates c ON i.id_candidate = c.id
		WHERE c.id_assigned_consultant = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND i.next_followup_date IS NOT NULL
		  AND DATE(i.next_followup_date) = CURRENT_DATE
		  AND i.id = (
			SELECT id FROM interactions
			WHERE id_candidate = i.id_candidate
			ORDER BY created_at DESC
			LIMIT 1
		  )
	`, consultantID).Scan(&stats.TodayTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to get today tasks: %w", err)
	}

	// Get unread suggestions count
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM interactions i
		JOIN candidates c ON i.id_candidate = c.id
		WHERE c.id_assigned_consultant = $1
		  AND i.supervisor_suggestion IS NOT NULL
		  AND i.suggestion_read_at IS NULL
	`, consultantID).Scan(&stats.UnreadSuggestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread suggestions count: %w", err)
	}

	// Get monthly stats (this month)
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	startOfMonth = time.Date(startOfMonth.Year(), startOfMonth.Month(), startOfMonth.Day(), 0, 0, 0, 0, time.Local)

	// Monthly new leads (candidates assigned this month)
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM candidates
		WHERE id_assigned_consultant = $1
		  AND created_at >= $2
	`, consultantID, startOfMonth).Scan(&stats.MonthlyNewLeads)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly leads: %w", err)
	}

	// Monthly interactions
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM interactions
		WHERE id_consultant = $1
		  AND created_at >= $2
	`, consultantID, startOfMonth).Scan(&stats.MonthlyInteractions)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly interactions: %w", err)
	}

	// Monthly commits (candidates who became committed this month)
	// We approximate by checking current committed candidates created this month
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM candidates
		WHERE id_assigned_consultant = $1
		  AND status IN ('committed', 'enrolled')
		  AND updated_at >= $2
	`, consultantID, startOfMonth).Scan(&stats.MonthlyCommits)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly commits: %w", err)
	}

	// Monthly enrollments
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM candidates
		WHERE id_assigned_consultant = $1
		  AND status = 'enrolled'
		  AND updated_at >= $2
	`, consultantID, startOfMonth).Scan(&stats.MonthlyEnrollments)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly enrollments: %w", err)
	}

	return &stats, nil
}

// GetOverdueCandidates returns candidates with overdue follow-ups for a consultant
func GetOverdueCandidates(ctx context.Context, consultantID string, limit int) ([]OverdueCandidate, error) {
	rows, err := pool.Query(ctx, `
		SELECT c.id, c.name, c.phone, p.name as prodi_name, c.status,
		       COALESCE(li.last_contact, c.created_at) as last_contact,
		       EXTRACT(DAY FROM NOW() - COALESCE(li.last_contact, c.created_at))::int as days_overdue
		FROM candidates c
		LEFT JOIN prodis p ON c.id_prodi = p.id
		LEFT JOIN LATERAL (
			SELECT MAX(created_at) as last_contact
			FROM interactions
			WHERE id_candidate = c.id
		) li ON true
		WHERE c.id_assigned_consultant = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND (li.last_contact IS NULL OR li.last_contact < NOW() - INTERVAL '7 days')
		ORDER BY li.last_contact ASC NULLS FIRST
		LIMIT $2
	`, consultantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue candidates: %w", err)
	}
	defer rows.Close()

	var candidates []OverdueCandidate
	for rows.Next() {
		var c OverdueCandidate
		var name, phone, prodiName *string
		err := rows.Scan(&c.ID, &name, &phone, &prodiName, &c.Status, &c.LastContact, &c.DaysOverdue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan overdue candidate: %w", err)
		}
		if name != nil {
			c.Name = *name
		}
		if phone != nil {
			c.Phone = *phone
		}
		if prodiName != nil {
			c.ProdiName = *prodiName
		}
		candidates = append(candidates, c)
	}

	return candidates, nil
}

// GetTodayTasks returns follow-up tasks scheduled for today for a consultant
func GetTodayTasks(ctx context.Context, consultantID string, limit int) ([]TodayTask, error) {
	rows, err := pool.Query(ctx, `
		SELECT DISTINCT ON (c.id) c.id, c.name, c.phone, p.name as prodi_name, c.status, i.next_followup_date
		FROM candidates c
		JOIN interactions i ON i.id_candidate = c.id
		LEFT JOIN prodis p ON c.id_prodi = p.id
		WHERE c.id_assigned_consultant = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND i.next_followup_date IS NOT NULL
		  AND DATE(i.next_followup_date) = CURRENT_DATE
		ORDER BY c.id, i.created_at DESC
		LIMIT $2
	`, consultantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get today tasks: %w", err)
	}
	defer rows.Close()

	var tasks []TodayTask
	for rows.Next() {
		var t TodayTask
		var name, phone, prodiName *string
		err := rows.Scan(&t.ID, &name, &phone, &prodiName, &t.Status, &t.FollowupDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan today task: %w", err)
		}
		if name != nil {
			t.Name = *name
		}
		if phone != nil {
			t.Phone = *phone
		}
		if prodiName != nil {
			t.ProdiName = *prodiName
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetUnreadSuggestions returns unread supervisor suggestions for a consultant
func GetUnreadSuggestions(ctx context.Context, consultantID string, limit int) ([]UnreadSuggestion, error) {
	rows, err := pool.Query(ctx, `
		SELECT i.id, i.id_candidate, c.name as candidate_name, i.supervisor_suggestion,
		       u.id as supervisor_id, u.name as supervisor_name, i.created_at
		FROM interactions i
		JOIN candidates c ON i.id_candidate = c.id
		LEFT JOIN users u ON i.id_consultant != $1 AND i.supervisor_suggestion IS NOT NULL
		WHERE c.id_assigned_consultant = $1
		  AND i.supervisor_suggestion IS NOT NULL
		  AND i.suggestion_read_at IS NULL
		ORDER BY i.created_at DESC
		LIMIT $2
	`, consultantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread suggestions: %w", err)
	}
	defer rows.Close()

	var suggestions []UnreadSuggestion
	for rows.Next() {
		var s UnreadSuggestion
		var candidateName, supervisorName *string
		var supervisorID *string
		err := rows.Scan(&s.InteractionID, &s.CandidateID, &candidateName, &s.Suggestion,
			&supervisorID, &supervisorName, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan suggestion: %w", err)
		}
		s.ID = s.InteractionID // Use interaction ID as the suggestion ID
		if candidateName != nil {
			s.CandidateName = *candidateName
		}
		if supervisorID != nil {
			s.SupervisorID = *supervisorID
		}
		if supervisorName != nil {
			s.SupervisorName = *supervisorName
		} else {
			s.SupervisorName = "Supervisor"
		}
		suggestions = append(suggestions, s)
	}

	return suggestions, nil
}

// FunnelStats represents the overall funnel data
type FunnelStats struct {
	Registered  int
	Prospecting int
	Committed   int
	Enrolled    int
	Lost        int
}

// FunnelByProdi represents funnel data for a specific prodi
type FunnelByProdi struct {
	ProdiID     string
	ProdiName   string
	Registered  int
	Prospecting int
	Committed   int
	Enrolled    int
	Lost        int
}

// GetFunnelStats returns the overall funnel statistics
func GetFunnelStats(ctx context.Context) (*FunnelStats, error) {
	var stats FunnelStats
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'registered') as registered,
			COUNT(*) FILTER (WHERE status = 'prospecting') as prospecting,
			COUNT(*) FILTER (WHERE status = 'committed') as committed,
			COUNT(*) FILTER (WHERE status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE status = 'lost') as lost
		FROM candidates
	`).Scan(&stats.Registered, &stats.Prospecting, &stats.Committed, &stats.Enrolled, &stats.Lost)
	if err != nil {
		return nil, fmt.Errorf("failed to get funnel stats: %w", err)
	}
	return &stats, nil
}

// GetFunnelByProdi returns funnel data grouped by prodi
func GetFunnelByProdi(ctx context.Context) ([]FunnelByProdi, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			p.id, p.name,
			COUNT(*) FILTER (WHERE c.status = 'registered') as registered,
			COUNT(*) FILTER (WHERE c.status = 'prospecting') as prospecting,
			COUNT(*) FILTER (WHERE c.status = 'committed') as committed,
			COUNT(*) FILTER (WHERE c.status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE c.status = 'lost') as lost
		FROM prodis p
		LEFT JOIN candidates c ON c.id_prodi = p.id
		WHERE p.is_active = true
		GROUP BY p.id, p.name
		ORDER BY p.name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get funnel by prodi: %w", err)
	}
	defer rows.Close()

	var result []FunnelByProdi
	for rows.Next() {
		var f FunnelByProdi
		err := rows.Scan(&f.ProdiID, &f.ProdiName, &f.Registered, &f.Prospecting, &f.Committed, &f.Enrolled, &f.Lost)
		if err != nil {
			return nil, fmt.Errorf("failed to scan funnel by prodi: %w", err)
		}
		result = append(result, f)
	}
	return result, nil
}

// CampaignStats represents campaign performance data
type CampaignStats struct {
	CampaignID   string
	CampaignName string
	CampaignType string
	Channel      string
	Registered   int
	Prospecting  int
	Committed    int
	Enrolled     int
	Lost         int
}

// GetCampaignStats returns performance statistics for each campaign
func GetCampaignStats(ctx context.Context) ([]CampaignStats, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			ca.id, ca.name, ca.type, COALESCE(ca.channel, ''),
			COALESCE((SELECT COUNT(*) FROM candidates WHERE id_campaign = ca.id AND status = 'registered'), 0),
			COALESCE((SELECT COUNT(*) FROM candidates WHERE id_campaign = ca.id AND status = 'prospecting'), 0),
			COALESCE((SELECT COUNT(*) FROM candidates WHERE id_campaign = ca.id AND status = 'committed'), 0),
			COALESCE((SELECT COUNT(*) FROM candidates WHERE id_campaign = ca.id AND status = 'enrolled'), 0),
			COALESCE((SELECT COUNT(*) FROM candidates WHERE id_campaign = ca.id AND status = 'lost'), 0)
		FROM campaigns ca
		WHERE ca.is_active = true
		ORDER BY ca.name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign stats: %w", err)
	}
	defer rows.Close()

	var result []CampaignStats
	for rows.Next() {
		var cs CampaignStats
		err := rows.Scan(&cs.CampaignID, &cs.CampaignName, &cs.CampaignType, &cs.Channel,
			&cs.Registered, &cs.Prospecting, &cs.Committed, &cs.Enrolled, &cs.Lost)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign stats: %w", err)
		}
		result = append(result, cs)
	}
	return result, nil
}

// ConsultantPerformanceData represents consultant performance for the leaderboard
type ConsultantPerformanceData struct {
	ConsultantID    string
	ConsultantName  string
	ConsultantEmail string
	SupervisorName  string
	TotalLeads      int
	Interactions    int
	Commits         int
	Enrollments     int
	AvgDaysToCommit float64
}

// GetConsultantPerformance returns performance data for all consultants
func GetConsultantPerformance(ctx context.Context, startDate, endDate *time.Time) ([]ConsultantPerformanceData, error) {
	// Build date filter
	dateFilter := ""
	args := []interface{}{}
	if startDate != nil && endDate != nil {
		dateFilter = "AND c.created_at BETWEEN $1 AND $2"
		args = append(args, *startDate, *endDate)
	}

	query := fmt.Sprintf(`
		SELECT
			u.id, u.name, u.email,
			COALESCE(sup.name, '') as supervisor_name,
			COUNT(DISTINCT c.id) as total_leads,
			COALESCE((SELECT COUNT(*) FROM interactions i WHERE i.id_consultant = u.id %s), 0) as interactions,
			COUNT(DISTINCT c.id) FILTER (WHERE c.status IN ('committed', 'enrolled')) as commits,
			COUNT(DISTINCT c.id) FILTER (WHERE c.status = 'enrolled') as enrollments,
			COALESCE(AVG(
				CASE WHEN c.status IN ('committed', 'enrolled')
				THEN EXTRACT(DAY FROM COALESCE(c.updated_at, NOW()) - c.created_at)
				END
			), 0) as avg_days_to_commit
		FROM users u
		LEFT JOIN users sup ON u.id_supervisor = sup.id
		LEFT JOIN candidates c ON c.id_assigned_consultant = u.id %s
		WHERE u.role = 'consultant' AND u.is_active = true
		GROUP BY u.id, u.name, u.email, sup.name
		ORDER BY COUNT(DISTINCT c.id) FILTER (WHERE c.status = 'enrolled') DESC, COUNT(DISTINCT c.id) DESC
	`, func() string {
		if startDate != nil && endDate != nil {
			return "AND i.created_at BETWEEN $1 AND $2"
		}
		return ""
	}(), dateFilter)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultant performance: %w", err)
	}
	defer rows.Close()

	var result []ConsultantPerformanceData
	for rows.Next() {
		var cp ConsultantPerformanceData
		var name, email *string
		err := rows.Scan(&cp.ConsultantID, &name, &email, &cp.SupervisorName,
			&cp.TotalLeads, &cp.Interactions, &cp.Commits, &cp.Enrollments, &cp.AvgDaysToCommit)
		if err != nil {
			return nil, fmt.Errorf("failed to scan consultant performance: %w", err)
		}
		if name != nil {
			cp.ConsultantName = *name
		}
		if email != nil {
			cp.ConsultantEmail = *email
		}
		result = append(result, cp)
	}
	return result, nil
}

// ReferrerPerformanceData represents referrer performance for the leaderboard
type ReferrerPerformanceData struct {
	ReferrerID     string
	ReferrerName   string
	ReferrerType   string
	Institution    string
	TotalReferrals int
	Enrolled       int
	Pending        int
	CommissionPaid int64
}

// GetReferrerStats returns performance data for all referrers
func GetReferrerStats(ctx context.Context) ([]ReferrerPerformanceData, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			r.id, r.name, r.type, COALESCE(r.institution, ''),
			COUNT(c.id) as total_referrals,
			COUNT(c.id) FILTER (WHERE c.status = 'enrolled') as enrolled,
			COUNT(c.id) FILTER (WHERE c.status IN ('registered', 'prospecting', 'committed')) as pending,
			COALESCE(SUM(cl.amount) FILTER (WHERE cl.status = 'paid'), 0) as commission_paid
		FROM referrers r
		LEFT JOIN candidates c ON c.id_referrer = r.id
		LEFT JOIN commission_ledger cl ON cl.id_referrer = r.id
		WHERE r.is_active = true
		GROUP BY r.id, r.name, r.type, r.institution
		ORDER BY COUNT(c.id) FILTER (WHERE c.status = 'enrolled') DESC, COUNT(c.id) DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get referrer stats: %w", err)
	}
	defer rows.Close()

	var result []ReferrerPerformanceData
	for rows.Next() {
		var rp ReferrerPerformanceData
		err := rows.Scan(&rp.ReferrerID, &rp.ReferrerName, &rp.ReferrerType, &rp.Institution,
			&rp.TotalReferrals, &rp.Enrolled, &rp.Pending, &rp.CommissionPaid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referrer stats: %w", err)
		}
		result = append(result, rp)
	}
	return result, nil
}

// SupervisorDashboardStats for supervisor team view
type SupervisorDashboardStats struct {
	TeamRegistered   int
	TeamProspecting  int
	TeamCommitted    int
	TeamEnrolled     int
	TeamLost         int
	StuckCandidates  int // candidates with no interaction > 7 days
	TodayFollowups   int
	MonthlyNewLeads  int
	MonthlyEnrolled  int
}

// GetSupervisorDashboardStats returns dashboard statistics for a supervisor's team
func GetSupervisorDashboardStats(ctx context.Context, supervisorID string) (*SupervisorDashboardStats, error) {
	var stats SupervisorDashboardStats

	// Get team candidate counts by status
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE c.status = 'registered') as registered,
			COUNT(*) FILTER (WHERE c.status = 'prospecting') as prospecting,
			COUNT(*) FILTER (WHERE c.status = 'committed') as committed,
			COUNT(*) FILTER (WHERE c.status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE c.status = 'lost') as lost
		FROM candidates c
		WHERE c.id_assigned_consultant IN (
			SELECT id FROM users WHERE id_supervisor = $1 OR id = $1
		)
	`, supervisorID).Scan(
		&stats.TeamRegistered, &stats.TeamProspecting, &stats.TeamCommitted,
		&stats.TeamEnrolled, &stats.TeamLost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get team candidate counts: %w", err)
	}

	// Get stuck candidates (no interaction > 7 days)
	err = pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT c.id)
		FROM candidates c
		LEFT JOIN LATERAL (
			SELECT MAX(created_at) as last_contact
			FROM interactions
			WHERE id_candidate = c.id
		) li ON true
		WHERE c.id_assigned_consultant IN (
			SELECT id FROM users WHERE id_supervisor = $1 OR id = $1
		)
		AND c.status NOT IN ('enrolled', 'lost')
		AND (li.last_contact IS NULL OR li.last_contact < NOW() - INTERVAL '7 days')
	`, supervisorID).Scan(&stats.StuckCandidates)
	if err != nil {
		return nil, fmt.Errorf("failed to get stuck candidates: %w", err)
	}

	// Get today's follow-ups for the team
	err = pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT c.id)
		FROM candidates c
		JOIN interactions i ON i.id_candidate = c.id
		WHERE c.id_assigned_consultant IN (
			SELECT id FROM users WHERE id_supervisor = $1 OR id = $1
		)
		AND c.status NOT IN ('enrolled', 'lost')
		AND i.next_followup_date IS NOT NULL
		AND DATE(i.next_followup_date) = CURRENT_DATE
		AND i.id = (
			SELECT id FROM interactions
			WHERE id_candidate = i.id_candidate
			ORDER BY created_at DESC
			LIMIT 1
		)
	`, supervisorID).Scan(&stats.TodayFollowups)
	if err != nil {
		return nil, fmt.Errorf("failed to get today followups: %w", err)
	}

	// Get monthly stats
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	startOfMonth = time.Date(startOfMonth.Year(), startOfMonth.Month(), startOfMonth.Day(), 0, 0, 0, 0, time.Local)

	err = pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE c.created_at >= $2) as monthly_leads,
			COUNT(*) FILTER (WHERE c.status = 'enrolled' AND c.updated_at >= $2) as monthly_enrolled
		FROM candidates c
		WHERE c.id_assigned_consultant IN (
			SELECT id FROM users WHERE id_supervisor = $1 OR id = $1
		)
	`, supervisorID, startOfMonth).Scan(&stats.MonthlyNewLeads, &stats.MonthlyEnrolled)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	return &stats, nil
}

// TeamConsultantSummary for supervisor's team view
type TeamConsultantSummary struct {
	ConsultantID   string
	ConsultantName string
	ActiveLeads    int
	TodayTasks     int
	Overdue        int
	MonthlyCommits int
}

// GetTeamConsultants returns summary for each consultant in a supervisor's team
func GetTeamConsultants(ctx context.Context, supervisorID string) ([]TeamConsultantSummary, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			u.id, u.name,
			COUNT(c.id) FILTER (WHERE c.status NOT IN ('enrolled', 'lost')) as active_leads,
			(
				SELECT COUNT(DISTINCT c2.id)
				FROM candidates c2
				JOIN interactions i ON i.id_candidate = c2.id
				WHERE c2.id_assigned_consultant = u.id
				AND c2.status NOT IN ('enrolled', 'lost')
				AND i.next_followup_date IS NOT NULL
				AND DATE(i.next_followup_date) = CURRENT_DATE
				AND i.id = (SELECT id FROM interactions WHERE id_candidate = i.id_candidate ORDER BY created_at DESC LIMIT 1)
			) as today_tasks,
			(
				SELECT COUNT(DISTINCT c3.id)
				FROM candidates c3
				LEFT JOIN LATERAL (SELECT MAX(created_at) as lc FROM interactions WHERE id_candidate = c3.id) li ON true
				WHERE c3.id_assigned_consultant = u.id
				AND c3.status NOT IN ('enrolled', 'lost')
				AND (li.lc IS NULL OR li.lc < NOW() - INTERVAL '7 days')
			) as overdue,
			COUNT(c.id) FILTER (
				WHERE c.status IN ('committed', 'enrolled')
				AND c.updated_at >= DATE_TRUNC('month', CURRENT_DATE)
			) as monthly_commits
		FROM users u
		LEFT JOIN candidates c ON c.id_assigned_consultant = u.id
		WHERE u.id_supervisor = $1 OR u.id = $1
		GROUP BY u.id, u.name
		ORDER BY u.name
	`, supervisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team consultants: %w", err)
	}
	defer rows.Close()

	var result []TeamConsultantSummary
	for rows.Next() {
		var tc TeamConsultantSummary
		var name *string
		err := rows.Scan(&tc.ConsultantID, &name, &tc.ActiveLeads, &tc.TodayTasks, &tc.Overdue, &tc.MonthlyCommits)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team consultant: %w", err)
		}
		if name != nil {
			tc.ConsultantName = *name
		}
		result = append(result, tc)
	}
	return result, nil
}

// StuckCandidate represents a candidate stuck without recent interaction
type StuckCandidate struct {
	CandidateID    string
	CandidateName  string
	ProdiName      string
	Status         string
	ConsultantName string
	DaysStuck      int
}

// GetStuckCandidatesForTeam returns candidates with no interaction > 7 days for a supervisor's team
func GetStuckCandidatesForTeam(ctx context.Context, supervisorID string, limit int) ([]StuckCandidate, error) {
	rows, err := pool.Query(ctx, `
		SELECT c.id, c.name, COALESCE(p.name, '') as prodi_name, c.status,
		       COALESCE(u.name, '') as consultant_name,
		       EXTRACT(DAY FROM NOW() - COALESCE(li.last_contact, c.created_at))::int as days_stuck
		FROM candidates c
		LEFT JOIN prodis p ON c.id_prodi = p.id
		LEFT JOIN users u ON c.id_assigned_consultant = u.id
		LEFT JOIN LATERAL (
			SELECT MAX(created_at) as last_contact
			FROM interactions
			WHERE id_candidate = c.id
		) li ON true
		WHERE c.id_assigned_consultant IN (
			SELECT id FROM users WHERE id_supervisor = $1 OR id = $1
		)
		AND c.status NOT IN ('enrolled', 'lost')
		AND (li.last_contact IS NULL OR li.last_contact < NOW() - INTERVAL '7 days')
		ORDER BY li.last_contact ASC NULLS FIRST
		LIMIT $2
	`, supervisorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get stuck candidates for team: %w", err)
	}
	defer rows.Close()

	var result []StuckCandidate
	for rows.Next() {
		var sc StuckCandidate
		var name *string
		err := rows.Scan(&sc.CandidateID, &name, &sc.ProdiName, &sc.Status, &sc.ConsultantName, &sc.DaysStuck)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stuck candidate: %w", err)
		}
		if name != nil {
			sc.CandidateName = *name
		}
		result = append(result, sc)
	}
	return result, nil
}
