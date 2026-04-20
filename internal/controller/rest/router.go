package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/board"
	"github.com/sater-151/todo-list/internal/controller/rest/column"
	"github.com/sater-151/todo-list/internal/controller/rest/task"
	tasktype "github.com/sater-151/todo-list/internal/controller/rest/type"
	"github.com/sater-151/todo-list/internal/controller/rest/user"
	"github.com/sater-151/todo-list/internal/service"
)

type Rest struct {
	s *service.TodoList
}

func New(s *service.TodoList) *Rest {
	return &Rest{
		s: s,
	}
}

func (r *Rest) Run() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	apiGroup := router.Group("/api")

	user.Router(r.s, apiGroup)
	board.Router(r.s, apiGroup)
	column.Router(r.s, apiGroup)
	task.Router(r.s, apiGroup)
	tasktype.Router(r.s, apiGroup)

	return router
}
