package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/proto"
	"google.golang.org/grpc"
)

type GRPCServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCServer struct {
	proto.UnimplementedCalendarServiceServer
	logger     logger.Logger
	app        app.Application
	grpcServer *grpc.Server
	listener   net.Listener
	addr       string
}

func NewServer(logger logger.Logger, app app.Application, cfg GRPCServerConf) *GRPCServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	grpcSrv := grpc.NewServer()

	s := &GRPCServer{
		logger:     logger,
		app:        app,
		grpcServer: grpcSrv,
		addr:       addr,
	}

	proto.RegisterCalendarServiceServer(grpcSrv, s)

	return s
}

func (s *GRPCServer) PrepareListener() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	s.listener = lis
	return nil
}

func (s *GRPCServer) Start(chanErr chan struct{}) {
	s.logger.Info("grpc server starting", logger.Fields{"addr": s.addr})

	if err := s.grpcServer.Serve(s.listener); err != nil && errors.Is(err, grpc.ErrServerStopped) {
		s.logger.Error(fmt.Sprintf("grpc server start failed: %v", err), logger.Fields{"addr": s.addr})
	}

	select {
	case chanErr <- struct{}{}:
	default:
	}
}

func (s *GRPCServer) Stop(ctx context.Context) {
	s.logger.Info("grpc server stopping", logger.Fields{"addr": s.addr})

	gracefulDone := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(gracefulDone)
	}()

	// пытаемся остановить сервер корректно,
	// но по завершении контекста завершаем принудительно
	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
	case <-gracefulDone:
	}
}
