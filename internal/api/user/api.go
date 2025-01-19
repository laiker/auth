package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/auth/internal/converter"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}

func (s *Server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {

	userID, err := s.userService.Create(ctx, converter.ToUserFromCreateRequest(request))

	if err != nil {
		return nil, err
	}

	return &user_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	user, err := s.userService.Get(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	nt := sql.NullTime{
		Time:  time.Time{},
		Valid: true, // Указываем, что значение действительно
	}

	var ts *timestamppb.Timestamp
	if nt.Valid {
		ts = timestamppb.New(nt.Time) // Преобразуем time.Time в timestamp.Timestamp
	} else {
		ts = nil // Если значение не действительно, устанавливаем в nil
	}

	return &user_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user_v1.Role(user_v1.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: ts,
	}, nil
}

func (s *Server) Update(ctx context.Context, request *user_v1.UpdateRequest) (*empty.Empty, error) {

	err := s.userService.Update(ctx, converter.ToUserFromUpdateRequest(request))

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	return nil, nil
}

func (s *Server) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*empty.Empty, error) {

	err := s.userService.Delete(ctx, request.GetId())

	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	return &empty.Empty{}, nil
}
