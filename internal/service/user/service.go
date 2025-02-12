package user

import (
	"github.com/laiker/auth/client/db"
	log "github.com/laiker/auth/internal/logger"
	"github.com/laiker/auth/internal/logger/logger"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
	"github.com/laiker/auth/internal/service"
	"golang.org/x/net/context"
)

type serv struct {
	repo      repository.UserRepository
	txManager db.TxManager
	logger    logger.DBLoggerInterface
}

func NewService(repo repository.UserRepository, manager db.TxManager, logger logger.DBLoggerInterface) service.UserService {
	return &serv{repo: repo, txManager: manager, logger: logger}
}

func (s *serv) Create(ctx context.Context, userInfo *model.UserInfo) (int64, error) {

	id, err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) (int64, error) {
		var id int64
		var errTx error
		id, errTx = s.repo.Create(ctx, userInfo)
		if errTx != nil {
			return int64(0), errTx
		}

		logData := log.LogData{
			Name:     "create user",
			EntityID: id,
		}

		errTx = s.logger.Log(ctx, logData)

		if errTx != nil {
			return int64(0), errTx
		}

		return id, nil
	})

	if err != nil {
		return int64(0), err
	}

	return id, nil
}

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *serv) Update(ctx context.Context, info *model.User) error {
	return s.repo.Update(ctx, info)
}
