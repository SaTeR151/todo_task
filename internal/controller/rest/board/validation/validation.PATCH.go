package validation

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func ValidateBoardUpdate(boardCreate dto.BoardPATCH, boards entity.Boards) error {
	if boardCreate.Name != nil {
		name := *boardCreate.Name

		if name == "" {
			return fmt.Errorf("name is empty")
		}

		if len(name) > 50 {
			return fmt.Errorf("name is too long")
		}

		for _, board := range boards {
			if board.Name == name {
				return fmt.Errorf("board with name %s already exists", name)
			}
		}
	}

	return nil
}
