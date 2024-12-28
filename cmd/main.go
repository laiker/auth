package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/config/env"
	"github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
	db *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	fmt.Printf("%+v %v", request, ctx)
	return &user_v1.CreateResponse{}, nil
}

func (s *server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	sBuilder := sq.Select("id", "name", "email", "role", "createdAt", "updatedAt").From("auth_user").Where(sq.Eq{"id": request.Id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var name, email string
	var role int
	var createdAt time.Time
	var updatedAt sql.NullTime
	fmt.Println(s.db)
	err = s.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	return &user_v1.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      user_v1.Role(role),
		CreatedAt: &timestamp.Timestamp{},
		UpdatedAt: &timestamp.Timestamp{},
	}, nil
}

func (s *server) Update(ctx context.Context, request *user_v1.UpdateRequest) (*empty.Empty, error) {
	fmt.Printf("%+v %v", request, ctx)
	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*empty.Empty, error) {
	fmt.Printf("%+v %v", request, ctx)
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()

	errConfig := config.Load(configPath)
	ctx := context.Background()

	if errConfig != nil {
		log.Fatalf("failed to load config: %v", errConfig)
	}

	gConfig, err := env.NewGRPCConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	listener, err := net.Listen("tcp", gConfig.Address())

	if err != nil {
		log.Fatalf("file to start server: %v", err)
	}

	pgConfig, errc := env.NewPGConfig()

	if errc != nil {
		log.Fatalf("failed to load config: %v", errc)
	}

	p, errp := pgxpool.Connect(ctx, pgConfig.DSN())

	g := grpc.NewServer()
	reflection.Register(g)

	user_v1.RegisterUserV1Server(g, &server{db: p})

	log.Printf("server listening at %v", listener.Addr())

	if err = g.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if errp != nil {
		log.Fatalf("failed to connect: %v", errp)
	}

	defer p.Close()
}
