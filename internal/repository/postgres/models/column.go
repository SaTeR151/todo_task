package models

import "github.com/sater-151/todo-list/internal/entity"

type Column struct {
	ID          string `db:"id"`
	BoardID     string `db:"board_id"`
	Name        string `db:"name"`
	OrderNumber int    `db:"order_number"`
}

func (c *Column) ToEntity() entity.Column {
	return entity.Column{
		ID:          c.ID,
		BoardID:     c.BoardID,
		Name:        c.Name,
		OrderNumber: c.OrderNumber,
	}
}
