package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sater-151/todo-list/internal/models"
)

type TodoTaskRepo struct {
	pool *pgxpool.Pool
}

func NewTodoTaskRepo(pool *pgxpool.Pool) *TodoTaskRepo {
	return &TodoTaskRepo{
		pool: pool,
	}
}

func (r *TodoTaskRepo) InsertTask(ctx context.Context, task models.Task) (string, error) {
	var id int

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO scheduler (date, title, comment, repeat)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		task.Date, task.Title, task.Comment, task.Repeat,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

func (r *TodoTaskRepo) UpdateTask(ctx context.Context, task models.Task) error {
	res, err := r.pool.Exec(ctx, "UPDATE scheduler SET date = &1, title = &2, comment = &3, repeat = &4 WHERE id = &5",
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("data update error")
	}

	return nil
}

func (r *TodoTaskRepo) DeleteTask(ctx context.Context, uuid string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM scheduler WHERE id = &1", uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoTaskRepo) Select(ctx context.Context, selectConfig models.SelectConfig) ([]models.Task, error) {
	listTask := []models.Task{}
	row := fmt.Sprintf("SELECT * FROM %s", selectConfig.Table)
	if selectConfig.Search != "" || selectConfig.Date != "" || selectConfig.Id != "" {
		row += " WHERE"
	}
	if selectConfig.Search != "" {
		row += fmt.Sprintf(" title LIKE %s OR comment LIKE %s", "'%"+selectConfig.Search+"%'", "'%"+selectConfig.Search+"%'")
	}
	if selectConfig.Date != "" {
		row += fmt.Sprintf(" date = '%s'", selectConfig.Date)
	}
	if selectConfig.Id != "" {
		row += fmt.Sprintf(" id = %s", selectConfig.Id)
	}
	if selectConfig.Sort != "" {
		row += fmt.Sprintf(" ORDER BY %s %s", selectConfig.Sort, selectConfig.TypeSort)
	}
	if selectConfig.Limit != "" {
		row += fmt.Sprintf(" LIMIT %s", selectConfig.Limit)
	}
	res, err := r.pool.Query(ctx, row)
	if err != nil {
		return listTask, err
	}
	for res.Next() {
		task := models.Task{}
		err = res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return listTask, err
		}
		listTask = append(listTask, task)
	}
	defer res.Close()
	return listTask, nil
}
