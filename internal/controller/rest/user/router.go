package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/middlewares"
	"github.com/sater-151/todo-list/internal/service"
)

type UserController struct {
	s *service.TodoList
}

func Router(s *service.TodoList, router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	ctrl := UserController{
		s: s,
	}

	userRouter.GET("/", middlewares.CheckAuth(s.UserService), ctrl.Get)
	userRouter.POST("/", ctrl.POST)
	userRouter.PATCH("/password-change", middlewares.CheckAuth(s.UserService), ctrl.ChangePassword)
	userRouter.DELETE("/", middlewares.CheckAuth(s.UserService), ctrl.DELETE)
	userRouter.POST("/auth", ctrl.Auth)
	userRouter.POST("/refresh", middlewares.CheckAuth(s.UserService), ctrl.RefreshToken)
	userRouter.POST("/logout", middlewares.CheckAuth(s.UserService), ctrl.LogOut)
}
