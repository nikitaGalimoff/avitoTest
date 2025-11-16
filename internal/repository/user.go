package repository

import (
	"context"
	"database/sql"
	"fmt"

	"avitotest/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateOrUpdate(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) 
		DO UPDATE SET username = $2, team_name = $3, is_active = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByTeamName(ctx context.Context, teamName string) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `SELECT user_id, username, team_name, is_active FROM users WHERE team_name = $1`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by team: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}

func (r *userRepository) SetIsActive(ctx context.Context, userID string, isActive bool) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `UPDATE users SET is_active = $1 WHERE user_id = $2`

	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return fmt.Errorf("failed to update user activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}

	return nil
}
