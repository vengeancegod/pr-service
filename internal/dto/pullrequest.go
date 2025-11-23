package dto

import "time"

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReplaceReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReplaceReviewerResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type PullRequest struct {
	PullRequestID    string     `json:"pull_request_id"`
	PullRequestName  string     `json:"pull_request_name"`
	AuthorID         string     `json:"author_id"`
	Status           string     `json:"status"`
	ReplaceReviewers []string   `json:"assigned_reviewers"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	MergedAt         *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}
