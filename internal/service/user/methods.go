package user

import (
	"context"
	"fmt"
	"pr-service/internal/model"
)

func (s *service) CreateUser(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	if user.Username == "" {
		return fmt.Errorf("user name is required")
	}
	if user.TeamName == "" {
		return fmt.Errorf("team ID is required")
	}

	return s.userRepository.CreateUser(ctx, user)
}

func (s *service) SetIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.IsActive = isActive

	if err := s.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	return s.userRepository.GetUserByID(ctx, id)
}


func (s *service) UpdateUser(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}

	_, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.userRepository.UpdateUser(ctx, user)
}
