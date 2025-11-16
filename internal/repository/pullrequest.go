package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"avitotest/internal/domain"
)

type pullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) domain.PullRequestRepository {
	return &pullRequestRepository{db: db}
}

func (r *pullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return fmt.Errorf("failed to marshal reviewers: %w", err)
	}

	query := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		string(pr.Status),
		reviewersJSON,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	pr.CreatedAt = &now
	return nil
}

func (r *pullRequestRepository) GetByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, 
		       assigned_reviewers, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr domain.PullRequest
	var statusStr string
	var reviewersJSON []byte
	var createdAt, mergedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&statusStr,
		&reviewersJSON,
		&createdAt,
		&mergedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "pull request not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	pr.Status = domain.PRStatus(statusStr)
	if err := json.Unmarshal(reviewersJSON, &pr.AssignedReviewers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reviewers: %w", err)
	}

	if createdAt.Valid {
		pr.CreatedAt = &createdAt.Time
	}
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	return &pr, nil
}

func (r *pullRequestRepository) GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, 
		       assigned_reviewers, created_at, merged_at
		FROM pull_requests
		WHERE assigned_reviewers::jsonb ? $1
	`

	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var statusStr string
		var reviewersJSON []byte
		var createdAt, mergedAt sql.NullTime

		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&statusStr,
			&reviewersJSON,
			&createdAt,
			&mergedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan pull request: %w", err)
		}

		pr.Status = domain.PRStatus(statusStr)
		if err := json.Unmarshal(reviewersJSON, &pr.AssignedReviewers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal reviewers: %w", err)
		}

		prs = append(prs, &pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate pull requests: %w", err)
	}

	return prs, nil
}

func (r *pullRequestRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return fmt.Errorf("failed to marshal reviewers: %w", err)
	}

	query := `
		UPDATE pull_requests
		SET pull_request_name = $2, author_id = $3, status = $4, 
		    assigned_reviewers = $5, merged_at = $6
		WHERE pull_request_id = $1
	`

	_, err = r.db.ExecContext(ctx, query,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		string(pr.Status),
		reviewersJSON,
		pr.MergedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}

	return nil
}

func (r *pullRequestRepository) Exists(ctx context.Context, prID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1 LIMIT 1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check pull request existence: %w", err)
	}
	return exists, nil
}
