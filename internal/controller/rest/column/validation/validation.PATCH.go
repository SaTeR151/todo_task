package validation

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func ValidateColumnUpdate(allBoardColumns entity.Columns, newColumn dto.ColumnPATCH) error {
	if newColumn.Name != nil {
		name := *newColumn.Name
		if name == "" {
			return fmt.Errorf("name is empty")
		}

		nameGroup := utils.NameGroup(allBoardColumns)
		if _, ok := nameGroup[name]; ok {
			return fmt.Errorf("column with name %s already exists", name)
		}
	}

	if newColumn.OrderNumber != nil {
		orderNumber := *newColumn.OrderNumber
		if orderNumber <= 0 && orderNumber != -1 {
			return fmt.Errorf("order number can't be below or eq 0")
		}

		orderNumberGroup := utils.OrderNumberGroup(allBoardColumns)
		if _, ok := orderNumberGroup[orderNumber]; ok {
			return fmt.Errorf("column with order number %d already exists", newColumn.OrderNumber)
		}
	}

	return nil
}
