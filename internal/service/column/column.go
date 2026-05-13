package column

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/service/task"
	"github.com/sater-151/todo-list/pkg/utils"
)

type ColumnService struct {
	columns     Repository
	tasks       TaskRepository
	taskService task.Task
}

func (s *ColumnService) Get(ctx context.Context, opts entity.GetColumnsOpts) (columns entity.Columns, err error) {

	columns, err = s.columns.GetColumns(ctx, opts)
	if err != nil {
		return
	}

	if len(columns) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *ColumnService) GetByID(ctx context.Context, boardID, columnID string) (column entity.Column, err error) {
	colums, err := s.Get(ctx, entity.GetColumnsOpts{ID: columnID, BoardID: boardID})
	if err != nil {
		return
	}

	return colums[0], nil
}

func (s *ColumnService) GetByBoardID(ctx context.Context, boardID string) (columns entity.Columns, err error) {
	return s.Get(ctx, entity.GetColumnsOpts{BoardID: boardID})
}

func (s *ColumnService) CreateColumn(ctx context.Context, columnCreate entity.ColumnCreate) (column entity.Column, err error) {
	defer utils.AddFuncLabel("[service-create-column]", err)

	newColumnID, err := s.columns.CreateColumn(ctx, columnCreate)
	if err != nil {
		return
	}

	return s.GetByID(ctx, columnCreate.BoardID, newColumnID)
}

func (s *ColumnService) UpdateColumn(ctx context.Context, boardID string, columnUpdate entity.ColumnUpdate) (column entity.Column, err error) {
	defer utils.AddFuncLabel("[service-update-column]", err)

	if err = s.columns.UpdateColumn(ctx, columnUpdate); err != nil {
		return
	}

	return s.GetByID(ctx, boardID, columnUpdate.ID)
}

func (s *ColumnService) DeleteColumn(ctx context.Context, boardID, columnID string) (err error) {
	defer utils.AddFuncLabel("[service-delete-column]", err)

	tasks, err := s.tasks.Get(ctx, entity.GetTasksOpts{ColumnID: columnID})
	if err != nil {
		return
	}

	var backlogColumn entity.Column

	if len(tasks) > 0 {
		columns, err := s.Get(ctx, entity.GetColumnsOpts{Name: "backlog"})
		if err != nil {
			return err
		}

		if len(columns) == 0 {
			err = entity.ErrNotFound
			return err
		}

		backlogColumn = columns[0]
	}

	for _, task := range tasks {
		_, err := s.taskService.Move(ctx, boardID, task.ID, backlogColumn.ID)
		if err != nil {
			return err
		}
	}

	return s.columns.DeleteColumn(ctx, columnID)
}

func (s *ColumnService) SwapColumns(ctx context.Context, boardID, columnIDA, columnIDB string) (err error) {
	defer utils.AddFuncLabel("[service-swap-columns]", err)

	columnA, err := s.GetByID(ctx, boardID, columnIDA)
	if err != nil {
		return
	}

	orderColumnA := columnA.OrderNumber

	columnB, err := s.GetByID(ctx, boardID, columnIDB)
	if err != nil {
		return
	}

	orderColumnB := columnB.OrderNumber

	columnUpdateA := entity.ColumnUpdate{
		ID:          columnIDA,
		OrderNumber: &orderColumnB,
	}

	columnUpdateB := entity.ColumnUpdate{
		ID:          columnIDB,
		OrderNumber: &orderColumnA,
	}

	if err = s.columns.UpdateColumn(ctx, columnUpdateA); err != nil {
		return
	}

	return s.columns.UpdateColumn(ctx, columnUpdateB)

}
