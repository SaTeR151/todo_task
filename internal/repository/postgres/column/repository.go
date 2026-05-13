package column

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_COLUMNS = "columns"
)

type ColumnStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
}

func New(storage *pgxpool.Pool, scheme string) (*ColumnStorage, error) {
	return &ColumnStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
	}, nil
}
