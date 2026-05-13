package validation

import (
	"errors"
	"fmt"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func ValidateUserCreate(userCreate entity.UserCreate, users entity.Users) error {
	login := userCreate.Login
	if login == "" {
		return errors.New("login is empty")
	}

	for _, u := range users {
		if u.Login == login {
			return fmt.Errorf("user %s already exists", login)
		}
	}

	if !utils.IsMatchRegexp(login, "^[a-z0-9]+$") {
		return errors.New("login is must have only [a-z] or [0-9]")
	}

	if len(login) > 10 {
		return errors.New("login length is more then 10 symbols")
	}

	if userCreate.Password == "" {
		return errors.New("password is empty")
	}

	if len(userCreate.Password) > 30 {
		return errors.New("password length is more then 30 symbols")
	}

	return nil

}
