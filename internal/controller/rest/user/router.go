package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	userservice "github.com/sater-151/todo-list/internal/service/user"
)

type UserController struct {
	users userservice.User
}

func Router(users userservice.User, router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	ctrl := UserController{
		users: users,
	}

	userRouter.GET("/", middlewares.CheckAuth(users), ctrl.Get)
	userRouter.POST("/", ctrl.POST)
	userRouter.PATCH("/password-change", middlewares.CheckAuth(users), ctrl.ChangePassword)
	userRouter.DELETE("/", middlewares.CheckAuth(users), ctrl.DELETE)
	userRouter.POST("/auth", ctrl.Auth)
	userRouter.POST("/refresh", middlewares.CheckAuth(users), ctrl.RefreshToken)
	userRouter.POST("/logout", middlewares.CheckAuth(users), ctrl.LogOut)
}
