package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) error {
	_ = cmd
	_ = args
	// config
	cfg, err := NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// logger, storage, business logic
	logg, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("initial logger: %w", err)
	}

	var store storage.Storage
	switch cfg.Storage.Type {
	case "memory":
		store = memorystorage.New()
	case "sql":
		store = sqlstorage.New(cfg.Storage.Dsn)
		if err := store.Connect(); err != nil {
			logg.Error(fmt.Sprintf("connect to storage: %v", err))
			return fmt.Errorf("connect to storage: %w", err)
		}
		defer func() { _ = store.Close() }()
	default:
		return fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(logg, calendar, cfg.Server)
	grpcServer := internalgrpc.NewServer(logg, calendar, cfg.GRPCServer)

	// context with cancellation on signal
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")

	if err := grpcServer.PrepareListener(); err != nil {
		logg.Error(fmt.Sprintf("prepare listener: %v", err))
		return fmt.Errorf("prepare listener: %w", err)
	}

	// канал для сигнализации о падении любого сервера
	errCh := make(chan struct{}, 1)

	// запускаем серверы
	go server.Start(errCh)
	go grpcServer.Start(errCh)

	// ждём либо cancellation on signal, либо падение любого сервера
	select {
	case <-ctx.Done():
	case <-errCh:
		cancel()
	}

	// graceful shutdown
	shdCtx, shdcancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shdcancel()

	server.Stop(shdCtx)
	grpcServer.Stop(shdCtx)

	logg.Info("calendar stopped")

	return nil
}
