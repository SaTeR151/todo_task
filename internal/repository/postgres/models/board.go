package models

import "github.com/sater-151/todo-list/internal/entity"

type Board struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Name   string `db:"name"`
}

func (b *Board) ToEntity() entity.Board {
	entityBoard := entity.Board{}

	entityBoard.ID = b.ID
	entityBoard.UserID = b.UserID
	entityBoard.Name = b.Name

	return entityBoard
}
