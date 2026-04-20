package task

import (
	"context"
	"fmt"
	"slices"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/pkg/utils"
)

type TaskService struct {
	repo *postgres.Repository
}

func (s *TaskService) Get(ctx context.Context, opts entity.GetTasksOpts) (tasks entity.Tasks, err error) {
	tasks, err = s.repo.Task.Get(ctx, opts)
	if err != nil {
		return
	}

	if len(tasks) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *TaskService) GetByID(ctx context.Context, boardID, taskID string) (res entity.Task, err error) {
	columns, err := s.repo.Column.GetColumns(ctx, entity.GetColumnsOpts{BoardID: boardID})
	if err != nil {
		return
	}

	if len(columns) == 0 {
		return entity.Task{}, fmt.Errorf("invalid board id: %s", boardID)
	}

	IDs := columns.GetIDs()

	tasks, err := s.Get(ctx, entity.GetTasksOpts{ID: taskID, ColumnIDs: IDs})
	if err != nil {
		return
	}

	return tasks[0], nil
}

func (s *TaskService) GetByColumnID(ctx context.Context, boardID, columnID string) (res entity.Tasks, err error) {
	columns, err := s.repo.Column.GetColumns(ctx, entity.GetColumnsOpts{BoardID: boardID})
	if err != nil {
		return
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("invalid board id: %s", boardID)
	}

	IDs := columns.GetIDs()

	if !slices.Contains(IDs, columnID) {
		return nil, fmt.Errorf("invalid column id: %s", columnID)
	}

	tasks, err := s.Get(ctx, entity.GetTasksOpts{ColumnID: columnID})
	if err != nil {
		return
	}

	return tasks, nil
}

func (s *TaskService) GetByTypeID(ctx context.Context, typeID string) (res entity.Tasks, err error) {
	users, err := s.Get(ctx, entity.GetTasksOpts{TypeID: typeID})
	if err != nil {
		return
	}

	return users, nil
}

func (s *TaskService) GetByBoardID(ctx context.Context, boardID string) (res entity.Tasks, err error) {
	columns, err := s.repo.Column.GetColumns(ctx, entity.GetColumnsOpts{BoardID: boardID})
	if err != nil {
		return
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("invalid board id: %s", boardID)
	}

	columnIDs := columns.GetIDs()

	tasks, err := s.Get(ctx, entity.GetTasksOpts{ColumnIDs: columnIDs})
	if err != nil {
		return
	}

	return tasks, err
}

func (s *TaskService) Create(ctx context.Context, boardID string, taskCreate entity.TaskCreate) (res entity.Task, err error) {
	defer utils.AddFuncLabel("[service-create-task]", err)

	if taskCreate.ColumnID == "" {
		columns, err := s.repo.Column.GetColumns(ctx, entity.GetColumnsOpts{BoardID: boardID, Name: "backlog"})
		if err != nil {
			return entity.Task{}, err
		}

		backlogColumn := columns[0]

		taskCreate.ColumnID = backlogColumn.ID
	}

	newTypeID, err := s.repo.Task.Create(ctx, taskCreate)
	if err != nil {
		return
	}

	return s.GetByID(ctx, boardID, newTypeID)
}

func (s *TaskService) Update(ctx context.Context, boardID string, taskUpdate entity.TaskUpdate) (res entity.Task, err error) {
	defer utils.AddFuncLabel("[service-update-task]", err)

	_, err = s.GetByID(ctx, boardID, taskUpdate.ID)
	if err != nil {
		return
	}

	if err = s.repo.Task.Update(ctx, taskUpdate); err != nil {
		return
	}

	return s.GetByID(ctx, boardID, taskUpdate.ID)
}

func (s *TaskService) Delete(ctx context.Context, taskID string) (err error) {
	defer utils.AddFuncLabel("[service-delete-task]", err)

	return s.repo.Task.Delete(ctx, taskID)
}

func (s *TaskService) Move(ctx context.Context, boardID, taskID, newColumnID string) (res entity.Task, err error) {
	defer utils.AddFuncLabel("[service-move-task]", err)

	columns, err := s.repo.Column.GetColumns(ctx, entity.GetColumnsOpts{BoardID: boardID, ID: newColumnID})
	if err != nil {
		return
	}

	if len(columns) == 0 {
		return entity.Task{}, fmt.Errorf("invalid column id: %s", newColumnID)
	}

	task, err := s.GetByID(ctx, boardID, taskID)
	if err != nil {
		return
	}

	moveEventCreate := entity.MoveEventCreate{
		TaskID:       taskID,
		ToColumnID:   newColumnID,
		FromColumnID: task.ColumnID,
	}

	s.repo.MoveEvent.Create(ctx, moveEventCreate)

	taskUpdate := entity.TaskUpdate{
		ID:       taskID,
		ColumnID: &newColumnID,
	}

	if err = s.repo.Task.Update(ctx, taskUpdate); err != nil {
		return
	}

	movedTask, err := s.GetByID(ctx, boardID, taskID)
	if err != nil {
		return
	}

	return movedTask, nil
}
