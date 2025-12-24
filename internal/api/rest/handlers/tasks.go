package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/todo-list/internal/models"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/sater-151/todo-list/internal/utils/selectconfig"
)

type (
	ITodoTaskUsecase interface {
		AddTask(ctx context.Context, task *models.Task) (string, error)
		GetListTask(ctx context.Context, selectConfig *models.SelectConfig) ([]models.Task, error)
		TaskDone(ctx context.Context, selectConfig *models.SelectConfig) error
		UpdateTask(ctx context.Context, task *models.Task) error
		DeleteTask(ctx context.Context, uuid string) error
		Select(ctx context.Context, selectConfig *models.SelectConfig) ([]models.Task, error)
	}
)

type TodoTaskServerDependencies struct {
	TodoTaskUsecase ITodoTaskUsecase `validate:"required"`
	Password        string           `validate:"required"`
}

type TodoTaskServer struct {
	todoTaskUsecase ITodoTaskUsecase
	password        string
}

func NewTodoTaskHandlers(d *TodoTaskServerDependencies) (*TodoTaskServer, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewTodoTaskHandlers", d, err)
	}

	return &TodoTaskServer{
		todoTaskUsecase: d.TodoTaskUsecase,
		password:        d.Password,
	}, nil
}

func (s *TodoTaskServer) PostTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")

	var task models.Task
	if err := sonic.ConfigDefault.NewDecoder(req.Body).Decode(&task); err != nil {
		slog.Warn(err.Error())
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.todoTaskUsecase.AddTask(req.Context(), &task)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(res).Encode(id); err != nil {
		slog.Error(err.Error())
		http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusBadRequest)
		return
	}

}

func (s *TodoTaskServer) ListTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")

	search := req.FormValue("search")
	selectConfig := selectconfig.Default()
	if search != "" {
		selectConfig.Search = search
	}

	tasks, err := s.todoTaskUsecase.GetListTask(req.Context(), selectConfig)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(res).Encode(models.ListTask{Tasks: tasks}); err != nil {
		slog.Error(err.Error())
		http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *TodoTaskServer) GetTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	id := req.FormValue("id")
	if id == "" {
		slog.Warn("id required")
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	selectConfig := selectconfig.Default()
	selectConfig.ID = id

	tasks, err := s.todoTaskUsecase.Select(req.Context(), selectConfig)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	if len(tasks) == 0 {
		http.Error(res, "task not found", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(tasks[0]); err != nil {
		slog.Error(err.Error())
		http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *TodoTaskServer) PutTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")

	var task models.Task
	if err := sonic.ConfigDefault.NewDecoder(req.Body).Decode(&task); err != nil {
		slog.Warn(err.Error())
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	err := s.todoTaskUsecase.UpdateTask(req.Context(), &task)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (s *TodoTaskServer) PostTaskDone(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")

	id := req.FormValue("id")
	selectConfig := selectconfig.Default()
	if id == "" {
		slog.Warn("id required")
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	selectConfig.ID = id

	err := s.todoTaskUsecase.TaskDone(req.Context(), selectConfig)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (s *TodoTaskServer) DeleteTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	id := req.FormValue("id")
	if id == "" {
		slog.Warn("id required")
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	err := s.todoTaskUsecase.DeleteTask(req.Context(), id)
	if err != nil {
		if errors.Is(err, errorspkg.ErrBadRequest) {
			http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (s *TodoTaskServer) Sign(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	passTrue := s.password

	var passJS models.PasswordJS
	if err := sonic.ConfigDefault.NewDecoder(req.Body).Decode(&passJS); err != nil {
		slog.Warn("password required")
		http.Error(res, "password required", http.StatusBadRequest)
		return
	}

	if passTrue != passJS.Pass {
		slog.Warn("wrong password")
		http.Error(res, "wrong password", http.StatusBadRequest)
		return
	}

	token, err := jwt.New(jwt.SigningMethodHS256).SignedString([]byte(passTrue))
	if err != nil {
		slog.Error(err.Error())
		http.Error(res, errorspkg.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(models.JWTToken{Token: token}); err != nil {
		slog.Error(err.Error())
		http.Error(res, errorspkg.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}
