package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/config"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/email_service"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/sender"
	"github.com/spf13/cobra"
)

func runSender(_ *cobra.Command, _ []string) error {
	// config
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("config loading: %w", err)
	}

	// logger
	logg, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("logger initial: %w", err)
	}

	// queue
	q, err := queue.NewQueueService(cfg.Queue.Dsn)
	if err != nil {
		return fmt.Errorf("connect to queue server: %w", err)
	}
	defer func() { _ = q.Close() }()

	// email service
	var es emailservice.EmailSender
	if os.Getenv("INTEGRATION_TESTS") == "true" {
		es = &emailservice.IntegrationTestsEmailService{}
	} else {
		es = &emailservice.EmailService{}
	}

	// context + signal shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	logg.Info("sender is running...")

	// start sender in a goroutine
	go sender.Run(ctx, logg, q, es)

	// wait for signal
	<-ctx.Done()

	logg.Info("shutting down sender...")

	// graceful shutdown timeout
	shdCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-shdCtx.Done():
		logg.Info("forced shutdown")
	default:
	}

	logg.Info("sender stopped")
	return nil
}
