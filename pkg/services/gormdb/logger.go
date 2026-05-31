package gormdb

import (
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm/logger"
)

func NewLogger(
	logger *slog.Logger,
) logger.Interface {
	return slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
	)
}
