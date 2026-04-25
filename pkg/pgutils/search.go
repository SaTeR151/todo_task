package pgutils

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func SearchEq(query squirrel.SelectBuilder, field string, search any) squirrel.SelectBuilder {
	value := fmt.Sprint(search)

	if search == nil || value == "" || field == "" {
		return query
	}

	return query.Where(
		squirrel.Eq{
			field: value,
		},
	)
}

func SearchMultiEq(query squirrel.SelectBuilder, field string, search []string) squirrel.SelectBuilder {
	if len(search) == 0 || field == "" {
		return query
	}

	or := squirrel.Or{}

	for _, value := range search {

		if value == "" {
			continue
		}

		or = append(or, squirrel.Eq{
			field: value,
		})

	}

	return query.Where(or)

}
