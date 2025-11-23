package team

import (
	"pr-service/internal/repository"
	def "pr-service/internal/service"
)

var _ def.TeamService = (*service)(nil)

type service struct {
	teamRepository repository.TeamRepository
	userRepository repository.UserRepository
}

func NewService(teamRepository repository.TeamRepository, userRepository repository.UserRepository) *service {
	return &service{
		teamRepository: teamRepository,
		userRepository: userRepository,
	}
}
