package models

import (
	"database/sql"

	"github.com/sater-151/todo-list/internal/entity"
)

type MoveEvent struct {
	ID           string       `db:"id"`
	TaskID       string       `db:"task_id"`
	FromColumnID string       `db:"from_column_id"`
	ToColumnID   string       `db:"to_column_id"`
	Timestamp    sql.NullTime `db:"updated_at"`
}

func (me MoveEvent) ToEntity() entity.MoveEvent {
	return entity.MoveEvent{
		ID:           me.ID,
		TaskID:       me.TaskID,
		FromColumnID: me.FromColumnID,
		ToColumnID:   me.ToColumnID,
		Timestamp:    me.Timestamp.Time,
	}
}
