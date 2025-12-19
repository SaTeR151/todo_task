package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/sater-151/todo-list/internal/api/rest"
	"github.com/sater-151/todo-list/internal/api/rest/handlers"
	"github.com/sater-151/todo-list/internal/configuration"
	"github.com/sater-151/todo-list/internal/credentials"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
)

type (
	Dependencies struct {
		Configuration *configuration.Configurations `validate:"required"`
		Credentials   *credentials.Credentials      `validate:"required"`
	}

	App struct {
		repository *Repository
		rest       *rest.Server
		usecases   *Usecases
		logger     *slog.Logger
	}
)

func New(ctx context.Context, d *Dependencies) (*App, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewApp", d, err)
	}

	repo, err := NewRepo(ctx, *d.Credentials.Postgres)
	if err != nil {
		return nil, err
	}

	uc, err := NewUsecases(&UsecasesDependencies{
		Repository: repo,
	})
	if err != nil {
		return nil, err
	}

	todoTaskHandlers, err := handlers.NewTodoTaskHandlers(&handlers.TodoTaskServerDependencies{
		TodoTaskUsecase: uc.TodoTask,
		TodoTaskRepo:    repo.TodoTask,
		Password:        d.Credentials.Data.Password,
	})
	if err != nil {
		return nil, err
	}

	router, err := rest.NewRouter(&rest.RouterDependencies{
		Handlers: todoTaskHandlers,
	})
	if err != nil {
		return nil, err
	}

	server, err := rest.NewServer(rest.ServerDependencies{
		Handler: router,
		Config:  d.Configuration.HTTPServer,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		repository: repo,
		rest:       server,
		usecases:   uc,
		logger:     slog.With(slog.String("component", "app")),
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	errCh := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.rest.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- fmt.Errorf("rest.Start error: %w", err)
		}
	}()

	return errCh
}
