package models

import (
	"database/sql"

	"github.com/sater-151/todo-list/internal/entity"
)

type User struct {
	ID           string         `db:"id"`
	Login        string         `db:"login"`
	Password     string         `db:"password"`
	RefreshToken sql.NullString `db:"refresh_token"`
}

func (u *User) ToEntity() entity.User {
	return entity.User{
		ID:           u.ID,
		Login:        u.Login,
		Password:     u.Password,
		RefreshToken: u.RefreshToken.String,
	}
}
