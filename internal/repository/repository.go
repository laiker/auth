package repository

import (
	"context"

	"github.com/laiker/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, info *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type AccessRepository interface {
	GetEndpointPermission(ctx context.Context, endpoint string) (*model.Permission, error)
	GetRole(ctx context.Context, role string) (*model.Role, error)
}
