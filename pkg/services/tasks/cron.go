package tasks

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

var _ cron.Logger = (*CronSlogAdapter)(nil)

type CronSlogAdapter struct {
	l *slog.Logger
}

func WrapCronLogger(l *slog.Logger) *CronSlogAdapter {
	return &CronSlogAdapter{l: l}
}

func (a *CronSlogAdapter) Info(msg string, keysAndValues ...any) {
	a.l.Info(msg, keysAndValues...)
}

func (a *CronSlogAdapter) Error(err error, msg string, keysAndValues ...any) {
	args := append([]any{"error", err}, keysAndValues...)
	a.l.Error(msg, args...)
}
