package usecase

import (
	"context"
	"math/rand"
	"time"

	"avitotest/internal/domain"
)

type PullRequestUseCase struct {
	prRepo   domain.PullRequestRepository
	userRepo domain.UserRepository
	teamRepo domain.TeamRepository
}

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

func (uc *PullRequestUseCase) CreatePullRequest(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	exists, err := uc.prRepo.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewDomainError(domain.ErrorCodePRExists, "PR id already exists")
	}

	author, err := uc.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	teamUsers, err := uc.userRepo.GetByTeamName(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []*domain.User
	for _, user := range teamUsers {
		if user.IsActive && user.UserID != authorID {
			candidates = append(candidates, user)
		}
	}

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

func (uc *PullRequestUseCase) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := uc.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (uc *PullRequestUseCase) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	pr, err := uc.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}
	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.NewDomainError(domain.ErrorCodePRMerged, "cannot reassign on merged PR")
	}

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

	oldReviewer, err := uc.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}

	teamUsers, err := uc.userRepo.GetByTeamName(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	asResMap := make(map[string]struct{})
	for _, user := range pr.AssignedReviewers {
		asResMap[user] = struct{}{}
	}

	var candidates []*domain.User
	for _, user := range teamUsers {
		if user.IsActive && user.UserID != oldUserID && user.UserID != pr.AuthorID {
			if _, ok := asResMap[user.UserID]; !ok {
				candidates = append(candidates, user)
			}
		}
	}

	if len(candidates) == 0 {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNoCandidate, "no active replacement candidate in team")
	}

	newReviewerID := uc.selectReviewers(candidates, 1)[0]

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

func (uc *PullRequestUseCase) GetPullRequestsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequestShort, error) {

	prs, err := uc.prRepo.GetByReviewerID(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

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

func (uc *PullRequestUseCase) selectReviewers(candidates []*domain.User, maxCount int) []string {
	if len(candidates) == 0 {
		return []string{}
	}

	count := maxCount
	if len(candidates) < count {
		count = len(candidates)
	}

	shuffled := make([]*domain.User, len(candidates))
	copy(shuffled, candidates)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	reviewers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		reviewers = append(reviewers, shuffled[i].UserID)
	}

	return reviewers
}
