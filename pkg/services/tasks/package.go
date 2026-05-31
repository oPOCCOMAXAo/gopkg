package tasks

import (
	"log/slog"
	"time"

	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
)

type Config struct {
	WorkerCount int           `env:"WORKER_COUNT" envDefault:"1"`
	QueueSize   int           `env:"QUEUE_SIZE"   envDefault:"0"` // 0 means unbounded
	Timeout     time.Duration `env:"TIMEOUT"      envDefault:"1h"`
	MaxStack    int           `env:"MAX_STACK"    envDefault:"30"`
}

type Package struct {
	Service *Service
}

func MakePackage(
	config Config,
	lifecycle *lifecycle.Package,
	logger *logger.Package,
) *Package {
	var res Package

	res.Service = NewService(
		config,
		logger.Logger.With(slog.String("service", "tasks")),
		logger.Logger.With(slog.String("service", "cron")),
	)

	lifecycle.Service.RegisterService(res.Service)

	return &res
}
