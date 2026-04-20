package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sater-151/todo-list/internal/repository/postgres/board"
	"github.com/sater-151/todo-list/internal/repository/postgres/column"
	moveevent "github.com/sater-151/todo-list/internal/repository/postgres/move-event"
	"github.com/sater-151/todo-list/internal/repository/postgres/task"
	db_type "github.com/sater-151/todo-list/internal/repository/postgres/type"
	"github.com/sater-151/todo-list/internal/repository/postgres/user"
	"github.com/sater-151/todo-list/pkg/utils"
)

type Repository struct {
	Board     *board.BoardStorage
	Column    *column.ColumnStorage
	Type      *db_type.TypeStorage
	Task      *task.TaskStorage
	User      *user.UserStorage
	MoveEvent *moveevent.MoveEventStorage
	Database  *pgxpool.Pool
}

func NewRepository(storage *pgxpool.Pool, schema string, cryptoKey string) (repo *Repository, err error) {
	defer utils.AddFuncLabel("[init-repository]", err)

	boardStorage, err := board.New(storage, schema)
	if err != nil {
		return
	}

	columnStorage, err := column.New(storage, schema)
	if err != nil {
		return
	}

	typeStorage, err := db_type.New(storage, schema)
	if err != nil {
		return
	}

	taskStorage, err := task.New(storage, schema)
	if err != nil {
		return
	}

	userStorage, err := user.New(storage, schema, cryptoKey)
	if err != nil {
		return
	}

	moveEventStorage, err := moveevent.New(storage, schema)
	if err != nil {
		return
	}

	return &Repository{
		Board:     boardStorage,
		Column:    columnStorage,
		Type:      typeStorage,
		Task:      taskStorage,
		User:      userStorage,
		MoveEvent: moveEventStorage,
		Database:  storage,
	}, nil

}
