package domain

import "context"

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	CreateOrUpdate(ctx context.Context, user *User) error
	GetByID(ctx context.Context, userID string) (*User, error)
	GetByTeamName(ctx context.Context, teamName string) ([]*User, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) error
}

// TeamRepository определяет интерфейс для работы с командами
type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	GetByName(ctx context.Context, teamName string) (*Team, error)
	Exists(ctx context.Context, teamName string) (bool, error)
}

// PullRequestRepository определяет интерфейс для работы с Pull Request'ами
type PullRequestRepository interface {
	Create(ctx context.Context, pr *PullRequest) error
	GetByID(ctx context.Context, prID string) (*PullRequest, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]*PullRequest, error)
	Update(ctx context.Context, pr *PullRequest) error
	Exists(ctx context.Context, prID string) (bool, error)
}

