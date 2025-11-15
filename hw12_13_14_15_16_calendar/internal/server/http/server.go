package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type ServerConf struct {
	Host string `yaml:"host"` // listen interface
	Port int    `yaml:"port"` // listen port
}

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	logger     logger.Logger
	app        app.Application
}

func NewServer(logger logger.Logger, app app.Application, cfg ServerConf) *Server {
	router := mux.NewRouter()
	s := &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
		logger: logger,
		app:    app,
		router: router,
	}

	s.setRouterUses()
	s.setRouterHandlers()

	return s
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http server start failed", logger.Fields{"addr": s.httpServer.Addr})
		}
	}()
	s.logger.Info("http server starting", logger.Fields{"addr": s.httpServer.Addr})
	<-ctx.Done()

	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("http server stopping", logger.Fields{"addr": s.httpServer.Addr})

	return s.httpServer.Shutdown(ctx)
}
