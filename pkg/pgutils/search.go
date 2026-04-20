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
		squirrel.Expr(
			"? = ANY("+field+")", value,
		),
	)
}

func SearchMultiEq(query squirrel.SelectBuilder, field string, search []string) squirrel.SelectBuilder {
	if len(search) == 0 || field == "" {
		return query
	}

	or := squirrel.Or{}

	for _, v := range search {
		value := fmt.Sprint(v)

		if value == "" {
			continue
		}

		or = append(or, squirrel.Expr(
			"? = ANY("+field+")", value,
		))

	}

	return query.Where(or)

}
