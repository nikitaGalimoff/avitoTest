package usecase

import (
	"context"

	"avitotest/internal/domain"
)

type TeamUseCase struct {
	teamRepo domain.TeamRepository
	userRepo domain.UserRepository
}

func NewTeamUseCase(teamRepo domain.TeamRepository, userRepo domain.UserRepository) *TeamUseCase {
	return &TeamUseCase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (uc *TeamUseCase) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	exists, err := uc.teamRepo.Exists(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewDomainError(domain.ErrorCodeTeamExists, "")
	}

	err = uc.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *TeamUseCase) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	return uc.teamRepo.GetByName(ctx, teamName)
}
