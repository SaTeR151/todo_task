package validation

import (
	"errors"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func ValidateTypeUpdate(typeUpdate entity.TypeUpdate, types entity.Types) error {
	if typeUpdate.Name != nil {
		name := *typeUpdate.Name

		if name == "" {
			return errors.New("name is empty")
		}

		if len(name) > 10 {
			return errors.New("name is too long (max 10 symbols)")
		}

		if !utils.IsMatchRegexp(name, `^[a-z0-9_.-]+$`) {
			return errors.New("name must contain only a-z 0-9 '_' '-' '.'")
		}

		for _, t := range types {
			if t.Name == name {
				return errors.New("type already exists")
			}
		}
	}

	if typeUpdate.Color != nil {
		color := *typeUpdate.Color

		if color == "" {
			return errors.New("color is empty")
		}

		if !utils.IsMatchRegexp(color, `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`) {
			return errors.New("color must be in HEX format")
		}
	}

	return nil
}
