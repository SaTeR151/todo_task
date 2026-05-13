package moveevent

import (
	"context"
	"fmt"
	"time"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres/models"
	"github.com/sater-151/todo-list/pkg/pgutils"
	"github.com/sater-151/todo-list/pkg/utils"
	"github.com/sirupsen/logrus"
)

func (s *MoveEventStorage) Get(ctx context.Context, opts entity.GetMoveEventsOpts) (res entity.MoveEvents, err error) {
	defer utils.AddFuncLabel("[repo-get-move-events]", err)

	query := s.queryBuilder.Select(
		"id",
		"task_id",
		"from_column_id",
		"to_column_id",
		"timestamp",
	).From(fmt.Sprintf("%s.%s AS move_events", s.scheme, TABLE_MOVE_EVENTS))

	query = pgutils.SearchEq(query, "move_events.task_id", opts.TaskID)

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [get-move-events: INCOMING] \n OPTS: %+v \n", opts)
	logrus.Debugf("\n [get-move-events: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	rows, err := s.client.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		me := models.MoveEvent{}

		scanValues := []any{
			&me.ID,
			&me.TaskID,
			&me.FromColumnID,
			&me.ToColumnID,
			&me.Timestamp,
		}

		if err = rows.Scan(scanValues...); err != nil {
			return
		}

		meEntity := me.ToEntity()

		res = append(res, meEntity)
	}

	return
}

func (s *MoveEventStorage) Create(ctx context.Context, moveEventCreate entity.MoveEventCreate) (err error) {
	defer utils.AddFuncLabel("[repo-create-move-event]", err)

	sql, args, err := s.queryBuilder.
		Insert(fmt.Sprintf("%s.%s", s.scheme, TABLE_MOVE_EVENTS)).
		Columns(
			"task_id",
			"from_column_id",
			"to_column_id",
			"timestamp",
		).Values(
		moveEventCreate.TaskID,
		moveEventCreate.FromColumnID,
		moveEventCreate.ToColumnID,
		time.Now(),
	).
		ToSql()

	logrus.Debugf("\n [create-move-event: INCOMING] \n move-event: %+v \n", moveEventCreate)
	logrus.Debugf("\n [create-move-event: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}
