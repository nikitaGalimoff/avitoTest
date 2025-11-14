package repository

import (
	"context"
	"database/sql"
	"fmt"

	"avitotest/internal/domain"
)

type teamRepository struct {
	db *sql.DB
}

// NewTeamRepository создает новый репозиторий команд
func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {
	// Проверяем, существует ли команда
	exists, err := r.Exists(ctx, team.TeamName)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewDomainError(domain.ErrorCodeTeamExists, "team_name already exists")
	}

	// Создаем команду (в нашей модели команда - это просто группа пользователей)
	// Пользователи уже должны быть созданы через UserRepository
	// Здесь мы просто проверяем, что команда не существует
	return nil
}

func (r *teamRepository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	// Получаем всех пользователей команды
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
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE team_name = $1 LIMIT 1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return exists, nil
}

