package user

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/auth/client/db"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
	"golang.org/x/net/context"
)

const (
	tableName = "auth_user"

	idColumn        = "id"
	nameColumn      = "name"
	passwordColumn  = "password"
	roleColumn      = "role_id"
	emailColumn     = "email"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userInfo *model.UserInfo) (int64, error) {

	sBuilder := sq.Insert(tableName).
		Columns(emailColumn, nameColumn, passwordColumn, roleColumn).
		Values(userInfo.Email, userInfo.Name, userInfo.Password, userInfo.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	var userID int64

	q := db.Query{
		Name:     "user.create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)

	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {

	sBuilder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	q := db.Query{
		Name:     "user.get",
		QueryRaw: query,
	}

	user := model.User{}

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)

	if err != nil {
		log.Printf("failed to select user: %v\n", err)
	}

	return &user, nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {

	sBuilder := sq.Select(idColumn, nameColumn, emailColumn, "role_name as role", passwordColumn).
		From(tableName).
		Where(sq.Eq{emailColumn: email}).
		Join("user_role on auth_user.role_id = user_role.role_id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	q := db.Query{
		Name:     "user.getByEmail",
		QueryRaw: query,
	}

	user := model.User{}

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)

	if err != nil {
		log.Printf("failed to select user: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {

	sBuilder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	q := db.Query{
		Name:     "user.delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to delete user: %v\n", err)
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
		log.Printf("failed to build query: %v\n", err)
	}

	q := db.Query{
		Name:     "user.update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to update user: %v\n", err)
	}

	return nil
}
