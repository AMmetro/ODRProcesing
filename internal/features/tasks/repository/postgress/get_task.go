package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	taskId int,
) (domain.Task, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT id, version, title, description, completed, author_user_id, created_at, completed_at
	FROM ODRProcesing.tasks WHERE id = $1;
    `

	row := r.pool.QueryRow(ctx, query, taskId)

	var task TaskModel
	err := row.Scan(
		&task.ID,
		&task.Version,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.AuthorUserId,
		&task.CreatedAt,
		&task.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{},
				fmt.Errorf("task with id='%d': %w", taskId, core_errors.ErrNotFound)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := domain.Task{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		AuthorUserId: task.AuthorUserId,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
	}

	return taskDomain, nil
}
