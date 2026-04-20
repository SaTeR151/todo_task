package models

import "github.com/sater-151/todo-list/internal/entity"

type Type struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Color  string `db:"color"`
}

type Types []Type

func (t *Type) ToEntity() entity.Type {
	entityType := entity.Type{}

	entityType.ID = t.ID
	entityType.Color = t.Color
	entityType.Name = t.Name
	entityType.UserID = t.UserID

	return entityType
}
