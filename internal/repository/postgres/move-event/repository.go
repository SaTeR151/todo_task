package moveevent

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_MOVE_EVENTS = "move_events"
)

type MoveEventStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
}

func New(storage *pgxpool.Pool, scheme string) (*MoveEventStorage, error) {
	return &MoveEventStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
	}, nil
}
