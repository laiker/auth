package user

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
	"github.com/laiker/auth/pkg/user_v1"
	"golang.org/x/net/context"
)

const (
	tableName = "auth_user"

	idColumn        = "id"
	nameColumn      = "name"
	passwordColumn  = "password"
	roleColumn      = "role"
	emailColumn     = "email"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userInfo *model.UserInfo) (int64, error) {
	sRole := strconv.Itoa(int(user_v1.Role_value[userInfo.Role]))

	sBuilder := sq.Insert(tableName).
		Columns(emailColumn, nameColumn, passwordColumn, roleColumn).
		Values(userInfo.Email, userInfo.Name, userInfo.Password, sRole).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Println("failed to build query: %v", err)
	}

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)

	if err != nil {
		log.Println("failed to insert user: %v", err)
	}

	return 1, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {

	sBuilder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Println("failed to build query: %v", err)
	}

	var idc int64
	var name, email string
	var role string
	var createdAt time.Time
	var updatedAt sql.NullTime
	fmt.Println(query, args)
	err = r.db.QueryRow(ctx, query, args...).Scan(&idc, &name, &email, &role, &createdAt, &updatedAt)

	if err != nil {
		log.Println("failed to select user: %v", err)
	}

	srole, err := strconv.Atoi(role)

	if err != nil {
		log.Println("failed to convert role to int: %v", err)
	}

	return &model.User{
		Id:        idc,
		Name:      name,
		Email:     email,
		Role:      strconv.Itoa(srole),
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {

	sBuilder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Println("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to delete user: %v", err)
	}

	return nil
}

func (r *repo) Update(ctx context.Context, info *model.User) error {

	sBuilder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(emailColumn, info.Email).
		Set(nameColumn, info.Name).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: info.Id})

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Println("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		log.Println("failed to update user: %v", err)
	}

	return nil
}
