package board

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	"github.com/sater-151/todo-list/internal/service"
)

type BoardController struct {
	s *service.TodoList
}

func Router(s *service.TodoList, router *gin.RouterGroup) {
	boardsRouter := router.Group("/boards")

	ctrl := BoardController{
		s: s,
	}

	boardsRouter.POST("/", middlewares.CheckAuth(s.UserService), ctrl.POST)
	boardsRouter.GET("/", middlewares.CheckAuth(s.UserService), ctrl.LIST)
	boardsRouter.GET("/:board", middlewares.CheckAuth(s.UserService), ctrl.GET)
	boardsRouter.PATCH("/:board", middlewares.CheckAuth(s.UserService), ctrl.PATCH)
	boardsRouter.DELETE("/:board", middlewares.CheckAuth(s.UserService), ctrl.DELETE)
}
