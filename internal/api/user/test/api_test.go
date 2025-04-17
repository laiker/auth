package test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/auth/internal/api/user"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/service"
	"github.com/laiker/auth/pkg/user_v1"
	. "github.com/ovechkin-dm/mockio/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type fields struct {
	UnimplementedUserV1Server user_v1.UnimplementedUserV1Server
	UserService               service.UserService
}

type TestDependencies struct {
	UserServiceMock service.UserService
	context         context.Context
}

func SetupApiTest(t *testing.T) *TestDependencies {
	t.Helper()

	m := Mock[service.UserService]()

	deps := &TestDependencies{
		UserServiceMock: m,
		context:         context.Background(),
	}

	return deps
}

func TestServer_Create(t *testing.T) {
	deps := SetupApiTest(t)

	type args struct {
		ctx     context.Context
		request *user_v1.CreateRequest
	}

	pw := gofakeit.Password(true, true, true, true, true, 10)

	name := gofakeit.Name()
	email := gofakeit.Email()

	a := args{
		ctx: deps.context,
		request: &user_v1.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        pw,
			PasswordConfirm: pw,
			Role:            0,
		},
	}

	mi := &model.UserInfo{
		Name:     name,
		Email:    email,
		Role:     0,
		Password: pw,
	}

	WhenDouble(deps.UserServiceMock.Create(a.ctx, mi)).ThenReturn(1, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user_v1.CreateResponse
		wantErr bool
	}{
		{
			name: "Success Test",
			want: &user_v1.CreateResponse{
				Id: int64(1),
			},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedUserV1Server: user_v1.UnimplementedUserV1Server{},
				UserService:               deps.UserServiceMock,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &user.ServerUser{
				UnimplementedUserV1Server: tt.fields.UnimplementedUserV1Server,
				UserService:               tt.fields.UserService,
			}
			got, err := s.Create(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Delete(t *testing.T) {

	type args struct {
		ctx     context.Context
		request *user_v1.DeleteRequest
	}

	m := Mock[service.UserService]()

	a := args{
		ctx: context.Background(),
		request: &user_v1.DeleteRequest{
			Id: 1,
		},
	}

	When(m.Delete(a.ctx, 1)).ThenReturn(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    &empty.Empty{},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedUserV1Server: user_v1.UnimplementedUserV1Server{},
				UserService:               m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &user.ServerUser{
				UnimplementedUserV1Server: tt.fields.UnimplementedUserV1Server,
				UserService:               tt.fields.UserService,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Get(t *testing.T) {
	deps := SetupApiTest(t)

	type args struct {
		ctx     context.Context
		request *user_v1.GetRequest
	}

	a := args{
		ctx: deps.context,
		request: &user_v1.GetRequest{
			Id: 1,
		},
	}

	name := gofakeit.Name()
	email := gofakeit.Email()
	id := int64(gofakeit.Number(1, 10))
	ctime := time.Date(2025, 02, 9, 13, 29, 58, 0, time.UTC)

	mu := &model.User{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      "USER",
		CreatedAt: ctime,
		UpdatedAt: sql.NullTime{
			Time:  ctime,
			Valid: true,
		},
	}

	When(deps.UserServiceMock.Get(a.ctx, id)).ThenReturn(mu, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user_v1.GetResponse
		wantErr bool
	}{
		{
			name: "Success Test",
			want: &user_v1.GetResponse{
				Id:        id,
				Name:      name,
				Email:     email,
				Role:      user_v1.Role(0),
				CreatedAt: timestamppb.New(ctime),
				UpdatedAt: timestamppb.New(ctime),
			},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedUserV1Server: user_v1.UnimplementedUserV1Server{},
				UserService:               deps.UserServiceMock,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &user.ServerUser{
				UnimplementedUserV1Server: tt.fields.UnimplementedUserV1Server,
				UserService:               tt.fields.UserService,
			}

			got, err := s.Get(tt.args.ctx, &user_v1.GetRequest{
				Id: id,
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Update(t *testing.T) {
	m := Mock[service.UserService]()

	type args struct {
		ctx     context.Context
		request *user_v1.UpdateRequest
	}

	name := gofakeit.Name()
	email := gofakeit.Email()
	id := int64(gofakeit.Number(1, 10))

	mu := &model.User{
		Id:    id,
		Name:  name,
		Email: email,
	}

	a := args{
		ctx: context.Background(),
		request: &user_v1.UpdateRequest{
			Id:    wrapperspb.Int64(id),
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
		},
	}

	When(m.Update(a.ctx, mu)).ThenReturn(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    nil,
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedUserV1Server: user_v1.UnimplementedUserV1Server{},
				UserService:               m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &user.ServerUser{
				UnimplementedUserV1Server: tt.fields.UnimplementedUserV1Server,
				UserService:               tt.fields.UserService,
			}
			got, err := s.Update(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
