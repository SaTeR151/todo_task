package selectconfig

import "github.com/sater-151/todo-list/internal/models"

func Default() *models.SelectConfig {
	return &models.SelectConfig{
		Limit: "20",
		Sort:  "date",
		Table: "scheduler",
	}
}
