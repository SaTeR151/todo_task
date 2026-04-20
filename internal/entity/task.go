package entity

import "time"

type Task struct {
	ID          string    `json:"id"`
	TypeID      string    `json:"type_id"`
	ColumnID    string    `json:"column_id"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tasks []Task

type TaskCreate struct {
	TypeID      string `json:"type_id"`
	ColumnID    string `json:"column_id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type TaskUpdate struct {
	ID          string  `json:"id"`
	TypeID      *string `json:"type_id"`
	ColumnID    *string `json:"column_id"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

type GetTasksOpts struct {
	ID        string
	ColumnID  string
	ColumnIDs []string
	TypeID    string
}
