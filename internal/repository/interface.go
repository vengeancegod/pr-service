package repository

import (
	"context"
	"pr-service/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	GetActiveUserFromTeam(ctx context.Context, teamName string) ([]model.User, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByName(ctx context.Context, teamName string) (*model.Team, error)
}

type PullRequestRepository interface {
	CreatePR(ctx context.Context, pr *model.PullRequest) error
	GetPRByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequest, error)
	GetPRByPRID(ctx context.Context, id string) (*model.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
}
