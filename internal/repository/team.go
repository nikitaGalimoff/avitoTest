package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"avitotest/internal/domain"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &teamRepository{db: db}
}
func (r *teamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	query := `SELECT EXiSTS(
		SELECT 1
		FROM USERS
		WHERE team_name = $1
		) AS team_exists`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return true, fmt.Errorf("failed to get team: %w", err)
	}
	return exists, nil

}
func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if len(team.Members) == 0 {
		return nil
	}
	baseQuery := `
		INSERT INTO users (user_id, username, team_name, is_active) 
		VALUES `

	valuePlaceholders := []string{}
	params := []interface{}{}
	paramCounter := 1

	for _, member := range team.Members {
		valuePlaceholders = append(valuePlaceholders,
			fmt.Sprintf("($%d,$%d,$%d,$%d)",
				paramCounter, paramCounter+1, paramCounter+2, paramCounter+3))

		params = append(params, member.UserID, member.Username, team.TeamName, member.IsActive)
		paramCounter += 4
	}

	finalQuery := baseQuery + strings.Join(valuePlaceholders, ",") + `
		ON CONFLICT(user_id) DO UPDATE SET 
			team_name = EXCLUDED.team_name`

	_, err := r.db.ExecContext(ctx, finalQuery, params...)
	if err != nil {
		return fmt.Errorf("failed to update team name for users: %w", err)
	}

	return nil
}

func (r *teamRepository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `SELECT user_id, username, is_active FROM users WHERE team_name = $1`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var member domain.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate team members: %w", err)
	}

	if len(members) == 0 {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "team not found")
	}

	return &domain.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}
