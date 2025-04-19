package user

import (
	"time"

	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/internal/utils"
	"golang.org/x/net/context"
)

const JwtAccessExpireTime = 24 * time.Hour
const JwtRefreshExpireTime = 30 * 24 * time.Hour

type authService struct {
	jwtConfig config.JwtConfig
}

func NewService(config config.JwtConfig) service.AuthService {
	return &authService{
		jwtConfig: config,
	}
}

func (s *authService) GetAccessToken(ctx context.Context, claims model.UserJwt) (string, error) {
	token, err := utils.GenerateToken(claims, []byte(s.jwtConfig.GetAccessSecret()), JwtAccessExpireTime)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, claims model.UserJwt) (string, error) {
	token, err := utils.GenerateToken(claims, []byte(s.jwtConfig.GetRefreshSecret()), JwtRefreshExpireTime)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) VerifyRefreshToken(ctx context.Context, token string) (model.UserClaims, error) {
	claims, err := utils.VerifyToken(token, []byte(s.jwtConfig.GetRefreshSecret()))

	if err != nil || claims == nil {
		return model.UserClaims{}, err
	}

	return *claims, nil
}

func (s *authService) VerifyAccessToken(ctx context.Context, token string) (model.UserClaims, error) {
	claims, err := utils.VerifyToken(token, []byte(s.jwtConfig.GetAccessSecret()))

	if err != nil || claims == nil {
		return model.UserClaims{}, err
	}

	return *claims, nil
}
