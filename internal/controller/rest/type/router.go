package tasktype

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	typeservice "github.com/sater-151/todo-list/internal/service/type"
	userservice "github.com/sater-151/todo-list/internal/service/user"
)

type TypeController struct {
	types typeservice.Type
}

func Router(types typeservice.Type, users userservice.User, router *gin.RouterGroup) {
	ctrl := TypeController{types: types}

	typeRouter := router.Group("/types")
	{
		typeRouter.POST("/", middlewares.CheckAuth(users), ctrl.POST)
		typeRouter.GET("/", middlewares.CheckAuth(users), ctrl.LIST)
		typeRouter.GET("/:type", middlewares.CheckAuth(users), ctrl.GET)
		typeRouter.PATCH("/:type", middlewares.CheckAuth(users), ctrl.PATCH)
		typeRouter.DELETE("/:type", middlewares.CheckAuth(users), ctrl.DELETE)
	}
}
