package board

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_BOARDS = "boards"
)

type BoardStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
}

func New(storage *pgxpool.Pool, scheme string) (*BoardStorage, error) {
	return &BoardStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
	}, nil
}
