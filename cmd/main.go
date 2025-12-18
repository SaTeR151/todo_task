package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/calyrexx/zeroslog"
	"github.com/go-chi/chi/v5"
	"github.com/sater-151/todo-list/internal/app"
	"github.com/sater-151/todo-list/internal/configuration"
	"github.com/sater-151/todo-list/internal/credentials"
	"github.com/sater-151/todo-list/internal/database"
	"github.com/sater-151/todo-list/internal/handlers"
	"github.com/sater-151/todo-list/internal/service"
	"gl.iteco.com/technology/go_services/toolbox/sl"
)

const version = "v1.0.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := configuration.NewConfig()
	if err != nil {
		slog.Error("NewConfig initialization failed.", sl.Error(err))

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
		logger.Error("newCredentials initialization failed", sl.Error(err))

		return
	}

	application, err := app.New(ctx, &app.Dependencies{
		Configuration: config,
		Credentials:   creds,
	})
	if err != nil {
		logger.Error("application initialization failed", sl.Error(err))

		return
	}

	wg := &sync.WaitGroup{}

	errCh := application.Start(ctx, wg)

	select {
	case err := <-errCh:
		logger.Error("service failed", sl.Error(err))
		stop()
	case <-ctx.Done():
		logger.Info("received shutdown signal")
	}

	logger.Info("please wait, services are stopping...", "version", version)
	wg.Wait()

	application.ShutdownTracer(ctx)
	logger.Info("application is stopped correctly. The force will be with you")

	db, err := database.OpenDB(configuration.Database.DbFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer db.Close()

	service := service.New(db)

	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Get("/api/tasks", handlers.Auth(handlers.ListTask(service)))
	r.Get("/api/task", handlers.Auth(handlers.GetTask(db)))

	r.Post("/api/task", handlers.Auth(handlers.PostTask(service)))
	r.Post("/api/task/done", handlers.Auth(handlers.PostTaskDone(service)))
	r.Post("/api/signin", handlers.Sign)

	r.Put("/api/task", handlers.Auth(handlers.PutTask(service)))

	r.Delete("/api/task", handlers.Auth(handlers.DeleteTask(db)))

	log.Println("Server start at port:", configuration.HttpClient.Port)
	if err := http.ListenAndServe(":"+configuration.HttpClient.Port, r); err != nil {
		log.Println("Ошибка запуска сервера:", err.Error())
		return
	}
}
