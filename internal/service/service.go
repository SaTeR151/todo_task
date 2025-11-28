package service

import (
	"strings"
	"time"

	"github.com/sater-151/todo-list/internal/database"
	"github.com/sater-151/todo-list/internal/models"
	"github.com/sater-151/todo-list/internal/utils"
)

type Service struct {
	db *database.DBStruct
}

func New(db *database.DBStruct) *Service {
	return &Service{db: db}
}

func (s *Service) AddTask(task models.Task) (models.ID, error) {
	var Id models.ID
	var err error
	task, err = utils.CheckTask(task)
	if err != nil {
		return Id, err
	}
	Id.ID, err = s.db.InsertTask(task)
	if err != nil {
		return Id, err
	}
	return Id, nil
}

func (s *Service) UpdateTask(task models.Task) error {
	task, err := utils.CheckTask(task)
	if err != nil {
		return err
	}
	err = s.db.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) TaskDone(selectConfig models.SelectConfig) error {
	tasks, err := s.db.Select(selectConfig)
	if err != nil {
		return err
	}
	task := tasks[0]
	if task.Repeat == "" {
		err = s.db.DeleteTask(selectConfig.Id)
		if err != nil {
			return err
		}
		return nil
	}
	task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}
	err = s.db.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetListTask(selectConfig models.SelectConfig) ([]models.Task, error) {
	if selectConfig.Search != "" {
		date := strings.Split(selectConfig.Search, ".")
		if len(date) == 3 {
			var d string
			for i := 2; i >= 0; i-- {
				d += date[i]
			}
			_, err := time.Parse("20060102", d)
			if err == nil {
				selectConfig.Search = ""
				selectConfig.Date = d
			}
		}
	}
	listTask, err := s.db.Select(selectConfig)
	if err != nil {
		return listTask, err
	}
	return listTask, err
}
