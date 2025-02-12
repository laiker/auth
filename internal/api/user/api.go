package user

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/auth/internal/converter"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	user_v1.UnimplementedUserV1Server
	UserService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		UserService: userService,
	}
}

func (s *Server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {

	userID, err := s.UserService.Create(ctx, converter.ToUserFromCreateRequest(request))

	if err != nil {
		return nil, err
	}

	return &user_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	user, err := s.UserService.Get(ctx, request.Id)
	fmt.Println(user, err)
	if err != nil {
		return nil, err
	}

	return &user_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user_v1.Role(user_v1.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
	}, nil
}

func (s *Server) Update(ctx context.Context, request *user_v1.UpdateRequest) (*empty.Empty, error) {

	err := s.UserService.Update(ctx, converter.ToUserFromUpdateRequest(request))

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	return nil, nil
}

func (s *Server) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*empty.Empty, error) {

	err := s.UserService.Delete(ctx, request.GetId())

	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	return &empty.Empty{}, nil
}
