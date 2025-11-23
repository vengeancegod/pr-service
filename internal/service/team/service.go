package team

import (
	"pr-service/internal/repository"
	def "pr-service/internal/service"
)

var _ def.TeamService = (*service)(nil)

type service struct {
	teamRepository repository.TeamRepository
}

func NewService(teamRepository repository.TeamRepository) *service {
	return &service{
		teamRepository: teamRepository,
	}
}
