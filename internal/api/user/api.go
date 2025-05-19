package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/auth/internal/converter"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerUser struct {
	user_v1.UnimplementedUserV1Server
	UserService service.UserService
}

func NewUserServer(userService service.UserService) *ServerUser {
	return &ServerUser{
		UserService: userService,
	}
}

func (s *ServerUser) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {

	userID, err := s.UserService.Create(ctx, converter.ToUserFromCreateRequest(request))

	if err != nil {
		return nil, err
	}

	return &user_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *ServerUser) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	user, err := s.UserService.Get(ctx, request.Id)

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

func (s *ServerUser) Update(ctx context.Context, request *user_v1.UpdateRequest) (*empty.Empty, error) {

	err := s.UserService.Update(ctx, converter.ToUserFromUpdateRequest(request))

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	return nil, nil
}

func (s *ServerUser) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*empty.Empty, error) {

	err := s.UserService.Delete(ctx, request.GetId())

	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	return &empty.Empty{}, nil
}

func (s *ServerUser) FindByLogin(ctx context.Context, request *user_v1.FindByLoginRequest) (*user_v1.FindByLoginResponse, error) {
	users, err := s.UserService.FindByName(ctx, request.GetName())

	if err != nil {
		log.Fatalf("failed to find users: %v", err)
	}

	userList := make([]*user_v1.UserSearchItem, len(users))

	for k, user := range users {
		userList[k] = &user_v1.UserSearchItem{
			Id:   user.Id,
			Name: user.Name,
		}
	}

	return &user_v1.FindByLoginResponse{
		Results: userList,
	}, nil
}
