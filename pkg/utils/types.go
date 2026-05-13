package utils

import "github.com/sater-151/todo-list/internal/entity"

func TypesExcept(types entity.Types, typeID string) entity.Types {
	var result entity.Types
	for _, type_ := range types {
		if type_.ID != typeID {
			result = append(result, type_)
		}
	}
	return result
}
