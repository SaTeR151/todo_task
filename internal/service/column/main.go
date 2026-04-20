package column

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

type Column interface {
	Get(ctx context.Context, opts entity.GetColumnsOpts) (columns entity.Columns, err error)
	GetByID(ctx context.Context, boardID, columnID string) (column entity.Column, err error)
	GetByBoardID(ctx context.Context, boardID string) (columns entity.Columns, err error)
	CreateColumn(ctx context.Context, columnCreate entity.ColumnCreate) (column entity.Column, err error)
	UpdateColumn(ctx context.Context, boardID string, columnUpdate entity.ColumnUpdate) (column entity.Column, err error)
	DeleteColumn(ctx context.Context, columnID string) (err error)
	SwapColumns(ctx context.Context, boardID, columnIDA, columnIDB string) (err error)
}

func New(repo *postgres.Repository) Column {
	return &ColumnService{
		repo: repo,
	}
}
