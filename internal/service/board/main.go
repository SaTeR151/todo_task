package board

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
)

type BoardRepository interface {
	Get(ctx context.Context, opts entity.GetBoardsOpts) (entity.Boards, error)
	Create(ctx context.Context, boardCreate entity.BoardCreate) (string, error)
	Update(ctx context.Context, boardUpdate entity.BoardUpdate) error
	Delete(ctx context.Context, boardID string) error
}

type ColumnCreator interface {
	CreateColumn(ctx context.Context, columnCreate entity.ColumnCreate) (string, error)
}

type Board interface {
	Get(ctx context.Context, opts entity.GetBoardsOpts) (boards entity.Boards, err error)
	GetByID(ctx context.Context, userID, boardID string) (board entity.Board, err error)
	GetByUserID(ctx context.Context, userID string) (boards entity.Boards, err error)
	Create(ctx context.Context, boardCreate entity.BoardCreate) (board entity.Board, err error)
	Update(ctx context.Context, userID string, boardUpdate entity.BoardUpdate) (board entity.Board, err error)
	Delete(ctx context.Context, userID, boardID string) (err error)
}

func New(boardRepo BoardRepository, columnRepo ColumnCreator) Board {
	return &BoardService{
		boards:  boardRepo,
		columns: columnRepo,
	}
}
