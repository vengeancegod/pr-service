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
	if user.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	if user.TeamID == "" {
		return fmt.Errorf("team ID is required")
	}

	return s.userRepository.CreateUser(ctx, user)
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

func (s *service) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user ID is required")
	}

	return s.userRepository.DeleteUser(ctx, id)
}
