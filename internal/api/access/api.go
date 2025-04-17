package access

import (
	"context"
	"fmt"
	"strings"

	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/pkg/access_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const authPrefix = "Bearer "

type ServerAccess struct {
	access_v1.UnimplementedAccessV1Server
	AuthService   service.AuthService
	AccessService service.AccessService
}

func NewAccessServer(
	AuthService service.AuthService,
	AccessService service.AccessService,
) *ServerAccess {
	return &ServerAccess{
		AuthService:   AuthService,
		AccessService: AccessService,
	}
}

func (s *ServerAccess) HasAccess(ctx context.Context, req *access_v1.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := s.AuthService.VerifyAccessToken(ctx, accessToken)

	if err != nil {
		return nil, errors.New("access token is invalid")
	}
	fmt.Printf("%v", claims.Role)
	hasEndpointAccess, err := s.AccessService.HasAccessRight(ctx, req.EndpointAddress, claims.Role)

	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	if hasEndpointAccess {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New("access denied")
}
