package test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/laiker/auth/client/db"
	log "github.com/laiker/auth/internal/logger"
	"github.com/laiker/auth/internal/logger/logger"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
	serv "github.com/laiker/auth/internal/service/user"
	. "github.com/ovechkin-dm/mockio/mock"
	"golang.org/x/net/context"
)

type fields struct {
	repo      repository.UserRepository
	txManager db.TxManager
	logger    logger.DBLoggerInterface
}

type TestDependencies struct {
	UserRepositoryMock repository.UserRepository
	txManagerMock      db.TxManager
	loggerMock         logger.DBLoggerInterface
	contextMock        context.Context
}

func SetupServiceTest(t *testing.T) *TestDependencies {
	t.Helper()

	r := Mock[repository.UserRepository]()
	tx := Mock[db.TxManager]()
	dblogger := Mock[logger.DBLoggerInterface]()

	deps := &TestDependencies{
		UserRepositoryMock: r,
		txManagerMock:      tx,
		loggerMock:         dblogger,
		contextMock:        context.Background(),
	}

	return deps
}

func Test_serv_Create(t *testing.T) {
	deps := SetupServiceTest(t)

	type args struct {
		ctx      context.Context
		userInfo *model.UserInfo
	}

	ld := log.LogData{
		Name:     "create user",
		EntityID: int64(1),
	}

	When(deps.loggerMock.Log(deps.contextMock, ld)).ThenReturn(nil)

	callback := func(args []any) []any {
		fn := args[1].(func(context.Context) (int64, error))
		id, err := fn(deps.contextMock)
		return []any{id, err}
	}

	When(deps.txManagerMock.ReadCommitted(Any[context.Context](), Any[db.Handler]())).
		ThenReturn(int64(1), nil).
		ThenAnswer(callback)

	pw := gofakeit.Password(true, true, true, true, true, 10)
	name := gofakeit.Name()
	email := gofakeit.Email()

	mi := &model.UserInfo{
		Name:     name,
		Email:    email,
		Role:     "USER",
		Password: pw,
	}

	When(deps.UserRepositoryMock.Create(deps.contextMock, mi)).ThenReturn(int64(1), nil)

	a := args{
		ctx:      deps.contextMock,
		userInfo: mi,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    int64(1),
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.UserRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serv.NewService(tt.fields.repo, tt.fields.txManager, tt.fields.logger)
			got, err := s.Create(tt.args.ctx, tt.args.userInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

		})
	}

}

func Test_serv_Delete(t *testing.T) {
	deps := SetupServiceTest(t)

	When(deps.UserRepositoryMock.Delete(deps.contextMock, int64(1))).ThenReturn(nil)

	type args struct {
		ctx context.Context
		id  int64
	}

	a := args{
		ctx: deps.contextMock,
		id:  int64(1),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Success Test",
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.UserRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serv.NewService(tt.fields.repo, tt.fields.txManager, tt.fields.logger)
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_serv_Get(t *testing.T) {
	deps := SetupServiceTest(t)

	name := gofakeit.Name()
	email := gofakeit.Email()

	mi := &model.User{
		Id:        1,
		Name:      name,
		Email:     email,
		Role:      "USER",
		UpdatedAt: sql.NullTime{},
		CreatedAt: time.Time{},
	}

	When(deps.UserRepositoryMock.Get(deps.contextMock, int64(1))).ThenReturn(mi, nil)

	type args struct {
		ctx context.Context
		id  int64
	}

	a := args{
		ctx: deps.contextMock,
		id:  int64(1),
	}

	tests := []struct {
		name    string
		want    *model.User
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    mi,
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.UserRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serv.NewService(tt.fields.repo, tt.fields.txManager, tt.fields.logger)
			got, err := s.Get(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serv_Update(t *testing.T) {
	deps := SetupServiceTest(t)

	name := gofakeit.Name()
	email := gofakeit.Email()

	mi := &model.User{
		Id:        1,
		Name:      name,
		Email:     email,
		Role:      "USER",
		UpdatedAt: sql.NullTime{},
		CreatedAt: time.Time{},
	}

	When(deps.UserRepositoryMock.Update(deps.contextMock, mi)).ThenReturn(nil)

	type args struct {
		ctx       context.Context
		modelUser *model.User
	}

	a := args{
		ctx:       deps.contextMock,
		modelUser: mi,
	}

	tests := []struct {
		name    string
		want    *model.User
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Success Test",
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.UserRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serv.NewService(tt.fields.repo, tt.fields.txManager, tt.fields.logger)
			err := s.Update(tt.args.ctx, tt.args.modelUser)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
