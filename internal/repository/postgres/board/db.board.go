package board

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

func (s *BoardStorage) Get(ctx context.Context, opts entity.GetBoardsOpts) (res entity.Boards, err error) {
	defer utils.AddFuncLabel("[repo-get-board]", err)

	query := s.queryBuilder.Select(
		"id",
		"user_id",
		"name",
	).From(fmt.Sprintf("%s.%s AS boards", s.scheme, TABLE_BOARDS))

	query = pgutils.SearchEq(query, "boards.id", opts.ID)
	query = pgutils.SearchEq(query, "boards.user_id", opts.UserID)

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [get-board: INCOMING] \n OPTS: %+v \n", opts)
	logrus.Debugf("\n [get-board: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	rows, err := s.client.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		b := models.Board{}

		scanValues := []any{
			&b.ID,
			&b.UserID,
			&b.Name,
		}

		if err = rows.Scan(scanValues...); err != nil {
			return
		}

		boardEntity := b.ToEntity()

		res = append(res, boardEntity)
	}

	return
}

func (s *BoardStorage) Create(ctx context.Context, boardCreate entity.BoardCreate) (res string, err error) {
	defer utils.AddFuncLabel("[repo-create-board]", err)

	sql, args, err := s.queryBuilder.
		Insert(fmt.Sprintf("%s.%s", s.scheme, TABLE_BOARDS)).
		Columns(
			"user_id",
			"name",
		).Values(
		boardCreate.UserID,
		boardCreate.Name,
	).
		Suffix("RETURNING id").
		ToSql()

	logrus.Debugf("\n [create-board: INCOMING] \n board=create: %+v \n", boardCreate)
	logrus.Debugf("\n [create-board: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *BoardStorage) Delete(ctx context.Context, boardID string) (err error) {
	defer utils.AddFuncLabel("[repo-delete-board]", err)

	sql, args, err := s.queryBuilder.
		Delete(fmt.Sprintf("%s.%s", s.scheme, TABLE_BOARDS)).
		Where(squirrel.Eq{"id": boardID}).
		ToSql()

	logrus.Debugf("\n [delete-board: INCOMING] \n board-id: %s \n", boardID)
	logrus.Debugf("\n [delete-board: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}

func (s *BoardStorage) Update(ctx context.Context, boardUdate entity.BoardUpdate) (err error) {
	defer utils.AddFuncLabel("[repo-update-board]", err)

	query := s.queryBuilder.
		Update(fmt.Sprintf("%s.%s", s.scheme, TABLE_BOARDS)).
		Where(squirrel.Eq{"id": boardUdate.ID})

	if boardUdate.Name != nil {
		query = query.Set("name", *boardUdate.Name)
	}

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [update-board: INCOMING] \n board-update: %+v \n", boardUdate)
	logrus.Debugf("\n [update-board: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}
