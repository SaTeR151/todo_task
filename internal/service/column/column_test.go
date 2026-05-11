package column

import (
	"context"
	"errors"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

type fakeColumnRepo struct {
	columns entity.Columns
	created []entity.ColumnCreate
	updated []entity.ColumnUpdate
	deleted []string
	err     error
}

func (r *fakeColumnRepo) GetColumns(_ context.Context, opts entity.GetColumnsOpts) (entity.Columns, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Columns
	for _, c := range r.columns {
		if opts.ID != "" && c.ID != opts.ID {
			continue
		}
		if opts.BoardID != "" && c.BoardID != opts.BoardID {
			continue
		}
		if opts.Name != "" && c.Name != opts.Name {
			continue
		}
		if opts.OrderNumber != 0 && c.OrderNumber != opts.OrderNumber {
			continue
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *fakeColumnRepo) CreateColumn(_ context.Context, columnCreate entity.ColumnCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	id := "column-new"
	r.created = append(r.created, columnCreate)
	r.columns = append(r.columns, entity.Column{
		ID:          id,
		BoardID:     columnCreate.BoardID,
		Name:        columnCreate.Name,
		OrderNumber: columnCreate.OderNumber,
	})
	return id, nil
}

func (r *fakeColumnRepo) UpdateColumn(_ context.Context, columnUpdate entity.ColumnUpdate) error {
	if r.err != nil {
		return r.err
	}
	r.updated = append(r.updated, columnUpdate)
	for i := range r.columns {
		if r.columns[i].ID != columnUpdate.ID {
			continue
		}
		if columnUpdate.Name != nil {
			r.columns[i].Name = *columnUpdate.Name
		}
		if columnUpdate.OrderNumber != nil {
			r.columns[i].OrderNumber = *columnUpdate.OrderNumber
		}
	}
	return nil
}

func (r *fakeColumnRepo) DeleteColumn(_ context.Context, columnID string) error {
	if r.err != nil {
		return r.err
	}
	r.deleted = append(r.deleted, columnID)
	return nil
}

type fakeTaskRepo struct {
	tasks entity.Tasks
	err   error
}

func (r *fakeTaskRepo) Get(_ context.Context, opts entity.GetTasksOpts) (entity.Tasks, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Tasks
	for _, task := range r.tasks {
		if opts.ColumnID != "" && task.ColumnID != opts.ColumnID {
			continue
		}
		res = append(res, task)
	}
	return res, nil
}

type fakeTaskMover struct {
	moves []struct {
		boardID     string
		taskID      string
		newColumnID string
	}
	err error
}

func (m *fakeTaskMover) Get(context.Context, string, entity.GetTasksOpts) (entity.Tasks, error) {
	return nil, nil
}
func (m *fakeTaskMover) GetByID(context.Context, string, string) (entity.Task, error) {
	return entity.Task{}, nil
}
func (m *fakeTaskMover) GetByColumnID(context.Context, string, string) (entity.Tasks, error) {
	return nil, nil
}
func (m *fakeTaskMover) GetByTypeID(context.Context, string, string) (entity.Tasks, error) {
	return nil, nil
}
func (m *fakeTaskMover) GetByBoardID(context.Context, string) (entity.Tasks, error) {
	return nil, nil
}
func (m *fakeTaskMover) Create(context.Context, string, entity.TaskCreate) (entity.Task, error) {
	return entity.Task{}, nil
}
func (m *fakeTaskMover) Update(context.Context, string, entity.TaskUpdate) (entity.Task, error) {
	return entity.Task{}, nil
}
func (m *fakeTaskMover) Delete(context.Context, string) error { return nil }
func (m *fakeTaskMover) Move(_ context.Context, boardID, taskID, newColumnID string) (entity.Task, error) {
	if m.err != nil {
		return entity.Task{}, m.err
	}
	m.moves = append(m.moves, struct {
		boardID     string
		taskID      string
		newColumnID string
	}{boardID: boardID, taskID: taskID, newColumnID: newColumnID})
	return entity.Task{ID: taskID, ColumnID: newColumnID}, nil
}

func TestColumnServiceGetReturnsNotFoundWhenEmpty(t *testing.T) {
	service := New(&fakeColumnRepo{}, &fakeTaskRepo{}, &fakeTaskMover{})

	_, err := service.Get(context.Background(), entity.GetColumnsOpts{BoardID: "missing"})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestColumnServiceCreateAndUpdate(t *testing.T) {
	repo := &fakeColumnRepo{}
	service := New(repo, &fakeTaskRepo{}, &fakeTaskMover{})

	created, err := service.CreateColumn(context.Background(), entity.ColumnCreate{BoardID: "board-1", Name: "todo", OderNumber: 1})
	if err != nil {
		t.Fatalf("CreateColumn returned error: %v", err)
	}
	if created.ID != "column-new" || created.Name != "todo" {
		t.Fatalf("unexpected created column: %+v", created)
	}

	name := "done"
	order := 2
	updated, err := service.UpdateColumn(context.Background(), "board-1", entity.ColumnUpdate{ID: "column-new", Name: &name, OrderNumber: &order})
	if err != nil {
		t.Fatalf("UpdateColumn returned error: %v", err)
	}
	if updated.Name != name || updated.OrderNumber != order {
		t.Fatalf("unexpected updated column: %+v", updated)
	}
}

func TestColumnServiceDeleteMovesTasksToBacklog(t *testing.T) {
	repo := &fakeColumnRepo{columns: entity.Columns{
		{ID: "backlog", BoardID: "board-1", Name: "backlog", OrderNumber: -1},
		{ID: "todo", BoardID: "board-1", Name: "todo", OrderNumber: 1},
	}}
	tasks := &fakeTaskRepo{tasks: entity.Tasks{{ID: "task-1", ColumnID: "todo"}, {ID: "task-2", ColumnID: "todo"}}}
	mover := &fakeTaskMover{}
	service := New(repo, tasks, mover)

	if err := service.DeleteColumn(context.Background(), "board-1", "todo"); err != nil {
		t.Fatalf("DeleteColumn returned error: %v", err)
	}
	if len(mover.moves) != 2 {
		t.Fatalf("expected two task moves, got %+v", mover.moves)
	}
	for _, move := range mover.moves {
		if move.newColumnID != "backlog" {
			t.Fatalf("expected move to backlog, got %+v", move)
		}
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "todo" {
		t.Fatalf("unexpected deletes: %+v", repo.deleted)
	}
}

func TestColumnServiceSwapColumns(t *testing.T) {
	repo := &fakeColumnRepo{columns: entity.Columns{
		{ID: "a", BoardID: "board-1", Name: "a", OrderNumber: 1},
		{ID: "b", BoardID: "board-1", Name: "b", OrderNumber: 2},
	}}
	service := New(repo, &fakeTaskRepo{}, &fakeTaskMover{})

	if err := service.SwapColumns(context.Background(), "board-1", "a", "b"); err != nil {
		t.Fatalf("SwapColumns returned error: %v", err)
	}
	if len(repo.updated) != 2 {
		t.Fatalf("expected two updates, got %+v", repo.updated)
	}
	if repo.columns[0].OrderNumber != 2 || repo.columns[1].OrderNumber != 1 {
		t.Fatalf("columns were not swapped: %+v", repo.columns)
	}
}
