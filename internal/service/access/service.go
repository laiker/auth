package access

import (
	"context"

	"github.com/laiker/auth/internal/repository"
	"github.com/laiker/auth/internal/service"
)

type accessService struct {
	repo repository.AccessRepository
}

func NewService(repo repository.AccessRepository) service.AccessService {
	return &accessService{
		repo: repo,
	}
}

func (s *accessService) HasAccessRight(ctx context.Context, endpoint string, role string) (bool, error) {
	permission, err := s.repo.GetEndpointPermission(ctx, endpoint)

	if err != nil {
		return false, err
	}

	mrole, errs := s.repo.GetRole(ctx, role)

	if errs != nil {
		return false, errs
	}

	return permission.MinPriority <= mrole.Priority, nil
}
