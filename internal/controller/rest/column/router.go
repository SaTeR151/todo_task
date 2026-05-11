package column

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	boardservice "github.com/sater-151/todo-list/internal/service/board"
	columnservice "github.com/sater-151/todo-list/internal/service/column"
	userservice "github.com/sater-151/todo-list/internal/service/user"
)

type ColumnController struct {
	columns columnservice.Column
}

func Router(columns columnservice.Column, boards boardservice.Board, users userservice.User, router *gin.RouterGroup) {
	columnsRouter := router.Group("/boards/:board/columns")

	ctrl := &ColumnController{
		columns: columns,
	}

	columnsRouter.GET("/", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.LIST)
	columnsRouter.GET("/:column", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.GET)
	columnsRouter.POST("/", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.POST)
	columnsRouter.PATCH("/:column", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.PATCH)
	columnsRouter.DELETE("/:column", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.DELETE)
	columnsRouter.PUT("/swap", middlewares.CheckAuth(users), middlewares.CheckBoard(boards), ctrl.SWAP)
}
