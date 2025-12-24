package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sater-151/todo-list/internal/models"
)

type TodoTaskRepo struct {
	pool *pgxpool.Pool
}

func NewTodoTaskRepo(pool *pgxpool.Pool) (*TodoTaskRepo, error) {
	if pool == nil {
		return nil, fmt.Errorf("postgres.NewTodoTaskRepo: error = pool is nil")
	}

	return &TodoTaskRepo{
		pool: pool,
	}, nil
}

func (r *TodoTaskRepo) InsertTask(ctx context.Context, task *models.Task) (string, error) {
	taskUUID, err := uuid.NewV7()
	if err != nil {
		taskUUID = uuid.New()
	}

	_, err = r.pool.Exec(
		ctx,
		`INSERT INTO scheduler (
		uuid, 
		date, 
		title, 
		comment, 
		repeat
		)
		 VALUES ($1, $2, $3, $4, $5)`,
		taskUUID.String(), task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return "", err
	}

	return taskUUID.String(), nil
}

func (r *TodoTaskRepo) UpdateTask(ctx context.Context, task *models.Task) error {
	res, err := r.pool.Exec(ctx, "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE uuid = $5",
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

func (r *TodoTaskRepo) DeleteTask(ctx context.Context, taskUUID string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM scheduler WHERE uuid = $1", taskUUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoTaskRepo) Select(ctx context.Context, selectConfig *models.SelectConfig) ([]models.Task, error) {
	listTask := []models.Task{}
	row := fmt.Sprintf("SELECT * FROM %s", selectConfig.Table)
	if selectConfig.Search != "" || selectConfig.Date != "" || selectConfig.ID != "" {
		row += " WHERE"
	}
	if selectConfig.Search != "" {
		row += fmt.Sprintf(" title LIKE %s OR comment LIKE %s", "'%"+selectConfig.Search+"%'", "'%"+selectConfig.Search+"%'")
	}
	if selectConfig.Date != "" {
		row += fmt.Sprintf(" date = '%s'", selectConfig.Date)
	}
	if selectConfig.ID != "" {
		row += fmt.Sprintf(" uuid = '%s'", selectConfig.ID)
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
	defer res.Close()

	for res.Next() {
		task := models.Task{}
		err = res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return listTask, err
		}
		listTask = append(listTask, task)
	}

	return listTask, nil
}
