package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sater-151/todo-list/internal/models"
)

type ITodoTask interface {
	InsertTask(ctx context.Context, task models.Task) (string, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, uuid string) error
	Select(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error)
}

type Repository struct {
	TodoTask ITodoTask
}

func Ping(ctx context.Context, pinger *pgxpool.Pool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := pinger.Ping(ctx); err != nil {
		return fmt.Errorf("cannot ping with timeout [%v]: %w", timeout, err)
	}

	return nil
}
