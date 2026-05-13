package board

import (
	"context"
	"errors"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

type fakeBoardRepo struct {
	boards  entity.Boards
	created []entity.BoardCreate
	updated []entity.BoardUpdate
	deleted []string
	err     error
}

func (r *fakeBoardRepo) Get(_ context.Context, opts entity.GetBoardsOpts) (entity.Boards, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Boards
	for _, b := range r.boards {
		if opts.ID != "" && b.ID != opts.ID {
			continue
		}
		if opts.UserID != "" && b.UserID != opts.UserID {
			continue
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *fakeBoardRepo) Create(_ context.Context, boardCreate entity.BoardCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	r.created = append(r.created, boardCreate)
	id := "board-new"
	r.boards = append(r.boards, entity.Board{ID: id, UserID: boardCreate.UserID, Name: boardCreate.Name})
	return id, nil
}

func (r *fakeBoardRepo) Update(_ context.Context, boardUpdate entity.BoardUpdate) error {
	if r.err != nil {
		return r.err
	}
	r.updated = append(r.updated, boardUpdate)
	for i := range r.boards {
		if r.boards[i].ID == boardUpdate.ID && boardUpdate.Name != nil {
			r.boards[i].Name = *boardUpdate.Name
		}
	}
	return nil
}

func (r *fakeBoardRepo) Delete(_ context.Context, boardID string) error {
	if r.err != nil {
		return r.err
	}
	r.deleted = append(r.deleted, boardID)
	return nil
}

type fakeColumnCreator struct {
	created []entity.ColumnCreate
	err     error
}

func (r *fakeColumnCreator) CreateColumn(_ context.Context, columnCreate entity.ColumnCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	r.created = append(r.created, columnCreate)
	return "column-new", nil
}

func TestBoardServiceGetReturnsNotFoundWhenEmpty(t *testing.T) {
	service := New(&fakeBoardRepo{}, &fakeColumnCreator{})

	_, err := service.Get(context.Background(), entity.GetBoardsOpts{UserID: "missing"})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBoardServiceCreateCreatesBacklogAndReturnsBoard(t *testing.T) {
	boards := &fakeBoardRepo{}
	columns := &fakeColumnCreator{}
	service := New(boards, columns)

	got, err := service.Create(context.Background(), entity.BoardCreate{UserID: "user-1", Name: "Work"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got.ID != "board-new" || got.UserID != "user-1" || got.Name != "Work" {
		t.Fatalf("unexpected board: %+v", got)
	}
	if len(columns.created) != 1 {
		t.Fatalf("expected backlog column to be created once, got %d", len(columns.created))
	}
	if columns.created[0].Name != "backlog" || columns.created[0].OderNumber != -1 || columns.created[0].BoardID != "board-new" {
		t.Fatalf("unexpected backlog column: %+v", columns.created[0])
	}
}

func TestBoardServiceUpdateAndDeleteCheckOwnership(t *testing.T) {
	name := "Updated"
	boards := &fakeBoardRepo{boards: entity.Boards{{ID: "board-1", UserID: "user-1", Name: "Old"}}}
	service := New(boards, &fakeColumnCreator{})

	updated, err := service.Update(context.Background(), "user-1", entity.BoardUpdate{ID: "board-1", Name: &name})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Name != name {
		t.Fatalf("expected updated name %q, got %q", name, updated.Name)
	}

	if err := service.Delete(context.Background(), "user-1", "board-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(boards.deleted) != 1 || boards.deleted[0] != "board-1" {
		t.Fatalf("unexpected deletes: %+v", boards.deleted)
	}
}
