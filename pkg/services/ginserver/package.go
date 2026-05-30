package ginserver

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	pkgLogger "github.com/opoccomaxao/gopkg/pkg/services/logger"
)

type Config struct {
	Port int64 `env:"PORT" envDefault:"8080"`
}

type Package struct {
	Service *Service
	Engine  *gin.Engine
	Router  gin.IRouter
}

func MakePackage(
	config Config,
	lifecycle *lifecycle.Package,
	logger *pkgLogger.Package,
) *Package {
	var res Package

	res.Service = NewService(config)

	res.Engine = res.Service.Engine()

	res.Router = res.Service.Router()

	lifecycle.Service.RegisterService(res.Service)

	gin.DebugPrintFunc = pkgLogger.SlogToFmt(logger.Logger, slog.LevelDebug)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Logger.Debug("Route",
			slog.String("method", httpMethod),
			slog.String("path", absolutePath),
			slog.String("handler", handlerName),
			slog.Int("num_handlers", nuHandlers),
		)
	}

	return &res
}
