package task

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	"github.com/sater-151/todo-list/internal/service"
)

type TaskController struct {
	s *service.TodoList
}

func Router(s *service.TodoList, router *gin.RouterGroup) {
	ctrl := TaskController{s: s}

	taskRouter := router.Group("/boards/:board/tasks")

	taskRouter.GET("/", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.LIST)
	taskRouter.GET("/:task", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.GET)
	taskRouter.POST("/", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.POST)
	taskRouter.PATCH("/:task", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.PATCH)
	taskRouter.DELETE("/:task", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.DELETE)
	taskRouter.PUT("/:task/move", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.MOVE)
}
