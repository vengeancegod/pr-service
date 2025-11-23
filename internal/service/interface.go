package service

import (
	"context"
	"pr-service/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	SetIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error)
	//GetActiveUserFromTeam(ctx context.Context, teamID string) ([]model.User, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByName(ctx context.Context, teamName string) (*model.Team, error)
	GetActiveTeamMembers(ctx context.Context, teamName string) ([]model.User, error)
}

type PullRequestService interface {
	CreatePR(ctx context.Context, prID, namePR, authorID string) (*model.PullRequest, error)
	GetPRByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequest, error)
	GetPRByPRID(ctx context.Context, prID string) (*model.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID string) (*model.PullRequest, string, error)
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
}
