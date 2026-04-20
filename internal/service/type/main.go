package type_service

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

type Type interface {
	Create(ctx context.Context, typeCreate entity.TypeCreate) (res entity.Type, err error)
	Update(ctx context.Context, userID string, typeUpdate entity.TypeUpdate) (res entity.Type, err error)
	Delete(ctx context.Context, typeID string) (err error)
	Get(ctx context.Context, opts entity.GetTypesOpts) (res entity.Types, err error)
	GetByUserID(ctx context.Context, userID string) (res entity.Types, err error)
	GetByID(ctx context.Context, typeID string) (res entity.Type, err error)
}

func New(repo *postgres.Repository) Type {
	return &TypeService{
		repo: repo,
	}
}
