package logger

import (
	"context"
	"fmt"
	"log/slog"
)

//nolint:forbidigo
func SlogToFmt(
	logger *slog.Logger,
	level slog.Level,
) func(format string, values ...any) {
	return func(format string, values ...any) {
		logger.Log(context.Background(), level, fmt.Sprintf(format, values...))
	}
}
