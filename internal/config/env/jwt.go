package env

import (
	"os"

	"github.com/laiker/auth/internal/config"
	"github.com/pkg/errors"
)

const (
	jwtAccessSecret  = "JWT_ACCESS_SECRET"  //nolint:golint,gosec
	jwtRefreshSecret = "JWT_REFRESH_SECRET" //nolint:golint,gosec
)

var _ config.JwtConfig = (*JwtConfig)(nil)

type JwtConfig struct {
	accessSecret  string
	refreshSecret string
}

func NewJwtConfig() (*JwtConfig, error) {
	accessSecret := os.Getenv(jwtAccessSecret)
	if len(accessSecret) == 0 {
		return nil, errors.New("jwt access token not found")
	}

	refreshSecret := os.Getenv(jwtRefreshSecret)
	if len(refreshSecret) == 0 {
		return nil, errors.New("jwt refresh token not found")
	}

	return &JwtConfig{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}, nil
}

func (cfg *JwtConfig) GetAccessSecret() string {
	return cfg.accessSecret
}

func (cfg *JwtConfig) GetRefreshSecret() string {
	return cfg.refreshSecret
}
