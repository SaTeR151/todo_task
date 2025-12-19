package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/sater-151/todo-list/internal/models"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/sater-151/todo-list/internal/utils"
)

type (
	ITodoTaskRepo interface {
		InsertTask(ctx context.Context, task models.Task) (string, error)
		UpdateTask(ctx context.Context, task models.Task) error
		DeleteTask(ctx context.Context, uuid string) error
		Select(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error)
	}
)

type (
	TodoTaskDependencies struct {
		TodoTaskRepo ITodoTaskRepo
	}

	TodoTask struct {
		todoTaskRepo ITodoTaskRepo
	}
)

func NewTodoTask(d *TodoTaskDependencies) (*TodoTask, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("usecases.NewTodoTask", d, err)
	}

	return &TodoTask{
		todoTaskRepo: d.TodoTaskRepo,
	}, nil
}

func (s *TodoTask) AddTask(ctx context.Context, task models.Task) (models.ID, error) {
	var Id models.ID
	var err error
	task, err = utils.CheckTask(task)
	if err != nil {
		return Id, err
	}
	Id.ID, err = s.todoTaskRepo.InsertTask(ctx, task)
	if err != nil {
		return Id, err
	}
	return Id, nil
}

func (s *TodoTask) UpdateTask(ctx context.Context, task models.Task) error {
	task, err := utils.CheckTask(task)
	if err != nil {
		return err
	}
	err = s.todoTaskRepo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoTask) TaskDone(ctx context.Context, selectConfig models.SelectConfig) error {
	tasks, err := s.todoTaskRepo.Select(ctx, selectConfig)
	if err != nil {
		return err
	}
	task := tasks[0]
	if task.Repeat == "" {
		err = s.todoTaskRepo.DeleteTask(ctx, selectConfig.Id)
		if err != nil {
			return err
		}
		return nil
	}
	task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}
	err = s.todoTaskRepo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoTask) GetListTask(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error) {
	if selectConfig.Search != "" {
		date := strings.Split(selectConfig.Search, ".")
		if len(date) == 3 {
			var d string
			for i := 2; i >= 0; i-- {
				d += date[i]
			}
			_, err := time.Parse("20060102", d)
			if err == nil {
				selectConfig.Search = ""
				selectConfig.Date = d
			}
		}
	}
	listTask, err := s.todoTaskRepo.Select(ctx, selectConfig)
	if err != nil {
		return listTask, err
	}
	return listTask, err
}
