package ginserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/pkg/errors"
)

var (
	_ lifecycle.Servable     = (*Service)(nil)
	_ lifecycle.Shutdownable = (*Service)(nil)
)

type Service struct {
	server *http.Server
	engine *gin.Engine
}

func NewService(
	config Config,
) *Service {
	res := &Service{
		engine: gin.New(),
	}

	res.server = &http.Server{
		Handler:           res.engine,
		WriteTimeout:      time.Minute,
		ReadTimeout:       time.Minute,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		Addr:              ":" + strconv.FormatInt(config.Port, 10),
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	return res
}

func (s *Service) Serve(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)

	return errors.WithStack(err)
}

func (s *Service) Engine() *gin.Engine {
	return s.engine
}

func (s *Service) Router() gin.IRouter {
	return s.engine
}
