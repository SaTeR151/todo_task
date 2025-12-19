package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
)

const (
	webDir       = "web"
	PathApi      = "/api"
	PathNextDate = "/nextdate"
	PathTasks    = "/tasks"
	PathTask     = "/task"
	PathTaskDone = "/task/done"
	PathSignin   = "/signin"
)

type (
	ITodoTaskHandlers interface {
		GetNextDate(res http.ResponseWriter, req *http.Request)
		PostTask(res http.ResponseWriter, req *http.Request)
		ListTask(res http.ResponseWriter, req *http.Request)
		GetTask(res http.ResponseWriter, req *http.Request)
		PutTask(res http.ResponseWriter, req *http.Request)
		PostTaskDone(res http.ResponseWriter, req *http.Request)
		DeleteTask(res http.ResponseWriter, req *http.Request)
		Sign(res http.ResponseWriter, req *http.Request)
		Auth(n http.HandlerFunc) http.HandlerFunc
	}
)

type (
	RouterDependencies struct {
		Handlers ITodoTaskHandlers
	}
)

func NewRouter(d *RouterDependencies) (http.Handler, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewRouter", d, err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// --- API ---
	r.Route(PathApi, func(r chi.Router) {
		r.Get(PathNextDate, d.Handlers.GetNextDate)
		r.Get(PathTasks, d.Handlers.Auth(d.Handlers.ListTask))
		r.Get(PathTask, d.Handlers.Auth(d.Handlers.GetTask))

		r.Post(PathTask, d.Handlers.Auth(d.Handlers.PostTask))
		r.Post(PathTaskDone, d.Handlers.Auth(d.Handlers.PostTaskDone))
		r.Post(PathSignin, d.Handlers.Sign)

		r.Put(PathTask, d.Handlers.Auth(d.Handlers.PutTask))
		r.Delete(PathTask, d.Handlers.Auth(d.Handlers.DeleteTask))
	})

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	return r, nil
}
