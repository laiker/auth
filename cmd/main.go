package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	user_v1.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	fmt.Printf("%+v %v", request, ctx)
	return &user_v1.CreateResponse{}, nil
}

func (s *server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	fmt.Printf("%+v %v", request, ctx)
	return &user_v1.GetResponse{
		Id: 100,
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
	// nolint:gosec
	listener, err := net.Listen("tcp", ":50052")

	if err != nil {
		panic(err)
	}

	g := grpc.NewServer()
	reflection.Register(g)
	user_v1.RegisterUserV1Server(g, &server{})

	log.Printf("server listening at %v", listener.Addr())

	if err = g.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
