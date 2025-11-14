package usecase

import (
	"context"

	"avitotest/internal/domain"
)

// TeamUseCase определяет бизнес-логику для работы с командами
type TeamUseCase struct {
	teamRepo domain.TeamRepository
	userRepo domain.UserRepository
}

// NewTeamUseCase создает новый use case для команд
func NewTeamUseCase(teamRepo domain.TeamRepository, userRepo domain.UserRepository) *TeamUseCase {
	return &TeamUseCase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam создает команду с участниками
func (uc *TeamUseCase) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	// Проверяем, существует ли команда
	exists, err := uc.teamRepo.Exists(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewDomainError(domain.ErrorCodeTeamExists, "team_name already exists")
	}

	// Создаем/обновляем всех пользователей команды
	for _, member := range team.Members {
		user := &domain.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}
		if err := uc.userRepo.CreateOrUpdate(ctx, user); err != nil {
			return nil, err
		}
	}

	// Сохраняем команду (в нашей модели это просто проверка)
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

// GetTeam получает команду по имени
func (uc *TeamUseCase) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	return uc.teamRepo.GetByName(ctx, teamName)
}

