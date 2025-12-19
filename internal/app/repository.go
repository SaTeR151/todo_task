package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sater-151/todo-list/internal/credentials"
	"github.com/sater-151/todo-list/internal/repository"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

const (
	migrationsPath = "migrations/postgres"
)

type Repository struct {
	TodoTask repository.ITodoTask
}

func NewRepo(ctx context.Context, c credentials.Postgres) (*Repository, error) {
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

	if err = repository.Ping(ctx, pool, 5*time.Second); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	slog.Info("applying migrations", slog.String("path", migrationsPath))

	err = applyMigrations(
		ctx, goose.DialectPostgres, 2*time.Minute, postgresURL(pool), migrationsPath,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot apply migrations: %w", err)
	}

	return initRepoRegistry(pool)
}

func initRepoRegistry(postgresConnect *pgxpool.Pool) (*Repository, error) {
	todoTaskRepo := postgres.NewTodoTaskRepo(postgresConnect)

	return &Repository{
		TodoTask: todoTaskRepo,
	}, nil
}

func postgresURL(pool *pgxpool.Pool) string {
	conf := pool.Config().ConnConfig

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		conf.User, conf.Password, net.JoinHostPort(conf.Host, strconv.Itoa(int(conf.Port))), conf.Database,
	)
}

func applyMigrations(
	ctx context.Context,
	dialect goose.Dialect,
	timeout time.Duration,
	connString string,
	migrationPath string,
) (err error) {
	db, err := goose.OpenDBWithDriver(string(dialect), connString)
	if err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = goose.UpContext(ctx, db, migrationPath)
	if err != nil {
		return fmt.Errorf("cannot up migrations: %w", err)
	}

	return nil
}
