package env

import (
	"net"
	"os"

	"github.com/laiker/auth/internal/config"
	"github.com/pkg/errors"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

var _ config.HTTPConfig = (*HTTPConfig)(nil)

type HTTPConfig struct {
	host string
	port string
}

func NewHTTPConfig() (*HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &HTTPConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *HTTPConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
