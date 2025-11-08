package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server/http"
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
	logg := logger.New(cfg.Logger.Level)
	storage := memorystorage.New()
	calendar := app.New(logg, storage)

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
