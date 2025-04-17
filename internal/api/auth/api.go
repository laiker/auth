package auth

import (
	"context"
	"fmt"

	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/internal/utils"
	"github.com/laiker/auth/pkg/auth_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAuth struct {
	auth_v1.UnimplementedAuthV1Server
	AuthService service.AuthService
	UserService service.UserService
}

func NewAuthServer(
	AuthService service.AuthService,
	UserService service.UserService,
) *ServerAuth {
	return &ServerAuth{
		AuthService: AuthService,
		UserService: UserService,
	}
}

func (s *ServerAuth) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {

	user, err := s.UserService.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if !utils.VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("Неверный логин, пароль")
	}

	fmt.Printf("%+v\n", user)
	mu := model.UserJwt{
		UserId: user.Id,
		Role:   user.Role,
	}
	fmt.Printf("%+v\n", mu)
	accessToken, err := s.AuthService.GetAccessToken(ctx, mu)

	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.AuthService.GetRefreshToken(ctx, mu)

	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &auth_v1.LoginResponse{RefreshToken: refreshToken, AccessToken: accessToken}, nil
}

func (s *ServerAuth) GetRefreshToken(ctx context.Context, req *auth_v1.GetRefreshTokenRequest) (*auth_v1.GetRefreshTokenResponse, error) {
	claims, err := s.AuthService.VerifyRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	mu := model.UserJwt{
		UserId: claims.UserId,
		Role:   claims.Role,
	}

	refreshToken, err := s.AuthService.GetRefreshToken(ctx, mu)

	if err != nil {
		return nil, err
	}

	return &auth_v1.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (s *ServerAuth) GetAccessToken(ctx context.Context, req *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {
	claims, err := s.AuthService.VerifyRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid access token")
	}

	mu := model.UserJwt{
		UserId: claims.UserId,
		Role:   claims.Role,
	}

	accessToken, err := s.AuthService.GetRefreshToken(ctx, mu)

	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
