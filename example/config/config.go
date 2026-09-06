package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/opoccomaxao/gopkg/pkg/services/ginserver"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	pkgerr "github.com/pkg/errors"
)

type Config struct {
	Logger    logger.Config    `envPrefix:"LOGGER_"`
	Lifecycle lifecycle.Config `envPrefix:"LIFECYCLE_"`
	Server    ginserver.Config `envPrefix:"SERVER_"`
}

func Load() (*Config, error) {
	var res Config

	err := env.ParseWithOptions(&res, env.Options{
		RequiredIfNoDef:       false,
		UseFieldNameByDefault: false,
	})
	if err != nil {
		return nil, pkgerr.WithStack(err)
	}

	return &res, nil
}
