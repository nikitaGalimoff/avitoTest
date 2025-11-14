package usecase

import (
	"context"

	"avitotest/internal/domain"
)

// UserUseCase определяет бизнес-логику для работы с пользователями
type UserUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase создает новый use case для пользователей
func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// SetIsActive устанавливает флаг активности пользователя
func (uc *UserUseCase) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if err := uc.userRepo.SetIsActive(ctx, userID, isActive); err != nil {
		return nil, err
	}

	return uc.userRepo.GetByID(ctx, userID)
}

// GetReviewPullRequests получает список PR, где пользователь назначен ревьювером
// Примечание: этот метод не используется, вместо него используется PullRequestUseCase.GetPullRequestsByReviewer
func (uc *UserUseCase) GetReviewPullRequests(ctx context.Context, userID string) ([]*domain.PullRequestShort, error) {
	// Проверяем, существует ли пользователь
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Этот метод не реализован, так как используется PullRequestUseCase.GetPullRequestsByReviewer
	return nil, nil
}
