package gormdb

import (
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	"gorm.io/gorm"
)

type Config struct {
	// Type of database.
	// Supported values:
	//  - mysql
	//  - sqlite
	Type string `env:"TYPE" envDefault:"sqlite"`

	// DSN - data source name.
	//
	// Supported formats:
	//  - sqlite: "filepath?options" | ":memory:?options"
	//  - mysql: "username:password@tcp(ip:port)/dbname?options"
	//
	// SQLite options see here: https://github.com/mattn/go-sqlite3
	//
	// MySQL options see here: https://github.com/go-sql-driver/mysql
	DSN string `env:"DSN" envDefault:":memory:"`
}

type Package struct {
	DB *gorm.DB
}

func MakePackage(
	config Config,
	logger *logger.Package,
) (*Package, error) {
	var (
		res Package
		err error
	)

	switch config.Type {
	case "mysql":
		res.DB, err = NewMySQL(
			config,
			logger.Logger,
		)
		if err != nil {
			return nil, err
		}

	default:
		res.DB, err = NewSQLite(
			config,
			logger.Logger,
		)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}
