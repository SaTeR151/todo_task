package app

import (
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/sater-151/todo-list/internal/usecases"
)

type (
	UsecasesDependencies struct {
		Repository *Repository `validate:"required"`
	}

	Usecases struct {
		TodoTask *usecases.TodoTask
	}
)

func NewUsecases(d *UsecasesDependencies) (*Usecases, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewUsecases", d, err)
	}

	todoTask, err := usecases.NewTodoTask(&usecases.TodoTaskDependencies{})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		TodoTask: todoTask,
	}, nil
}
