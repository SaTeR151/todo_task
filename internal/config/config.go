package config

import (
	"os"

	"github.com/sater-151/todo-list/internal/models"
)

func GetConfig() models.Config {
	var config models.Config
	config.Port = os.Getenv("TODO_PORT")
	config.DbFilePath = os.Getenv("TODO_DBFILE")
	return config
}

func GetPass() (pass string) {
	pass = os.Getenv("TODO_PASSWORD")
	return pass
}
