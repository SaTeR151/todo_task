package db_type

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_TYPES = "types"
)

type TypeStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
}

func New(storage *pgxpool.Pool, scheme string) (*TypeStorage, error) {
	return &TypeStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
	}, nil
}
