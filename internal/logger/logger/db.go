package logger

import (
	"context"
	_ "context"
	"log"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/auth/client/db"
	"github.com/laiker/auth/internal/logger"
)

type DBLogger struct {
	*slog.Logger
	db db.Client
}

func NewDBLogger(db db.Client, logger *slog.Logger) *DBLogger {
	return &DBLogger{db: db, Logger: logger}
}

func (l *DBLogger) Log(ctx context.Context, data logger.LogData) error {

	sBuilder := sq.Insert("auth_user_log").
		Columns("name", "entity_id").
		Values(data.Name, data.EntityID).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return err
	}

	q := db.Query{
		Name:     "log",
		QueryRaw: query,
	}

	l.Logger.Info("Database Operation:", data)

	_, err = l.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to insert log: %v\n", err)
		return err
	}

	return nil
}
