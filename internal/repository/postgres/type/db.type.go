package db_type

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

func (s *TypeStorage) Get(ctx context.Context, opts entity.GetTypesOpts) (res entity.Types, err error) {
	defer utils.AddFuncLabel("[repo-get-type]", err)

	query := s.queryBuilder.Select(
		"id",
		"user_id",
		"name",
		"color",
	).From(fmt.Sprintf("%s.%s AS types", s.scheme, TABLE_TYPES))

	query = pgutils.SearchEq(query, "types.id", opts.ID)
	query = pgutils.SearchEq(query, "types.user_id", opts.UserID)

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [get-type: INCOMING] \n OPTS: %+v \n", opts)
	logrus.Debugf("\n [get-type: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return nil, err
	}

	rows, err := s.client.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		t := models.Type{}

		scanValues := []any{
			&t.ID,
			&t.UserID,
			&t.Name,
			&t.Color,
		}

		if err := rows.Scan(scanValues...); err != nil {
			return nil, err
		}

		typeEntity := t.ToEntity()

		res = append(res, typeEntity)
	}

	return
}

func (s *TypeStorage) Create(ctx context.Context, typeCreate entity.TypeCreate) (res string, err error) {
	defer utils.AddFuncLabel("[repo-type-create]", err)

	sql, args, err := s.queryBuilder.
		Insert(fmt.Sprintf("%s.%s", s.scheme, TABLE_TYPES)).
		Columns(
			"user_id",
			"name",
			"color",
		).Values(
		typeCreate.UserID,
		typeCreate.Name,
		typeCreate.Color,
	).
		Suffix("RETURNING id").
		ToSql()

	logrus.Debugf("\n [create-type: INCOMING] \n type-create: %+v \n", typeCreate)
	logrus.Debugf("\n [create-type: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return res, err
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *TypeStorage) Update(ctx context.Context, typeUpdate entity.TypeUpdate) (err error) {
	defer utils.AddFuncLabel("[repo-type-update]", err)

	query := s.queryBuilder.
		Update(fmt.Sprintf("%s.%s", s.scheme, TABLE_TYPES)).
		Where(squirrel.Eq{"id": typeUpdate.ID})

	if typeUpdate.Color != nil {
		query = query.Set("color", *typeUpdate.Color)
	}

	if typeUpdate.Name != nil {
		query = query.Set("name", *typeUpdate.Name)
	}

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [update-type: INCOMING] \n type-update: %+v \n", typeUpdate)
	logrus.Debugf("\n [update-type: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}

func (s *TypeStorage) Delete(ctx context.Context, typeID string) (err error) {
	defer utils.AddFuncLabel("[repo-type-delete]", err)

	sql, args, err := s.queryBuilder.
		Delete(fmt.Sprintf("%s.%s", s.scheme, TABLE_TYPES)).
		Where(squirrel.Eq{"id": typeID}).
		ToSql()

	logrus.Debugf("\n [delete-type: INCOMING] \n type-id: %s \n", typeID)
	logrus.Debugf("\n [delete-type: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}
