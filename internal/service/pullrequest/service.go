package pr

import (
	"pr-service/internal/repository"
	def "pr-service/internal/service"
)

var _ def.PullRequestService = (*service)(nil)

type service struct {
	prRepository   repository.PullRequestRepository
	userRepository repository.UserRepository
}

func NewService(prRepository repository.PullRequestRepository, userRepository repository.UserRepository) *service {
	return &service{
		prRepository:   prRepository,
		userRepository: userRepository,
	}
}
