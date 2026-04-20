package board

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

type Board interface {
	Get(ctx context.Context, opts entity.GetBoardsOpts) (boards entity.Boards, err error)
	GetByID(ctx context.Context, userID, boardID string) (board entity.Board, err error)
	GetByUserID(ctx context.Context, userID string) (boards entity.Boards, err error)
	Create(ctx context.Context, boardCreate entity.BoardCreate) (board entity.Board, err error)
	Update(ctx context.Context, userID string, boardUpdate entity.BoardUpdate) (board entity.Board, err error)
	Delete(ctx context.Context, userID, boardID string) (err error)
}

func New(repo *postgres.Repository) Board {
	return &BoardService{
		repo: repo,
	}
}
