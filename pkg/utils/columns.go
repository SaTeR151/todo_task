package utils

import "github.com/sater-151/todo-list/internal/entity"

func NameGroup(columns entity.Columns) map[string]struct{} {
	columnsMap := map[string]struct{}{}
	for _, column := range columns {
		columnsMap[column.Name] = struct{}{}
	}
	return columnsMap
}

func OrderNumberGroup(columns entity.Columns) map[int]struct{} {
	columnsMap := map[int]struct{}{}
	for _, column := range columns {
		if column.OrderNumber == -1 {
			continue
		}
		columnsMap[column.OrderNumber] = struct{}{}
	}
	return columnsMap
}

func ColumnsExcept(columns entity.Columns, columnID string) entity.Columns {
	var result entity.Columns
	for _, column := range columns {
		if column.ID != columnID {
			result = append(result, column)
		}
	}
	return result
}
