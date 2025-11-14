package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"avitotest/internal/domain"
)

type teamRepository struct {
	db *sql.DB
}

const timeout = 5 * time.Second

// NewTeamRepository создает новый репозиторий команд
func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {

	if len(team.Members) == 0 {
		return nil // Нет пользователей для обновления
	}

	args := []interface{}{team.TeamName}
	placeholders := make([]string, len(team.Members))

	for i, member := range team.Members {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, member.UserID)
	}
	//ПЕРЕПИСАТЬ НА UPSERt
	query := fmt.Sprintf(
		`UPDATE users SET team_name = $1 WHERE user_id IN (%s)`,
		strings.Join(placeholders, ", "),
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
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

func (r *teamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE team_name = $1 )`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return exists, nil
}
