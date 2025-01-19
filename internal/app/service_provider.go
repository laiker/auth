package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	api "github.com/laiker/auth/internal/api/user"
	"github.com/laiker/auth/internal/closer"
	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/config/env"
	"github.com/laiker/auth/internal/repository"
	repo "github.com/laiker/auth/internal/repository/user"
	"github.com/laiker/auth/internal/service"
	serv "github.com/laiker/auth/internal/service/user"
)

type ServiceProvider struct {
	pgConfig       config.PGConfig
	grpcConfig     config.GRPCConfig
	pgPool         *pgxpool.Pool
	userRepository repository.UserRepository
	userService    service.UserService
	userApi        *api.Server
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) PGConfig() config.PGConfig {

	if s.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()

		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.pgConfig = pgConfig
	}

	return s.pgConfig
}

func (s *ServiceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {

		gConfig, err := env.NewGRPCConfig()

		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.grpcConfig = gConfig

	}

	return s.grpcConfig
}

func (s *ServiceProvider) PGPool(ctx context.Context) *pgxpool.Pool {

	if s.pgPool == nil {
		p, err := pgxpool.Connect(ctx, s.PGConfig().DSN())

		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}

		s.pgPool = p

		closer.Add(func() error {
			s.pgPool.Close()
			return nil
		})
	}

	return s.pgPool
}

func (s *ServiceProvider) UserRepository(ctx context.Context) repository.UserRepository {

	if s.userRepository == nil {
		r := repo.NewRepository(s.PGPool(ctx))
		s.userRepository = r
	}

	return s.userRepository
}

func (s *ServiceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		r := serv.NewService(s.UserRepository(ctx))
		s.userService = r
	}

	return s.userService
}

func (s *ServiceProvider) UserApi(ctx context.Context) *api.Server {
	if s.userApi == nil {
		a := api.NewServer(s.UserService(ctx))
		s.userApi = a
	}

	return s.userApi
}
