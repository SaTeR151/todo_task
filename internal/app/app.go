package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sater-151/todo-list/internal/configuration"
	"github.com/sater-151/todo-list/internal/credentials"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/sirupsen/logrus"
	"gl.iteco.com/technology/go_general/errproc"
	"gl.iteco.com/technology/go_services/toolbox/cron"
	"gl.iteco.com/technology/go_services/toolbox/sl"
)

type (
	Dependencies struct {
		Configuration *configuration.Configurations `validate:"required"`
		Credentials   *credentials.Credentials      `validate:"required"`
	}

	App struct {
		repository *Repository
		rest       *rest.Server
		appCron    *cron.Cron
		usecases   *Usecases
		logger     *slog.Logger
	}
)

func New(ctx context.Context, d *Dependencies) (*App, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewApp", d, err)
	}

	errorProc, err := errproc.NewErrProc(logrus.New())
	if err != nil {
		return nil, err
	}

	repo, err := NewRepo(ctx, d.Credentials.Postgres, tr)
	if err != nil {
		return nil, err
	}

	uc, err := NewUsecases(&UsecasesDependencies{
		Repository: repo,
		GrpcClient: grpcClient,
		Onemsg:     integrations.OnemsgProvider,
		API1C:      integrations.API1C,
		Conf:       d.Configuration,
		ErrorProc:  errorProc,
		Metrics:    metricApp,
		Tracer:     tr,
	})
	if err != nil {
		return nil, err
	}

	router, err := rest.NewRouter(&rest.RouterDependencies{
		Metrics:              promhttp.Handler(),
		OnemsgWebhookHandler: onemsgWebhook,
		Tracer:               tr,
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

	appCron.RegisterJobs(
		cron.Job{
			Config: cron.JobConfig(d.Configuration.AppCron.Send),
			Fn:     uc.WhatsApp.SendIteration,
		},
		cron.Job{
			Config: cron.JobConfig(d.Configuration.AppCron.SendPushes),
			Fn:     uc.WhatsApp.SendPushesIteration,
		},
		cron.Job{
			Config: cron.JobConfig(d.Configuration.AppCron.Cleaner),
			Fn:     uc.WhatsApp.Cleaner,
		},
		cron.Job{
			Config: cron.JobConfig(d.Configuration.AppCron.Send1CIteration),
			Fn:     uc.WhatsApp.Send1CIteration,
		},
	)

	if d.Credentials.Otel.Environment != "local" {
		sentry.CaptureEvent(&sentry.Event{
			Message: "Service started",
			Release: d.Configuration.Version,
		})
	}

	return &App{
		httpClient: httpClient,
		appCron:    appCron,
		repository: repo,
		grpcClient: grpcClient,
		grpcServer: grpcServer,
		rest:       server,
		usecases:   uc,
		tracer:     tr,
		logger:     slog.With(sl.Component("app")),
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	errCh := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.grpcServer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- fmt.Errorf("grpcServer.Start error: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.appCron.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- fmt.Errorf("appCron.Start error: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.rest.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- fmt.Errorf("rest.Start error: %w", err)
		}
	}()

	return errCh
}

func (a *App) ShutdownTracer(ctx context.Context) {
	if a.tracer != nil {
		if err := a.tracer.Shutdown(ctx); err != nil {
			a.logger.Error("tracer shutdown failed", sl.Error(err))
		}
	}
}
