package validation

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func ValidateColumnCreate(allBoardColumns entity.Columns, newColumn dto.ColumnPOST) error {
	if newColumn.Name == "" {
		return fmt.Errorf("name is empty")
	}

	nameGroup := utils.NameGroup(allBoardColumns)
	if _, ok := nameGroup[newColumn.Name]; ok {
		return fmt.Errorf("column with name %s already exists", newColumn.Name)
	}

	if newColumn.OrderNumber <= 0 && newColumn.OrderNumber != -1 {
		return fmt.Errorf("order number can't be below or eq 0")
	}

	orderNumberGroup := utils.OrderNumberGroup(allBoardColumns)
	if _, ok := orderNumberGroup[newColumn.OrderNumber]; ok {
		return fmt.Errorf("column with order number %d already exists", newColumn.OrderNumber)
	}

	return nil
}
