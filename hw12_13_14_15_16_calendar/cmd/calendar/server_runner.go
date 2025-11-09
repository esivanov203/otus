package main

import (
	"context"
	"fmt"
	sqlstorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/sql"
	"os/signal"
	"syscall"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"

	internalhttp "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) error {
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
		if err := store.Connect(context.Background()); err != nil {
			return fmt.Errorf("connect to storage: %w", err)
		}
		defer func() { _ = store.Close(context.Background()) }()
	default:
		return fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}

	calendar := app.New(logg, store)

	// http server
	server := internalhttp.NewServer(logg, calendar)

	// context with cancellation on signal
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// graceful shutdown handler
	go func() {
		<-ctx.Done()
		shdCtx, shdcancel := context.WithTimeout(context.Background(), time.Second*3)
		defer shdcancel()
		if err := server.Stop(shdCtx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	// start http server (blocking)
	return server.Start(ctx)
}
