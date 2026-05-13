package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TABLE_USERS = "users"
)

type UserStorage struct {
	queryBuilder squirrel.StatementBuilderType
	client       *pgxpool.Pool
	scheme       string
	cryptoKey    string
}

func New(storage *pgxpool.Pool, scheme string, cryptoKey string) (*UserStorage, error) {
	return &UserStorage{
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		client:       storage,
		scheme:       scheme,
		cryptoKey:    cryptoKey,
	}, nil
}
