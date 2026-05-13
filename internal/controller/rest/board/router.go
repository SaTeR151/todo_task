package board

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	boardservice "github.com/sater-151/todo-list/internal/service/board"
	userservice "github.com/sater-151/todo-list/internal/service/user"
)

type BoardController struct {
	boards boardservice.Board
}

func Router(boards boardservice.Board, users userservice.User, router *gin.RouterGroup) {
	boardsRouter := router.Group("/boards")

	ctrl := BoardController{
		boards: boards,
	}

	boardsRouter.POST("/", middlewares.CheckAuth(users), ctrl.POST)
	boardsRouter.GET("/", middlewares.CheckAuth(users), ctrl.LIST)
	boardsRouter.GET("/:board", middlewares.CheckAuth(users), ctrl.GET)
	boardsRouter.PATCH("/:board", middlewares.CheckAuth(users), ctrl.PATCH)
	boardsRouter.DELETE("/:board", middlewares.CheckAuth(users), ctrl.DELETE)
}
