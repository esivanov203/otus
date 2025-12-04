package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/config"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

func runScheduler(_ *cobra.Command, _ []string) error {
	// config
	cfg, err := config.NewConfig(ConfigFile)
	if err != nil {
		return fmt.Errorf("config loading: %w", err)
	}

	// logger
	logg, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("logger initial: %w", err)
	}

	// storage
	var store storage.Storage
	switch cfg.Storage.Type {
	case "memory":
		store = memorystorage.New()
	case "sql":
		store = sqlstorage.New(cfg.Storage.Dsn)
		if err := store.Connect(); err != nil {
			return fmt.Errorf("connect to storage: %w", err)
		}
		defer func() { _ = store.Close() }()
	default:
		return fmt.Errorf("connect to storage: unknown storage type: %s", cfg.Storage.Type)
	}

	// queue
	q, err := queue.NewQueueService(cfg.Queue.Dsn)
	if err != nil {
		return fmt.Errorf("connect to queue server: %w", err)
	}
	defer func() { _ = q.Close() }()

	// context + signal shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	logg.Info("scheduler is running...")

	t, err := time.ParseDuration(cfg.Queue.Interval)
	if err != nil {
		return fmt.Errorf("parse config queue interval: %w", err)
	}

	// start scheduler in a goroutine
	go scheduler.Run(ctx, logg, store, q, t)

	// wait for signal
	<-ctx.Done()

	logg.Info("shutting down scheduler...")

	// graceful shutdown timeout
	shdCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-shdCtx.Done():
		logg.Info("forced shutdown")
	default:
	}

	logg.Info("scheduler stopped")
	return nil
}
