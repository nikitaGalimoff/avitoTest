package usecase

import (
	"context"

	"avitotest/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if err := uc.userRepo.SetIsActive(ctx, userID, isActive); err != nil {
		return nil, err
	}

	return uc.userRepo.GetByID(ctx, userID)
}
