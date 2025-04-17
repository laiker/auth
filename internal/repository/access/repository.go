package access

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/auth/client/db"
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/internal/repository"
)

const (
	tableName = "permission"

	idColumn              = "permission_id"
	resourceNameColumn    = "resource_name"
	minRolePriorityColumn = "min_role_priority"
)

type accessRepo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {

	return &accessRepo{db: db}
}

func (r *accessRepo) GetEndpointPermission(ctx context.Context, endpoint string) (*model.Permission, error) {
	fmt.Println("1")
	sBuilder := sq.Select(idColumn, resourceNameColumn, minRolePriorityColumn).
		From(tableName).
		Where(sq.Eq{"resource_name": endpoint}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	q := db.Query{
		Name:     "access.GetPermissionByRole",
		QueryRaw: query,
	}
	fmt.Println("2")
	permission := model.Permission{}

	err = r.db.DB().ScanOneContext(ctx, &permission, q, args...)

	if err != nil {
		log.Printf("failed to select user: %v\n", err)
	}

	return &permission, nil
}

func (r *accessRepo) GetRole(ctx context.Context, role string) (*model.Role, error) {

	sBuilder := sq.Select("role_id", "role_name", "priority").
		From("user_role").
		Where(sq.Eq{"role_name": role}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	q := db.Query{
		Name:     "access.GetRole",
		QueryRaw: query,
	}

	mrole := model.Role{}

	err = r.db.DB().ScanOneContext(ctx, &mrole, q, args...)

	if err != nil {
		log.Printf("failed to select user: %v\n", err)
	}

	return &mrole, nil
}
