package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sater-151/todo-list/config"
	"github.com/sater-151/todo-list/database"
	"github.com/sater-151/todo-list/models"
	"github.com/sater-151/todo-list/service"
	"github.com/sater-151/todo-list/utils"
)

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

func GetNextDate(res http.ResponseWriter, req *http.Request) {
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

func PostTask(s *service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
		idJS, err = s.AddTask(task)
		if err != nil {
			log.Println(err.Error())
			ErrorHandler(res, err, http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(idJS)
	})
}

func ListTask(s *service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		var err error
		search := req.FormValue("search")
		selectConfig := CreateDefaultSelectConfig()
		if search != "" {
			selectConfig.Search = search
		}
		tasks, err := s.GetListTask(selectConfig)
		if err != nil {
			log.Println(err.Error())
			ErrorHandler(res, err, http.StatusBadRequest)
			return
		}
		listTask := models.ListTask{Tasks: tasks}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(listTask)
	})
}

func GetTask(db *database.DBStruct) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		id := req.FormValue("id")
		if id == "" {
			log.Println("id required")
			ErrorHandler(res, fmt.Errorf("id required"), http.StatusBadRequest)
			return
		}
		selectConfig := CreateDefaultSelectConfig()
		selectConfig.Id = id
		tasks, err := db.Select(selectConfig)
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
	})
}

func PutTask(s *service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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

		err = s.UpdateTask(task)
		if err != nil {
			log.Println(err.Error())
			ErrorHandler(res, err, http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(struct{}{})
	})
}

func PostTaskDone(s *service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		id := req.FormValue("id")
		selectConfig := CreateDefaultSelectConfig()
		selectConfig.Id = id
		err := s.TaskDone(selectConfig)
		if err != nil {
			log.Println(err.Error())
			ErrorHandler(res, err, http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(struct{}{})
	})
}

func DeleteTask(db *database.DBStruct) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		id := req.FormValue("id")
		if id == "" {
			log.Println("id required")
			ErrorHandler(res, fmt.Errorf("id required"), http.StatusBadRequest)
			return
		}
		err := db.DeleteTask(id)
		if err != nil {
			log.Println(err.Error())
			ErrorHandler(res, err, http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(struct{}{})
	})
}

func Sign(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	passTrue := config.GetPass()
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

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pass := config.GetPass()
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
