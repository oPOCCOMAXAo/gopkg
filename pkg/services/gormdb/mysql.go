package gormdb

import (
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(
	config Config,
	logger *slog.Logger,
) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: NewLogger(logger),

		TranslateError: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
