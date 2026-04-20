package validation

import (
	"errors"
	"slices"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
)

func ValidateTaskUpdate(taskUpdate dto.TaskPATCH, validTypeIDs []string) error {
	if taskUpdate.Label != nil {
		label := *taskUpdate.Label
		if len(label) > 50 {
			return errors.New("label is too long")
		}

		if label == "" {
			return errors.New("label is empty")
		}
	}

	if taskUpdate.Description != nil {
		description := *taskUpdate.Description
		if len(description) > 300 {
			return errors.New("description is too long")
		}
	}

	if taskUpdate.TypeID != nil && len(validTypeIDs) > 0 {
		if slices.Contains(validTypeIDs, *taskUpdate.TypeID) {
			return nil
		}
		return errors.New("type_id incorrect")
	}

	return nil
}
