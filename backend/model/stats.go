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
		WHERE assigned_consultant_id = $1
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
			WHERE candidate_id = c.id
		) li ON true
		WHERE c.assigned_consultant_id = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND (li.last_contact IS NULL OR li.last_contact < NOW() - INTERVAL '7 days')
	`, consultantID).Scan(&stats.OverdueCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue count: %w", err)
	}

	// Get today's tasks (follow-ups scheduled for today)
	err = pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT i.candidate_id)
		FROM interactions i
		JOIN candidates c ON i.candidate_id = c.id
		WHERE c.assigned_consultant_id = $1
		  AND c.status NOT IN ('enrolled', 'lost')
		  AND i.next_followup_date IS NOT NULL
		  AND DATE(i.next_followup_date) = CURRENT_DATE
		  AND i.id = (
			SELECT id FROM interactions
			WHERE candidate_id = i.candidate_id
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
		JOIN candidates c ON i.candidate_id = c.id
		WHERE c.assigned_consultant_id = $1
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
		WHERE assigned_consultant_id = $1
		  AND created_at >= $2
	`, consultantID, startOfMonth).Scan(&stats.MonthlyNewLeads)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly leads: %w", err)
	}

	// Monthly interactions
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM interactions
		WHERE consultant_id = $1
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
		WHERE assigned_consultant_id = $1
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
		WHERE assigned_consultant_id = $1
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
		LEFT JOIN prodis p ON c.prodi_id = p.id
		LEFT JOIN LATERAL (
			SELECT MAX(created_at) as last_contact
			FROM interactions
			WHERE candidate_id = c.id
		) li ON true
		WHERE c.assigned_consultant_id = $1
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
		JOIN interactions i ON i.candidate_id = c.id
		LEFT JOIN prodis p ON c.prodi_id = p.id
		WHERE c.assigned_consultant_id = $1
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
		SELECT i.id, i.candidate_id, c.name as candidate_name, i.supervisor_suggestion,
		       u.id as supervisor_id, u.name as supervisor_name, i.created_at
		FROM interactions i
		JOIN candidates c ON i.candidate_id = c.id
		LEFT JOIN users u ON i.consultant_id != $1 AND i.supervisor_suggestion IS NOT NULL
		WHERE c.assigned_consultant_id = $1
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
