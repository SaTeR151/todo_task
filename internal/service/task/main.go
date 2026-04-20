package task

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

type Task interface {
	Get(ctx context.Context, opts entity.GetTasksOpts) (res entity.Tasks, err error)
	GetByID(ctx context.Context, boardID, taskID string) (res entity.Task, err error)
	GetByColumnID(ctx context.Context, boardID, columnID string) (res entity.Tasks, err error)
	GetByTypeID(ctx context.Context, typeID string) (res entity.Tasks, err error)
	GetByBoardID(ctx context.Context, boardID string) (res entity.Tasks, err error)
	Create(ctx context.Context, boardID string, taskCreate entity.TaskCreate) (res entity.Task, err error)
	Update(ctx context.Context, boardID string, taskUpdate entity.TaskUpdate) (res entity.Task, err error)
	Delete(ctx context.Context, taskID string) (err error)
	Move(ctx context.Context, boardID, taskID, newColumnID string) (res entity.Task, err error)
}

func New(repo *postgres.Repository) Task {
	return &TaskService{
		repo: repo,
	}
}
