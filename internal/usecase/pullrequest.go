package usecase

import (
	"context"
	"math/rand"
	"time"

	"avitotest/internal/domain"
)

// PullRequestUseCase определяет бизнес-логику для работы с Pull Request'ами
type PullRequestUseCase struct {
	prRepo   domain.PullRequestRepository
	userRepo domain.UserRepository
	teamRepo domain.TeamRepository
}

// NewPullRequestUseCase создает новый use case для Pull Request'ов
func NewPullRequestUseCase(
	prRepo domain.PullRequestRepository,
	userRepo domain.UserRepository,
	teamRepo domain.TeamRepository,
) *PullRequestUseCase {
	return &PullRequestUseCase{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// CreatePullRequest создает PR и автоматически назначает до 2 ревьюверов
func (uc *PullRequestUseCase) CreatePullRequest(
	ctx context.Context,
	prID, prName, authorID string,
) (*domain.PullRequest, error) {
	// Проверяем, существует ли PR
	exists, err := uc.prRepo.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewDomainError(domain.ErrorCodePRExists, "PR id already exists")
	}

	// Получаем автора
	author, err := uc.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	// Получаем всех активных пользователей команды автора
	teamUsers, err := uc.userRepo.GetByTeamName(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	// Фильтруем активных пользователей, исключая автора
	var candidates []*domain.User
	for _, user := range teamUsers {
		if user.IsActive && user.UserID != authorID {
			candidates = append(candidates, user)
		}
	}

	// Назначаем до 2 ревьюверов
	reviewers := uc.selectReviewers(candidates, 2)

	pr := &domain.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: reviewers,
	}

	if err := uc.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// MergePullRequest помечает PR как MERGED (идемпотентная операция)
func (uc *PullRequestUseCase) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := uc.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	// Если уже MERGED, просто возвращаем текущее состояние (идемпотентность)
	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	// Помечаем как MERGED
	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// ReassignReviewer переназначает ревьювера
func (uc *PullRequestUseCase) ReassignReviewer(
	ctx context.Context,
	prID, oldUserID string,
) (*domain.PullRequest, string, error) {
	// Получаем PR
	pr, err := uc.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	// Проверяем, что PR не MERGED
	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.NewDomainError(domain.ErrorCodePRMerged, "cannot reassign on merged PR")
	}

	// Проверяем, что старый ревьювер назначен
	isAssigned := false
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
	}

	// Получаем информацию о старом ревьювере
	oldReviewer, err := uc.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}

	// Получаем всех активных пользователей команды старого ревьювера
	teamUsers, err := uc.userRepo.GetByTeamName(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	// Фильтруем активных кандидатов, исключая старого ревьювера и автора
	var candidates []*domain.User
	for _, user := range teamUsers {
		if user.IsActive && user.UserID != oldUserID && user.UserID != pr.AuthorID {
			// Исключаем уже назначенных ревьюверов
			isAlreadyAssigned := false
			for _, reviewerID := range pr.AssignedReviewers {
				if reviewerID == user.UserID {
					isAlreadyAssigned = true
					break
				}
			}
			if !isAlreadyAssigned {
				candidates = append(candidates, user)
			}
		}
	}

	if len(candidates) == 0 {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNoCandidate, "no active replacement candidate in team")
	}

	// Выбираем случайного кандидата
	newReviewerID := uc.selectReviewers(candidates, 1)[0]

	// Заменяем старого ревьювера на нового
	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			pr.AssignedReviewers[i] = newReviewerID
			break
		}
	}

	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, "", err
	}

	return pr, newReviewerID, nil
}

// GetPullRequestsByReviewer получает список PR, где пользователь назначен ревьювером
func (uc *PullRequestUseCase) GetPullRequestsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequestShort, error) {

	// Получаем все PR, где пользователь является ревьювером
	prs, err := uc.prRepo.GetByReviewerID(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в короткий формат
	result := make([]*domain.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		result = append(result, &domain.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return result, nil
}

// selectReviewers выбирает случайных ревьюверов из кандидатов (до maxCount)
func (uc *PullRequestUseCase) selectReviewers(candidates []*domain.User, maxCount int) []string {
	if len(candidates) == 0 {
		return []string{}
	}

	count := maxCount
	if len(candidates) < count {
		count = len(candidates)
	}

	// Перемешиваем кандидатов
	shuffled := make([]*domain.User, len(candidates))
	copy(shuffled, candidates)
	// Используем новый источник случайных чисел (Go 1.21+)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Выбираем первых count
	reviewers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		reviewers = append(reviewers, shuffled[i].UserID)
	}

	return reviewers
}
