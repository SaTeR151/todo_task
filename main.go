package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sater-151/todo-list/internal/config"
	"github.com/sater-151/todo-list/internal/database"
	"github.com/sater-151/todo-list/internal/handlers"
	"github.com/sater-151/todo-list/internal/service"
)

func main() {
	config := config.GetConfig()

	db, err := database.OpenDB(config.DbFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer db.Close()

	service := service.New(db)

	r := chi.NewRouter()

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

	log.Println("Server start at port:", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Println("Ошибка запуска сервера:", err.Error())
		return
	}
}
