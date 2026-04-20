package column

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	"github.com/sater-151/todo-list/internal/service"
)

type ColumnController struct {
	s *service.TodoList
}

func Router(s *service.TodoList, router *gin.RouterGroup) {
	columnsRouter := router.Group("/boards/:board/columns")

	ctrl := &ColumnController{
		s: s,
	}

	columnsRouter.GET("/", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.LIST)
	columnsRouter.GET("/:column", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.GET)
	columnsRouter.POST("/", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.POST)
	columnsRouter.PATCH("/:column", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.PATCH)
	columnsRouter.DELETE("/:column", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.DELETE)
	columnsRouter.PUT("/swap", middlewares.CheckAuth(s.UserService), middlewares.CheckBoard(s.BoardService), ctrl.SWAP)
}
