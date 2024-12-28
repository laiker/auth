package env

import (
	"net"
	"os"

	"github.com/laiker/auth/internal/config"
	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

var _ config.GRPCConfig = (*GrpcConfig)(nil)

type GrpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*GrpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &GrpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *GrpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
