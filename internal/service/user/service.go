package user

import (
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
	"github.com/laiker/auth/internal/service"
	"golang.org/x/net/context"
)

type serv struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) service.UserService {
	return &serv{repo: repo}
}

func (s *serv) Create(ctx context.Context, userInfo *model.UserInfo) (int64, error) {
	return s.repo.Create(ctx, userInfo)
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
