package tasktype

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	"github.com/sater-151/todo-list/internal/service"
)

type TypeController struct {
	s *service.TodoList
}

func Router(s *service.TodoList, router *gin.RouterGroup) {
	ctrl := TypeController{s}

	typeRouter := router.Group("/types")
	{
		typeRouter.POST("/", middlewares.CheckAuth(s.UserService), ctrl.POST)
		typeRouter.GET("/", middlewares.CheckAuth(s.UserService), ctrl.LIST)
		typeRouter.GET("/:type", middlewares.CheckAuth(s.UserService), ctrl.GET)
		typeRouter.PATCH("/:type", middlewares.CheckAuth(s.UserService), ctrl.PATCH)
		typeRouter.DELETE("/:type", middlewares.CheckAuth(s.UserService), ctrl.DELETE)
	}
}
