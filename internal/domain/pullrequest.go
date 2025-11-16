package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name" db:"pull_request_name"`
	AuthorID          string     `json:"author_id" db:"author_id"`
	Status            PRStatus   `json:"status" db:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers" db:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty" db:"created_at"`
	MergedAt          *time.Time `json:"mergedAt,omitempty" db:"merged_at"`
}

type PullRequestShort struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PRStatus `json:"status"`
}
