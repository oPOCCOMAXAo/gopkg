package mysql

import (
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(
	config Config,
	logger *slog.Logger,
) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
		),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
