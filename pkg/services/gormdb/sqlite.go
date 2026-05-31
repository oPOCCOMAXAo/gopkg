package gormdb

import (
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite(
	config Config,
	logger *slog.Logger,
) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: NewLogger(logger),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
