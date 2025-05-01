package env

import (
	"net"
	"os"

	"github.com/laiker/auth/internal/config"
	"github.com/pkg/errors"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HOST"
	prometheusPortEnvName = "PROMETHEUS_PORT"
)

var _ config.PrometheusConfig = (*prometheusConfig)(nil)

type prometheusConfig struct {
	Host string
	Port string
}

func NewPrometheusConfig() (*prometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		Host: host,
		Port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}
