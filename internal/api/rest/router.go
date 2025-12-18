package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sater-151/todo-list/internal/handlers"
	"gl.iteco.com/technology/go_services/iteco_notifgatev2/whatsapp/internal/pkg/errorspkg"
	"gl.iteco.com/technology/go_services/toolbox/validate"
)

const (
	PathMetrics = "/metrics"

	PathWebhookWhatsApp = "/api/webhooks/whatsapp"
)

type (
	OnemsgWebhook interface {
		Handle(w http.ResponseWriter, r *http.Request)
	}
)

type (
	RouterDependencies struct {
		TodoTask
	}
)

func NewRouter(d *RouterDependencies) (http.Handler, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewRouter", d, err)
	}

	r := mux.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Get("/api/tasks", handlers.Auth(handlers.ListTask(service)))
	r.Get("/api/task", handlers.Auth(handlers.GetTask(db)))

	r.Post("/api/task", handlers.Auth(handlers.PostTask(service)))
	r.Post("/api/task/done", handlers.Auth(handlers.PostTaskDone(service)))
	r.Post("/api/signin", handlers.Sign)

	r.Put("/api/task", handlers.Auth(handlers.PutTask(service)))

	r.Delete("/api/task", handlers.Auth(handlers.DeleteTask(db)))

	return r, nil
}
