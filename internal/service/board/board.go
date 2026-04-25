package board

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/pkg/utils"
)

type BoardService struct {
	repo *postgres.Repository
}

func (s *BoardService) Get(ctx context.Context, opts entity.GetBoardsOpts) (boards entity.Boards, err error) {
	boards, err = s.repo.Board.Get(ctx, opts)
	if err != nil {
		return
	}

	if len(boards) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *BoardService) GetByID(ctx context.Context, userID, boardID string) (board entity.Board, err error) {
	defer utils.AddFuncLabel("[service-get-board-by-id]", err)

	boards, err := s.Get(ctx, entity.GetBoardsOpts{ID: boardID, UserID: userID})
	if err != nil {
		return
	}

	return boards[0], nil
}

func (s *BoardService) GetByUserID(ctx context.Context, userID string) (boards entity.Boards, err error) {
	defer utils.AddFuncLabel("[service-get-boards-by-user-id]", err)

	return s.Get(ctx, entity.GetBoardsOpts{UserID: userID})
}

func (s *BoardService) Create(ctx context.Context, boardCreate entity.BoardCreate) (board entity.Board, err error) {
	defer utils.AddFuncLabel("[service-create-board]", err)

	newBoardID, err := s.repo.Board.Create(ctx, boardCreate)
	if err != nil {
		return
	}

	// Таблица содержит дефолтную колонку backlog
	columnCreate := entity.ColumnCreate{
		BoardID:    newBoardID,
		Name:       "backlog",
		OderNumber: -1,
	}

	_, err = s.repo.Column.CreateColumn(ctx, columnCreate)
	if err != nil {
		return
	}

	return s.GetByID(ctx, boardCreate.UserID, newBoardID)
}

func (s *BoardService) Update(ctx context.Context, userID string, boardUpdate entity.BoardUpdate) (board entity.Board, err error) {
	defer utils.AddFuncLabel("[service-update-board]", err)

	if err = s.repo.Board.Update(ctx, boardUpdate); err != nil {
		return
	}

	return s.GetByID(ctx, userID, boardUpdate.ID)
}

func (s *BoardService) Delete(ctx context.Context, userID, boardID string) (err error) {
	defer utils.AddFuncLabel("[service-delete-board]", err)

	_, err = s.GetByID(ctx, userID, boardID)
	if err != nil {
		return
	}

	return s.repo.Board.Delete(ctx, boardID)
}
