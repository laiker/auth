package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	dbName = "PG_DATABASE_NAME"
	dbUser = "PG_USER"
	dbPass = "PG_PASSWORD"
	dbPort = "PG_PORT"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	name string
	user string
	pass string
	port string
}

func NewPGConfig() (PGConfig, error) {
	name := os.Getenv(dbName)
	if len(name) == 0 {
		return nil, errors.New("pg name not found")
	}
	user := os.Getenv(dbUser)
	if len(user) == 0 {
		return nil, errors.New("pg user not found")
	}
	pass := os.Getenv(dbPass)
	if len(pass) == 0 {
		return nil, errors.New("pg pass not found")
	}
	port := os.Getenv(dbPort)
	if len(port) == 0 {
		return nil, errors.New("pg port not found")
	}

	return &pgConfig{
		name: name,
		user: user,
		pass: pass,
		port: port,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("host=localhost port=%v dbname=%v user=%v password=%v sslmode=disable", cfg.name, cfg.user, cfg.pass, cfg.port)
}
