package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"pr-service/internal/api/handlers"
	"pr-service/internal/config/db"
	"pr-service/internal/config/server"
	infra "pr-service/internal/db"
	"pr-service/internal/repository"
	prRepo "pr-service/internal/repository/pullrequest"
	teamRepo "pr-service/internal/repository/team"
	userRepo "pr-service/internal/repository/user"
	"pr-service/internal/service"
	prService "pr-service/internal/service/pullrequest"
	teamService "pr-service/internal/service/team"
	userService "pr-service/internal/service/user"
)

type serviceProvider struct {
	httpConfig server.HTTPConfig
	dbConfig   db.DBConfig

	db *pgxpool.Pool

	userRepository repository.UserRepository
	teamRepository repository.TeamRepository
	prRepository   repository.PullRequestRepository

	userService service.UserService
	teamService service.TeamService
	prService   service.PullRequestService

	userHandler *handlers.UserHandler
	teamHandler *handlers.TeamHandler
	prHandler   *handlers.PullRequestHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HTTPConfig() server.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := server.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}
		s.httpConfig = cfg
	}
	return s.httpConfig
}

func (s *serviceProvider) DBConfig() db.DBConfig {
	if s.dbConfig == nil {
		cfg, err := db.NewDBConfig()
		if err != nil {
			log.Fatalf("failed to get db config: %s", err.Error())
		}
		s.dbConfig = cfg
	}
	return s.dbConfig
}

func (s *serviceProvider) DB(ctx context.Context) *pgxpool.Pool {
	if s.db == nil {
		pool, err := infra.InitDB(s.DBConfig())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}
		s.db = pool
	}
	return s.db
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.DB(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) TeamRepository(ctx context.Context) repository.TeamRepository {
	if s.teamRepository == nil {
		s.teamRepository = teamRepo.NewRepository(s.DB(ctx))
	}
	return s.teamRepository
}

func (s *serviceProvider) PRRepository(ctx context.Context) repository.PullRequestRepository {
	if s.prRepository == nil {
		s.prRepository = prRepo.NewRepository(s.DB(ctx))
	}
	return s.prRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		)
	}
	return s.userService
}

func (s *serviceProvider) TeamService(ctx context.Context) service.TeamService {
	if s.teamService == nil {
		s.teamService = teamService.NewService(
			s.TeamRepository(ctx),
			s.UserRepository(ctx),
		)
	}
	return s.teamService
}

func (s *serviceProvider) PRService(ctx context.Context) service.PullRequestService {
	if s.prService == nil {
		s.prService = prService.NewService(
			s.PRRepository(ctx),
			s.UserRepository(ctx),
		)
	}
	return s.prService
}

func (s *serviceProvider) UserHandler(ctx context.Context) *handlers.UserHandler {
	if s.userHandler == nil {
		s.userHandler = handlers.NewUserHandler(s.UserService(ctx), s.PRService(ctx))
	}

	return s.userHandler
}

func (s *serviceProvider) TeamHandler(ctx context.Context) *handlers.TeamHandler {
	if s.teamHandler == nil {
		s.teamHandler = handlers.NewTeamHandler(s.TeamService(ctx))
	}

	return s.teamHandler
}

func (s *serviceProvider) PullRequestHandler(ctx context.Context) *handlers.PullRequestHandler {
	if s.prHandler == nil {
		s.prHandler = handlers.NewPullRequestHandler(s.PRService(ctx))
	}
	return s.prHandler
}

func (s *serviceProvider) Close() {
	if s.db != nil {
		s.db.Close()
		log.Println("db conn close")
	}
}
