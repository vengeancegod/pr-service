package team

import (
	"context"
	"fmt"

	"pr-service/internal/model"
)

func (s *service) CreateTeam(ctx context.Context, team *model.Team) error {
	if team.ID == "" {
		return fmt.Errorf("team ID is required")
	}
	if team.TeamName == "" {
		return fmt.Errorf("team name is required")
	}

	return s.teamRepository.CreateTeam(ctx, team)
}

func (s *service) UpdateTeam(ctx context.Context, team *model.Team) error {
	if team.ID == "" {
		return fmt.Errorf("team ID is required")
	}

	return s.teamRepository.UpdateTeam(ctx, team)
}

func (s *service) DeleteTeam(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("team ID is required")
	}

	return s.teamRepository.DeleteTeam(ctx, id)
}

func (s *service) AddMemberInTeam(ctx context.Context, teamID, userID string) error {
	if teamID == "" {
		return fmt.Errorf("team ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	return s.teamRepository.AddMemberInTeam(ctx, teamID, userID)
}

func (s *service) RemoveMember(ctx context.Context, teamID, userID string) error {
	if teamID == "" {
		return fmt.Errorf("team ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	return s.teamRepository.RemoveMember(ctx, teamID, userID)
}
