package service

import (
	"context"
	"pr-service/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	//GetActiveUserFromTeam(ctx context.Context, teamID string) ([]model.User, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id string) error
	AddMemberInTeam(ctx context.Context, teamID, userID string) error
	RemoveMember(ctx context.Context, teamID, userID string) error
}

type PullRequestService interface {
	CreatePR(ctx context.Context, namePR, authorID string) (*model.PullRequest, error)
	GetPRByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequest, error)
	GetPRByPRID(ctx context.Context, prID string) (*model.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
}
