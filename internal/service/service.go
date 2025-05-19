package service

import (
	"context"

	"github.com/laiker/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, info *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	FindByName(ctx context.Context, name string) ([]*model.UserName, error)
}

type AuthService interface {
	GetAccessToken(ctx context.Context, model model.UserJwt) (string, error)
	GetRefreshToken(ctx context.Context, model model.UserJwt) (string, error)
	VerifyRefreshToken(ctx context.Context, token string) (model.UserClaims, error)
	VerifyAccessToken(ctx context.Context, token string) (model.UserClaims, error)
}

type AccessService interface {
	HasAccessRight(ctx context.Context, endpoint string, role string) (bool, error)
}
