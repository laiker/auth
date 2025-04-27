package logger

import (
	_ "context"
	"log/slog"

	slogmulti "github.com/samber/slog-multi"
)

func InitLogger(handlers ...slog.Handler) *slog.Logger {
	slogger := slog.New(slogmulti.Fanout(handlers...))
	slog.SetDefault(slogger)

	return slogger
}
