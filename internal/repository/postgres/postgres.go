package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sater-151/todo-list/internal/credentials"
	"github.com/sater-151/todo-list/internal/repository"
)

type Repository struct {
	WhatsApp repository.ITodoTask
}

const (
	migrationsPath = "migrations/postgres"
)

func NewPostgres(ctx context.Context, c credentials.Postgres) (*Repository, error) {
	pool, err := pgxpool.New(ctx, c.ConnString)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	defer func() {
		if err != nil && pool != nil {
			pool.Close()
		}
	}()

	slog.Info("connection opened")

	if err = Ping(ctx, pool, 5*time.Second); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	slog.Info("applying migrations", slog.String("path", migrationsPath))

	return InitRepoRegistry(pool)
}

func InitRepoRegistry(postgresConnect *pgxpool.Pool) (*Repository, error) {
	whatsAppRepo := NewTodoTaskRepo(postgresConnect)

	return &Repository{
		WhatsApp: whatsAppRepo,
	}, nil
}

func Ping(ctx context.Context, pinger *pgxpool.Pool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := pinger.Ping(ctx); err != nil {
		return fmt.Errorf("cannot ping with timeout [%v]: %w", timeout, err)
	}

	return nil
}
