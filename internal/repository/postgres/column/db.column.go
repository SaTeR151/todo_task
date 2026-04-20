package column

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres/models"
	"github.com/sater-151/todo-list/pkg/pgutils"
	"github.com/sater-151/todo-list/pkg/utils"
	"github.com/sirupsen/logrus"
)

func (s *ColumnStorage) GetColumns(ctx context.Context, opts entity.GetColumnsOpts) (res entity.Columns, err error) {
	defer utils.AddFuncLabel("[repo-get-columns]", err)

	query := s.queryBuilder.Select(
		"id",
		"label",
		"board_id",
		"order_number",
	).From(fmt.Sprintf("%s.%s AS columns", s.scheme, TABLE_COLUMNS))

	query = pgutils.SearchEq(query, "columns.board_id", opts.BoardID)
	query = pgutils.SearchEq(query, "columns.order_number", opts.OrderNumber)

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [get-columns: INCOMING] \n opts: %v \n", opts)
	logrus.Debugf("\n [get-columns: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	rows, err := s.client.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		c := models.Column{}

		scanValues := []any{
			&c.ID,
			&c.Name,
			&c.BoardID,
			&c.OrderNumber,
		}

		if err = rows.Scan(scanValues...); err != nil {
			return nil, err
		}

		columnEntity := c.ToEntity()

		res = append(res, columnEntity)
	}

	return
}

func (s *ColumnStorage) CreateColumn(ctx context.Context, column entity.ColumnCreate) (res string, err error) {
	defer utils.AddFuncLabel("[repo-create-column]", err)

	sql, args, err := s.queryBuilder.
		Insert(fmt.Sprintf("%s.%s", s.scheme, TABLE_COLUMNS)).
		Columns(
			"label",
			"board_id",
			"order_number",
		).Values(
		column.Name,
		column.BoardID,
		column.OderNumber,
	).
		Suffix("RETURNING id").
		ToSql()

	logrus.Debugf("\n [create-column: INCOMING] \n column-create: %+v \n", column)
	logrus.Debugf("\n [create-column: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return res, err
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *ColumnStorage) UpdateColumn(ctx context.Context, column entity.ColumnUpdate) (err error) {
	defer utils.AddFuncLabel("[repo-update-column]", err)

	query := s.queryBuilder.
		Update(fmt.Sprintf("%s.%s", s.scheme, TABLE_COLUMNS)).
		Where(squirrel.Eq{"id": column.ID})

	if column.Name != nil {
		query = query.Set("name", *column.Name)
	}

	if column.OrderNumber != nil {
		query = query.Set("order_number", *column.OrderNumber)
	}

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [update-column: INCOMING] \n column-update: %+v \n", column)
	logrus.Debugf("\n [update-column: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}

func (s *ColumnStorage) DeleteColumn(ctx context.Context, columnID string) (err error) {
	defer utils.AddFuncLabel("[repo-delete-column]", err)

	sql, args, err := s.queryBuilder.
		Delete(fmt.Sprintf("%s.%s", s.scheme, TABLE_COLUMNS)).
		Where(squirrel.Eq{"id": columnID}).
		ToSql()

	logrus.Debugf("\n [delete-column: INCOMING] \n column-id: %s \n", columnID)
	logrus.Debugf("\n [delete-column: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}
