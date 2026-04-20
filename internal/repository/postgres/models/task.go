package models

import (
	"database/sql"

	"github.com/sater-151/todo-list/internal/entity"
)

type Task struct {
	ID          string       `db:"id"`
	TypeID      string       `db:"type_id"`
	ColumnID    string       `db:"column_id"`
	Label       string       `db:"label"`
	Description string       `db:"description"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

func (t Task) ToEntity() entity.Task {
	return entity.Task{
		ID:          t.ID,
		TypeID:      t.TypeID,
		ColumnID:    t.ColumnID,
		Label:       t.Label,
		Description: t.Description,
		UpdatedAt:   t.UpdatedAt.Time,
		CreatedAt:   t.CreatedAt.Time,
	}
}
