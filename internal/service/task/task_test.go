package task

import (
	"context"
	"errors"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

type fakeBoardRepo struct {
	boards entity.Boards
	err    error
}

func (r *fakeBoardRepo) Get(_ context.Context, opts entity.GetBoardsOpts) (entity.Boards, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Boards
	for _, board := range r.boards {
		if opts.ID != "" && board.ID != opts.ID {
			continue
		}
		if opts.UserID != "" && board.UserID != opts.UserID {
			continue
		}
		res = append(res, board)
	}
	return res, nil
}

type fakeColumnRepo struct {
	columns entity.Columns
	err     error
}

func (r *fakeColumnRepo) GetColumns(_ context.Context, opts entity.GetColumnsOpts) (entity.Columns, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Columns
	for _, column := range r.columns {
		if opts.ID != "" && column.ID != opts.ID {
			continue
		}
		if opts.BoardID != "" && column.BoardID != opts.BoardID {
			continue
		}
		if opts.Name != "" && column.Name != opts.Name {
			continue
		}
		res = append(res, column)
	}
	return res, nil
}

type fakeTypeRepo struct {
	types entity.Types
	err   error
}

func (r *fakeTypeRepo) Get(_ context.Context, opts entity.GetTypesOpts) (entity.Types, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Types
	for _, typ := range r.types {
		if opts.ID != "" && typ.ID != opts.ID {
			continue
		}
		if opts.UserID != "" && typ.UserID != opts.UserID {
			continue
		}
		if opts.Name != "" && typ.Name != opts.Name {
			continue
		}
		res = append(res, typ)
	}
	return res, nil
}

type fakeTaskRepo struct {
	tasks   entity.Tasks
	created []entity.TaskCreate
	updated []entity.TaskUpdate
	deleted []string
	err     error
}

func (r *fakeTaskRepo) Get(_ context.Context, opts entity.GetTasksOpts) (entity.Tasks, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Tasks
	for _, task := range r.tasks {
		if opts.ID != "" && task.ID != opts.ID {
			continue
		}
		if opts.ColumnID != "" && task.ColumnID != opts.ColumnID {
			continue
		}
		if opts.TypeID != "" && task.TypeID != opts.TypeID {
			continue
		}
		if len(opts.ColumnIDs) > 0 {
			found := false
			for _, id := range opts.ColumnIDs {
				if task.ColumnID == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		res = append(res, task)
	}
	return res, nil
}

func (r *fakeTaskRepo) Create(_ context.Context, taskCreate entity.TaskCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	id := "task-new"
	r.created = append(r.created, taskCreate)
	r.tasks = append(r.tasks, entity.Task{
		ID:          id,
		TypeID:      taskCreate.TypeID,
		ColumnID:    taskCreate.ColumnID,
		Label:       taskCreate.Label,
		Description: taskCreate.Description,
	})
	return id, nil
}

func (r *fakeTaskRepo) Update(_ context.Context, taskUpdate entity.TaskUpdate) error {
	if r.err != nil {
		return r.err
	}
	r.updated = append(r.updated, taskUpdate)
	for i := range r.tasks {
		if r.tasks[i].ID != taskUpdate.ID {
			continue
		}
		if taskUpdate.TypeID != nil {
			r.tasks[i].TypeID = *taskUpdate.TypeID
		}
		if taskUpdate.ColumnID != nil {
			r.tasks[i].ColumnID = *taskUpdate.ColumnID
		}
		if taskUpdate.Label != nil {
			r.tasks[i].Label = *taskUpdate.Label
		}
		if taskUpdate.Description != nil {
			r.tasks[i].Description = *taskUpdate.Description
		}
	}
	return nil
}

func (r *fakeTaskRepo) Delete(_ context.Context, taskID string) error {
	if r.err != nil {
		return r.err
	}
	r.deleted = append(r.deleted, taskID)
	return nil
}

type fakeMoveEventRepo struct {
	created []entity.MoveEventCreate
	err     error
}

func (r *fakeMoveEventRepo) Create(_ context.Context, moveEventCreate entity.MoveEventCreate) error {
	if r.err != nil {
		return r.err
	}
	r.created = append(r.created, moveEventCreate)
	return nil
}

func newTestTaskService(taskRepo *fakeTaskRepo) *TaskService {
	return New(
		&fakeBoardRepo{boards: entity.Boards{{ID: "board-1", UserID: "user-1", Name: "board"}}},
		&fakeColumnRepo{columns: entity.Columns{
			{ID: "backlog", BoardID: "board-1", Name: "backlog", OrderNumber: -1},
			{ID: "todo", BoardID: "board-1", Name: "todo", OrderNumber: 1},
		}},
		&fakeTypeRepo{types: entity.Types{{ID: "null-type", UserID: "user-1", Name: "null"}}},
		taskRepo,
		&fakeMoveEventRepo{},
	).(*TaskService)
}

func TestTaskServiceGetReturnsNotFoundWhenEmpty(t *testing.T) {
	service := newTestTaskService(&fakeTaskRepo{})

	_, err := service.Get(context.Background(), "board-1", entity.GetTasksOpts{})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskServiceGetByColumnIDRejectsForeignColumn(t *testing.T) {
	service := newTestTaskService(&fakeTaskRepo{})

	_, err := service.GetByColumnID(context.Background(), "board-1", "foreign")
	if err == nil {
		t.Fatal("expected invalid column error")
	}
}

func TestTaskServiceCreateUsesBacklogAndNullTypeDefaults(t *testing.T) {
	taskRepo := &fakeTaskRepo{}
	service := newTestTaskService(taskRepo)

	created, err := service.Create(context.Background(), "board-1", entity.TaskCreate{Label: "Task"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != "task-new" || created.ColumnID != "backlog" || created.TypeID != "null-type" {
		t.Fatalf("unexpected created task: %+v", created)
	}
	if len(taskRepo.created) != 1 || taskRepo.created[0].ColumnID != "backlog" || taskRepo.created[0].TypeID != "null-type" {
		t.Fatalf("unexpected create payloads: %+v", taskRepo.created)
	}
}

func TestTaskServiceUpdateDeleteAndMove(t *testing.T) {
	taskRepo := &fakeTaskRepo{tasks: entity.Tasks{{ID: "task-1", ColumnID: "todo", TypeID: "null-type", Label: "Old"}}}
	moveEvents := &fakeMoveEventRepo{}
	service := New(
		&fakeBoardRepo{boards: entity.Boards{{ID: "board-1", UserID: "user-1"}}},
		&fakeColumnRepo{columns: entity.Columns{
			{ID: "backlog", BoardID: "board-1", Name: "backlog", OrderNumber: -1},
			{ID: "todo", BoardID: "board-1", Name: "todo", OrderNumber: 1},
		}},
		&fakeTypeRepo{types: entity.Types{{ID: "null-type", UserID: "user-1", Name: "null"}}},
		taskRepo,
		moveEvents,
	)

	label := "New"
	updated, err := service.Update(context.Background(), "board-1", entity.TaskUpdate{ID: "task-1", Label: &label})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Label != label {
		t.Fatalf("expected updated label %q, got %+v", label, updated)
	}

	moved, err := service.Move(context.Background(), "board-1", "task-1", "backlog")
	if err != nil {
		t.Fatalf("Move returned error: %v", err)
	}
	if moved.ColumnID != "backlog" {
		t.Fatalf("expected moved task in backlog, got %+v", moved)
	}
	if len(moveEvents.created) != 1 || moveEvents.created[0].FromColumnID != "todo" || moveEvents.created[0].ToColumnID != "backlog" {
		t.Fatalf("unexpected move events: %+v", moveEvents.created)
	}

	if err := service.Delete(context.Background(), "task-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(taskRepo.deleted) != 1 || taskRepo.deleted[0] != "task-1" {
		t.Fatalf("unexpected deletes: %+v", taskRepo.deleted)
	}
}
