package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    INSERT INTO ODRProcesing.tasks (title, description, completed, author_user_id, created_at, completed_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, version, title, description, completed, author_user_id, created_at, completed_at;
    `

	row := r.pool.QueryRow(
		ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.AuthorUserId,
		task.CreatedAt,
		task.CompletedAt,
	)

	var taskModel TasksModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.AuthorUserId,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user with id=%d: %w",
				err,
				task.AuthorUserId,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.AuthorUserId,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
	)

	return taskDomain, nil
}
