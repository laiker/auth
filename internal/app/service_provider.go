package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/laiker/auth/client/db"
	"github.com/laiker/auth/client/db/pg"
	"github.com/laiker/auth/client/db/transaction"
	api "github.com/laiker/auth/internal/api/user"
	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/config/env"
	"github.com/laiker/auth/internal/logger/logger"
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
	db             db.Client
	txManager      db.TxManager
	dbLogger       *logger.DBLogger
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

func (s *ServiceProvider) DB(ctx context.Context) db.Client {
	if s.db == nil {
		p, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}

		s.db = p
	}
	return s.db
}

func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DB(ctx).DB())
	}

	return s.txManager
}

func (s *ServiceProvider) UserRepository(ctx context.Context) repository.UserRepository {

	if s.userRepository == nil {
		r := repo.NewRepository(s.DB(ctx))
		s.userRepository = r
	}

	return s.userRepository
}

func (s *ServiceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		r := serv.NewService(s.UserRepository(ctx), s.TxManager(ctx), *s.DBLogger(ctx))
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

func (s *ServiceProvider) DBLogger(ctx context.Context) *logger.DBLogger {
	if s.dbLogger == nil {
		l := logger.NewDBLogger(s.DB(ctx))
		s.dbLogger = l
	}

	return s.dbLogger
}
