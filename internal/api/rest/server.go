package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/sater-151/todo-list/internal/configuration"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
)

type (
	ServerDependencies struct {
		Handler http.Handler              `validate:"required"`
		Config  *configuration.HTTPServer `validate:"required"`
	}

	Server struct {
		srvHTTP         *http.Server
		shutdownTimeout time.Duration
	}
)

func NewServer(dep ServerDependencies) (*Server, error) {
	if err := validate.Struct(dep); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewServer", dep, err)
	}

	return &Server{
		srvHTTP: &http.Server{
			Addr:              ":" + dep.Config.Port,
			Handler:           dep.Handler,
			ReadTimeout:       dep.Config.ReadTimeout,
			ReadHeaderTimeout: dep.Config.ReadHeaderTimeout,
			WriteTimeout:      dep.Config.WriteTimeout,
			IdleTimeout:       dep.Config.IdleTimeout,
			MaxHeaderBytes:    dep.Config.MaxHeaderBytes,
		},
		shutdownTimeout: dep.Config.ShutdownTimeout,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		errChan <- s.srvHTTP.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.Stop(ctx)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srvHTTP.Shutdown(ctx)
}
