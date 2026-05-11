package task

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	boardservice "github.com/sater-151/todo-list/internal/service/board"
	columnservice "github.com/sater-151/todo-list/internal/service/column"
	taskservice "github.com/sater-151/todo-list/internal/service/task"
	typeservice "github.com/sater-151/todo-list/internal/service/type"
	userservice "github.com/sater-151/todo-list/internal/service/user"
)

type TaskController struct {
	tasks   taskservice.Task
	columns columnservice.Column
	types   typeservice.Type
}

func Router(
	tasks taskservice.Task,
	columns columnservice.Column,
	types typeservice.Type,
	boards boardservice.Board,
	users userservice.User,
	router *gin.RouterGroup,
) {
	ctrl := TaskController{
		tasks:   tasks,
		columns: columns,
		types:   types,
	}

	taskRouter := router.Group("/boards/:board/tasks")

	taskRouter.GET("/", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.LIST)
	taskRouter.GET("/:task", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.GET)
	taskRouter.POST("/", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.POST)
	taskRouter.PATCH("/:task", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.PATCH)
	taskRouter.DELETE("/:task", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.DELETE)
	taskRouter.PUT("/:task/move", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.MOVE)
}
