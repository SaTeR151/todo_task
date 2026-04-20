package service

import (
	"github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/internal/service/board"
	"github.com/sater-151/todo-list/internal/service/column"
	"github.com/sater-151/todo-list/internal/service/task"
	type_service "github.com/sater-151/todo-list/internal/service/type"
	"github.com/sater-151/todo-list/internal/service/user"
)

type TodoList struct {
	BoardService  board.Board
	ColumnService column.Column
	TaskService   task.Task
	TypeService   type_service.Type
	UserService   user.User
}

func New(repo *postgres.Repository, secretKey string) *TodoList {
	return &TodoList{
		BoardService:  board.New(repo),
		ColumnService: column.New(repo),
		TaskService:   task.New(repo),
		TypeService:   type_service.New(repo),
		UserService:   user.New(repo, secretKey),
	}
}
