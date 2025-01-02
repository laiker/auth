package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
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

	sRole := strconv.Itoa(int(user_v1.Role_value[string(request.Role)]))

	sBuilder := sq.Insert("auth_user").
		Columns("email", "name", "password", "role").
		Values(request.Email, request.Name, request.Password, sRole).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)

	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	return &user_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	sBuilder := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth_user").
		Where(sq.Eq{"id": request.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var name, email string
	var role string
	var createdAt time.Time
	var updatedAt sql.NullTime
	fmt.Println(query, args)
	err = s.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)

	if err != nil {
		log.Fatalf("failed to select user: %v", err)
	}

	srole, err := strconv.Atoi(role)

	if err != nil {
		log.Fatalf("failed to convert role to int: %v", err)
	}

	return &user_v1.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      user_v1.Role(srole),
		CreatedAt: &timestamp.Timestamp{},
		UpdatedAt: &timestamp.Timestamp{},
	}, nil
}

func (s *server) Update(ctx context.Context, request *user_v1.UpdateRequest) (*empty.Empty, error) {

	sBuilder := sq.Update("auth_user").
		PlaceholderFormat(sq.Dollar).
		Set("email", request.Email.GetValue()).
		Set("name", request.Name.GetValue()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": request.Id.GetValue()})

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.db.Exec(ctx, query, args...)

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	return nil, nil
}

func (s *server) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*empty.Empty, error) {
	sBuilder := sq.Delete("auth_user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": request.Id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

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

	if errp != nil {
		log.Fatalf("failed to connect: %v", errp)
	}

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
