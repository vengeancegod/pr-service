package team

import (
	"context"
	"fmt"

	"pr-service/internal/model"
)

func (s *service) CreateTeam(ctx context.Context, team *model.Team) error {
	err := s.teamRepository.CreateTeam(ctx, team)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	for i := range team.Members {
		user := &team.Members[i]
		user.TeamName = team.TeamName

		err := s.userRepository.CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.ID, err)
		}
	}

	return nil
}

func (s *service) GetTeamByName(ctx context.Context, teamName string) (*model.Team, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	return s.teamRepository.GetTeamByName(ctx, teamName)
}

func (s *service) GetActiveTeamMembers(ctx context.Context, teamName string) ([]model.User, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	_, err := s.teamRepository.GetTeamByName(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	return s.userRepository.GetActiveUserFromTeam(ctx, teamName)
}
