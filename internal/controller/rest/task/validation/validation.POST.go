package validation

import (
	"errors"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func ValidateTaskCreate(taskCreate dto.TaskPOST, boardColumns entity.Columns, userTypes entity.Types) error {
	if taskCreate.Label == "" {
		return errors.New("label is empty")
	}

	if len(taskCreate.Label) > 50 {
		return errors.New("label is too long")
	}

	if taskCreate.ColumnID != "" {
		columnFound := false
		for _, column := range boardColumns {
			if column.ID == taskCreate.ColumnID {
				columnFound = true
				break
			}
		}
		if !columnFound {
			return errors.New("column_id incorrect")
		}
	}

	if len(taskCreate.Description) > 300 {
		return errors.New("description is too long")
	}

	if taskCreate.TypeID != "" {
		typeFound := false
		for _, taskType := range userTypes {
			if taskType.ID == taskCreate.TypeID {
				typeFound = true
				break
			}
		}
		if !typeFound {
			return errors.New("type_id incorrect")
		}
	}

	return nil
}
