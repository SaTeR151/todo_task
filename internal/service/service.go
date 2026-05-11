package service

import (
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

type Repositories struct {
	Board     board.BoardRepository
	Column    column.Repository
	Type      type_service.Repository
	Task      task.Repository
	User      user.Repository
	MoveEvent task.MoveEventRepository
}

func New(repo Repositories, secretKey string) *TodoList {
	boardService := board.New(repo.Board, repo.Column)
	taskService := task.New(repo.Board, repo.Column, repo.Type, repo.Task, repo.MoveEvent)
	columnService := column.New(repo.Column, repo.Task, taskService)
	userService := user.New(repo.User, repo.Type, secretKey)
	typeService := type_service.New(repo.Type)
	return &TodoList{
		BoardService:  boardService,
		ColumnService: columnService,
		TaskService:   taskService,
		TypeService:   typeService,
		UserService:   userService,
	}
}
