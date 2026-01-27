package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// User represents a staff user in the system
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	GoogleID     *string    `json:"google_id,omitempty"`
	Role         string     `json:"role"`
	IDSupervisor *string    `json:"id_supervisor,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// UserWithSupervisor includes supervisor name for display
type UserWithSupervisor struct {
	User
	SupervisorName *string `json:"supervisor_name,omitempty"`
}

// CreateUser creates a new user
func CreateUser(ctx context.Context, email, name, googleID, role string) (*User, error) {
	var user User

	// Encrypt fields before storing
	emailEnc, nameEnc, googleIDEnc, err := encryptUserFields(email, name, googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt user fields: %w", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, name, google_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, google_id, role, id_supervisor, is_active, last_login_at, created_at, updated_at
	`, emailEnc, nameEnc, googleIDEnc, role).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID, &user.Role,
		&user.IDSupervisor, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Decrypt fields before returning
	user.Email, user.Name, user.GoogleID, err = decryptUserFields(user.Email, user.Name, user.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt user: %w", err)
	}

	return &user, nil
}

// FindUserByEmail finds a user by email
func FindUserByEmail(ctx context.Context, email string) (*User, error) {
	// Encrypt email for search (deterministic encryption allows equality match)
	emailEnc, err := encryptEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt email for search: %w", err)
	}

	var user User
	err = pool.QueryRow(ctx, `
		SELECT id, email, name, google_id, role, id_supervisor, is_active, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`, emailEnc).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID, &user.Role,
		&user.IDSupervisor, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Decrypt fields
	user.Email, user.Name, user.GoogleID, err = decryptUserFields(user.Email, user.Name, user.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt user: %w", err)
	}

	return &user, nil
}

// FindUserByGoogleID finds a user by Google ID
func FindUserByGoogleID(ctx context.Context, googleID string) (*User, error) {
	// Encrypt googleID for search (deterministic encryption allows equality match)
	googleIDEnc, err := encryptEmail(googleID) // Using encryptEmail as it's deterministic
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt google_id for search: %w", err)
	}

	var user User
	err = pool.QueryRow(ctx, `
		SELECT id, email, name, google_id, role, id_supervisor, is_active, last_login_at, created_at, updated_at
		FROM users WHERE google_id = $1
	`, googleIDEnc).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID, &user.Role,
		&user.IDSupervisor, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by google_id: %w", err)
	}

	// Decrypt fields
	user.Email, user.Name, user.GoogleID, err = decryptUserFields(user.Email, user.Name, user.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt user: %w", err)
	}

	return &user, nil
}

// FindUserByID finds a user by ID
func FindUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := pool.QueryRow(ctx, `
		SELECT id, email, name, google_id, role, id_supervisor, is_active, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID, &user.Role,
		&user.IDSupervisor, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	// Decrypt fields
	user.Email, user.Name, user.GoogleID, err = decryptUserFields(user.Email, user.Name, user.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt user: %w", err)
	}

	return &user, nil
}

// ListUsers returns all users with optional filters
func ListUsers(ctx context.Context, role string, activeOnly bool) ([]UserWithSupervisor, error) {
	query := `
		SELECT u.id, u.email, u.name, u.google_id, u.role, u.id_supervisor, u.is_active,
			   u.last_login_at, u.created_at, u.updated_at, s.name as supervisor_name
		FROM users u
		LEFT JOIN users s ON u.id_supervisor = s.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if role != "" {
		argCount++
		query += fmt.Sprintf(" AND u.role = $%d", argCount)
		args = append(args, role)
	}

	if activeOnly {
		query += " AND u.is_active = true"
	}

	query += " ORDER BY u.created_at DESC"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []UserWithSupervisor
	for rows.Next() {
		var u UserWithSupervisor
		err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.GoogleID, &u.Role, &u.IDSupervisor,
			&u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.SupervisorName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		// Decrypt fields
		u.Email, u.Name, u.GoogleID, _ = decryptUserFields(u.Email, u.Name, u.GoogleID)
		u.SupervisorName, _ = decryptNullableP(u.SupervisorName)

		users = append(users, u)
	}
	return users, nil
}

// ListSupervisors returns users who can be supervisors (admin and supervisor roles)
func ListSupervisors(ctx context.Context) ([]User, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, email, name, google_id, role, id_supervisor, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE role IN ('admin', 'supervisor') AND is_active = true
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list supervisors: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.GoogleID, &u.Role, &u.IDSupervisor,
			&u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supervisor: %w", err)
		}

		// Decrypt fields
		u.Email, u.Name, u.GoogleID, _ = decryptUserFields(u.Email, u.Name, u.GoogleID)

		users = append(users, u)
	}
	return users, nil
}

// UpdateUserRole updates a user's role
func UpdateUserRole(ctx context.Context, id, role string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET role = $1 WHERE id = $2`, role, id)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}

// UpdateUserSupervisor updates a user's supervisor
func UpdateUserSupervisor(ctx context.Context, id string, supervisorID *string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET id_supervisor = $1 WHERE id = $2`, supervisorID, id)
	if err != nil {
		return fmt.Errorf("failed to update user supervisor: %w", err)
	}
	return nil
}

// ToggleUserActive toggles user active status
func ToggleUserActive(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET is_active = NOT is_active WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to toggle user active: %w", err)
	}
	return nil
}

// UpdateLastLogin updates user's last login timestamp
func UpdateLastLogin(ctx context.Context, id string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// CountUsersByRole returns count of users by role
func CountUsersByRole(ctx context.Context) (map[string]int, error) {
	rows, err := pool.Query(ctx, `
		SELECT role, COUNT(*) FROM users WHERE is_active = true GROUP BY role
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var role string
		var count int
		if err := rows.Scan(&role, &count); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[role] = count
	}
	return counts, nil
}
