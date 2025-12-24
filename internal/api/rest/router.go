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
	PathAPI      = "/api"
	PathTasks    = "/tasks"
	PathTask     = "/task"
	PathTaskDone = "/task/done"
	PathSignin   = "/signin"
)

type (
	ITodoTaskHandlers interface {
		PostTask(res http.ResponseWriter, req *http.Request)
		ListTask(res http.ResponseWriter, req *http.Request)
		GetTask(res http.ResponseWriter, req *http.Request)
		PutTask(res http.ResponseWriter, req *http.Request)
		PostTaskDone(res http.ResponseWriter, req *http.Request)
		DeleteTask(res http.ResponseWriter, req *http.Request)
		Sign(res http.ResponseWriter, req *http.Request)
	}

	IInternalMW interface {
		Auth(n http.Handler) http.Handler
	}
)

type (
	RouterDependencies struct {
		Handlers   ITodoTaskHandlers
		InternalMW IInternalMW
	}
)

func NewRouter(d *RouterDependencies) (http.Handler, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewRouter", d, err)
	}

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	// --- API ---
	apiR := r.Route(PathAPI, func(r chi.Router) {})

	apiR.Post(PathSignin, d.Handlers.Sign)

	authR := apiR.With(d.InternalMW.Auth)

	authR.Get(PathTasks, d.Handlers.ListTask)
	authR.Get(PathTask, d.Handlers.GetTask)

	authR.Post(PathTask, d.Handlers.PostTask)
	authR.Post(PathTaskDone, d.Handlers.PostTaskDone)

	authR.Put(PathTask, d.Handlers.PutTask)
	authR.Delete(PathTask, d.Handlers.DeleteTask)

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	return r, nil
}
