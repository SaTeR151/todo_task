package task

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
)

type BoardRepository interface {
	Get(ctx context.Context, opts entity.GetBoardsOpts) (entity.Boards, error)
}

type ColumnRepository interface {
	GetColumns(ctx context.Context, opts entity.GetColumnsOpts) (entity.Columns, error)
}

type TypeRepository interface {
	Get(ctx context.Context, opts entity.GetTypesOpts) (entity.Types, error)
}

type Repository interface {
	Get(ctx context.Context, opts entity.GetTasksOpts) (entity.Tasks, error)
	Create(ctx context.Context, taskCreate entity.TaskCreate) (string, error)
	Update(ctx context.Context, taskUpdate entity.TaskUpdate) error
	Delete(ctx context.Context, taskID string) error
}

type MoveEventRepository interface {
	Create(ctx context.Context, moveEventCreate entity.MoveEventCreate) error
}

type Task interface {
	Get(ctx context.Context, boardID string, opts entity.GetTasksOpts) (res entity.Tasks, err error)
	GetByID(ctx context.Context, boardID, taskID string) (res entity.Task, err error)
	GetByColumnID(ctx context.Context, boardID, columnID string) (res entity.Tasks, err error)
	GetByTypeID(ctx context.Context, boardID, typeID string) (res entity.Tasks, err error)
	GetByBoardID(ctx context.Context, boardID string) (res entity.Tasks, err error)
	Create(ctx context.Context, boardID string, taskCreate entity.TaskCreate) (res entity.Task, err error)
	Update(ctx context.Context, boardID string, taskUpdate entity.TaskUpdate) (res entity.Task, err error)
	Delete(ctx context.Context, taskID string) (err error)
	Move(ctx context.Context, boardID, taskID, newColumnID string) (res entity.Task, err error)
}

func New(
	boardRepo BoardRepository,
	columnRepo ColumnRepository,
	typeRepo TypeRepository,
	taskRepo Repository,
	moveEventRepo MoveEventRepository,
) Task {
	return &TaskService{
		boards:     boardRepo,
		columns:    columnRepo,
		types:      typeRepo,
		tasks:      taskRepo,
		moveEvents: moveEventRepo,
	}
}
