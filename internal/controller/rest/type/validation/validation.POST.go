package validation

import (
	"errors"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func ValidateTypeCreate(typeCreate entity.TypeCreate, types entity.Types) error {

	if typeCreate.Name == "" {
		return errors.New("name is empty")
	}

	for _, t := range types {
		if t.Name == typeCreate.Name {
			return errors.New("type already exists")
		}
	}

	if typeCreate.Color == "" {
		return errors.New("color is empty")
	}

	if len(typeCreate.Name) > 10 {
		return errors.New("name is too long (max 10 symbols)")
	}

	if !utils.IsMatchRegexp(typeCreate.Name, `^[a-z0-9_.-]+$`) {
		return errors.New("name must contain only a-z 0-9 '_' '-' '.'")
	}

	if !utils.IsMatchRegexp(typeCreate.Color, `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`) {
		return errors.New("color must be in HEX format")
	}

	return nil
}
