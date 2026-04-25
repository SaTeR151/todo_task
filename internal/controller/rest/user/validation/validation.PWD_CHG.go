package validation

import (
	"errors"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
)

func ValidateUserPasswordChange(pwdChange dto.UserPasswordChange, userCurrentPassword string) error {

	if userCurrentPassword != pwdChange.OldPassword {
		return errors.New("old password is incorrect")
	}

	if pwdChange.OldPassword == "" {
		return errors.New("old password is empty")
	}

	if pwdChange.NewPassword == "" {
		return errors.New("new password is empty")
	}

	if pwdChange.NewPassword == pwdChange.OldPassword {
		return errors.New("new password is equal to old password")
	}

	if len(pwdChange.NewPassword) > 30 {
		return errors.New("new password length is more then 30 symbols")
	}

	return nil
}
