package validation

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/entity"
)

func ValidateColumnSwap(columnA entity.Column, columnB entity.Column) error {
	if columnA.ID == columnB.ID {
		return fmt.Errorf("can't swap column with itself")
	}

	return nil
}
