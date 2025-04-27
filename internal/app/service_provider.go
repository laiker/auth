package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/laiker/auth/client/db"
	"github.com/laiker/auth/client/db/pg"
	"github.com/laiker/auth/client/db/transaction"
	accessApi "github.com/laiker/auth/internal/api/access"
	authApi "github.com/laiker/auth/internal/api/auth"
	userApi "github.com/laiker/auth/internal/api/user"
	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/config/env"
	"github.com/laiker/auth/internal/logger/logger"
	"github.com/laiker/auth/internal/repository"
	accessRepository "github.com/laiker/auth/internal/repository/access"
	repo "github.com/laiker/auth/internal/repository/user"
	"github.com/laiker/auth/internal/service"
	accessService "github.com/laiker/auth/internal/service/access"
	authService "github.com/laiker/auth/internal/service/auth"
	serv "github.com/laiker/auth/internal/service/user"
	"github.com/lmittmann/tint"
)

type ServiceProvider struct {
	//Configs
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	jwtConfig     config.JwtConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	//User
	userApi        *userApi.ServerUser
	userService    service.UserService
	userRepository repository.UserRepository

	//Auth
	authApi     *authApi.ServerAuth
	authService service.AuthService

	//Access
	accessApi        *accessApi.ServerAccess
	accessService    service.AccessService
	accessRepository repository.AccessRepository

	//Database
	db        db.Client
	txManager db.TxManager

	//Loggers
	dbLogger *logger.DBLogger
	logger   *slog.Logger
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

func (s *ServiceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {

		hConfig, err := env.NewHTTPConfig()

		if err != nil {
			s.Logger().Error("failed to load config: %v", err)
			os.Exit(1)
		}

		s.httpConfig = hConfig

	}

	return s.httpConfig
}

func (s *ServiceProvider) SwaggerConfig() config.HTTPConfig {
	if s.swaggerConfig == nil {

		sConfig, err := env.NewSwaggerConfig()

		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.swaggerConfig = sConfig

	}

	return s.swaggerConfig
}

func (s *ServiceProvider) DB(ctx context.Context) db.Client {
	if s.db == nil {
		p, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			s.Logger().Error("failed to connect: %v", err)
			os.Exit(1)
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

func (s *ServiceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		r := authService.NewService(s.JwtConfig())
		s.authService = r
	}

	return s.authService
}

func (s *ServiceProvider) UserApi(ctx context.Context) *userApi.ServerUser {
	if s.userApi == nil {
		a := userApi.NewUserServer(s.UserService(ctx))
		s.userApi = a
	}

	return s.userApi
}

func (s *ServiceProvider) AccessApi(ctx context.Context) *accessApi.ServerAccess {
	if s.accessApi == nil {
		a := accessApi.NewAccessServer(s.AuthService(ctx), s.AccessService(ctx), s.Logger())
		s.accessApi = a
	}

	return s.accessApi
}

func (s *ServiceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		r := accessService.NewService(s.AccessRepository(ctx))
		s.accessService = r
	}

	return s.accessService
}

func (s *ServiceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {

	if s.accessRepository == nil {
		r := accessRepository.NewRepository(s.DB(ctx), s.Logger())
		s.accessRepository = r
	}

	return s.accessRepository
}

func (s *ServiceProvider) JwtConfig() config.JwtConfig {
	if s.jwtConfig == nil {

		jwtConfig, err := env.NewJwtConfig()

		if err != nil {
			s.Logger().Error("failed to load config: %v", err)
			os.Exit(1)
		}

		s.jwtConfig = jwtConfig

	}

	return s.jwtConfig
}

func (s *ServiceProvider) AuthApi(ctx context.Context) *authApi.ServerAuth {
	if s.authApi == nil {
		a := authApi.NewAuthServer(
			s.AuthService(ctx),
			s.UserService(ctx),
		)
		s.authApi = a
	}

	return s.authApi
}

func (s *ServiceProvider) DBLogger(ctx context.Context) *logger.DBLogger {
	if s.dbLogger == nil {
		l := logger.NewDBLogger(s.DB(ctx), s.Logger())
		s.dbLogger = l
	}

	return s.dbLogger
}

func (s *ServiceProvider) Logger() *slog.Logger {
	if s.logger == nil {
		color := tint.NewHandler(os.Stdout, &tint.Options{
			Level:     slog.LevelDebug,
			AddSource: true,
		})

		l := logger.InitLogger(color)
		s.logger = l
	}

	return s.logger
}
