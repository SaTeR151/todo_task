package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/todo-list/internal/models"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/sater-151/todo-list/internal/utils"
)

type (
	ITodoTaskUsecase interface {
		AddTask(ctx context.Context, task models.Task) (models.ID, error)
		GetListTask(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error)
		TaskDone(ctx context.Context, selectConfig models.SelectConfig) error
	}

	ITodoTaskRepo interface {
		InsertTask(ctx context.Context, task models.Task) (string, error)
		UpdateTask(ctx context.Context, task models.Task) error
		DeleteTask(ctx context.Context, uuid string) error
		Select(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error)
	}
)

type TodoTaskServerDependencies struct {
	TodoTaskUsecase ITodoTaskUsecase `validate:"required"`
	TodoTaskRepo    ITodoTaskRepo    `validate:"required"`
	Password        string           `validate:"required"`
}

type TodoTaskServer struct {
	todoTaskUsecase ITodoTaskUsecase
	todoTaskRepo    ITodoTaskRepo
	password        string
}

func NewTodoTaskHandlers(d *TodoTaskServerDependencies) (*TodoTaskServer, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewTodoTaskHandlers", d, err)
	}

	return &TodoTaskServer{
		todoTaskUsecase: d.TodoTaskUsecase,
		todoTaskRepo:    d.TodoTaskRepo,
		password:        d.Password,
	}, nil
}

func ErrorHandler(res http.ResponseWriter, err error, status int) {
	var errJS models.Error
	errJS.Err = err.Error()
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(errJS)
}

func CreateDefaultSelectConfig() models.SelectConfig {
	var selectConfig models.SelectConfig
	selectConfig.Limit = "20"
	selectConfig.Sort = "date"
	selectConfig.Table = "scheduler"
	return selectConfig
}

func (s *TodoTaskServer) GetNextDate(res http.ResponseWriter, req *http.Request) {
	slog.Debug("1")
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	now := req.FormValue("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	nextDate, err := utils.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}

func (s *TodoTaskServer) PostTask(res http.ResponseWriter, req *http.Request) {
	slog.Debug("2")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var task models.Task
	var buf bytes.Buffer
	var idJS models.ID
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	idJS, err = s.todoTaskUsecase.AddTask(req.Context(), task)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(idJS)

}

func (s *TodoTaskServer) ListTask(res http.ResponseWriter, req *http.Request) {
	slog.Debug("3")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var err error
	search := req.FormValue("search")
	selectConfig := CreateDefaultSelectConfig()
	if search != "" {
		selectConfig.Search = search
	}
	tasks, err := s.todoTaskUsecase.GetListTask(req.Context(), selectConfig)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	listTask := models.ListTask{Tasks: tasks}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(listTask)
}

func (s *TodoTaskServer) GetTask(res http.ResponseWriter, req *http.Request) {
	slog.Debug("4")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	id := req.FormValue("id")
	if id == "" {
		log.Println("id required")
		ErrorHandler(res, fmt.Errorf("id required"), http.StatusBadRequest)
		return
	}
	selectConfig := CreateDefaultSelectConfig()
	selectConfig.Id = id
	tasks, err := s.todoTaskRepo.Select(req.Context(), selectConfig)
	if err != nil {
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}

	if len(tasks) == 0 {
		ErrorHandler(res, fmt.Errorf("not task"), http.StatusBadRequest)
		return
	}
	task := tasks[0]
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(task)

}

func (s *TodoTaskServer) PutTask(res http.ResponseWriter, req *http.Request) {
	slog.Debug("5")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var task models.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}

	err = s.todoTaskRepo.UpdateTask(req.Context(), task)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})
}

func (s *TodoTaskServer) PostTaskDone(res http.ResponseWriter, req *http.Request) {
	slog.Debug("6")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	id := req.FormValue("id")
	selectConfig := CreateDefaultSelectConfig()
	selectConfig.Id = id
	err := s.todoTaskUsecase.TaskDone(req.Context(), selectConfig)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})

}

func (s *TodoTaskServer) DeleteTask(res http.ResponseWriter, req *http.Request) {
	slog.Debug("7")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	id := req.FormValue("id")
	if id == "" {
		log.Println("id required")
		ErrorHandler(res, fmt.Errorf("id required"), http.StatusBadRequest)
		return
	}
	err := s.todoTaskRepo.DeleteTask(req.Context(), id)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})
}

func (s *TodoTaskServer) Sign(res http.ResponseWriter, req *http.Request) {
	slog.Debug("8")

	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	passTrue := s.password
	var passJS models.PasswordJS
	var token models.JWTToken
	err := json.NewDecoder(req.Body).Decode(&passJS)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	if passTrue != passJS.Pass {
		log.Println("wrong password")
		ErrorHandler(res, fmt.Errorf("wrong password"), http.StatusBadRequest)
		return
	}
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token.Token, err = jwtToken.SignedString([]byte(passTrue))
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(res, err, http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(token)
}

func (s *TodoTaskServer) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pass := s.password
		if len(pass) > 0 {
			cookie, err := req.Cookie("token")
			if err != nil {
				log.Println(err.Error())
				ErrorHandler(res, err, http.StatusUnauthorized)
				return
			}
			jwtCookie := cookie.Value
			jwtToken, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil {
				log.Println(err.Error())
				ErrorHandler(res, err, http.StatusUnauthorized)
				return
			}
			if !jwtToken.Valid {
				log.Println("token is invalid")
				ErrorHandler(res, fmt.Errorf("token is invalid"), http.StatusUnauthorized)
				return
			}
		}
		next(res, req)
	})
}
