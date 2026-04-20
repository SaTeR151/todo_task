package task

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_TASKS = "tasks"
)

type TaskStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
}

func New(storage *pgxpool.Pool, scheme string) (*TaskStorage, error) {
	return &TaskStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
	}, nil
}
