package mysql

import (
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	"gorm.io/gorm"
)

type Config struct {
	DSN string `env:"DSN,required"`
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

	res.DB, err = New(
		config,
		logger.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
