package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/calyrexx/zeroslog"
	"github.com/sater-151/todo-list/internal/app"
	"github.com/sater-151/todo-list/internal/configuration"
	"github.com/sater-151/todo-list/internal/credentials"
)

const version = "v1.0.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := configuration.NewConfig()
	if err != nil {
		slog.Error("NewConfig initialization failed.", slog.String("error", err.Error()))

		return
	}
	config.Version = version

	logger := slog.New(zeroslog.New(
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
		zeroslog.WithOutput(os.Stderr),
		zeroslog.WithColors(),
		zeroslog.WithMinLevel(config.Logger.Level),
	))

	slog.SetDefault(logger)

	logger.Info("starting application...", "version", version)

	creds, err := credentials.NewCredentials()
	if err != nil {
		logger.Error("newCredentials initialization failed", slog.String("error", err.Error()))

		return
	}

	application, err := app.New(ctx, &app.Dependencies{
		Configuration: config,
		Credentials:   creds,
	})
	if err != nil {
		logger.Error("application initialization failed", slog.String("error", err.Error()))

		return
	}

	wg := &sync.WaitGroup{}

	errCh := application.Start(ctx, wg)

	select {
	case err := <-errCh:
		logger.Error("service failed", slog.String("error", err.Error()))
		stop()
	case <-ctx.Done():
		logger.Info("received shutdown signal")
	}

	logger.Info("please wait, services are stopping...", "version", version)
	wg.Wait()

	logger.Info("application is stopped correctly. The force will be with you")
}
