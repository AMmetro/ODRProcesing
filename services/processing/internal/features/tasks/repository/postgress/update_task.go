package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/shared/pkg/core/repository/postgres/pool"
)

func (r *TasksRepository) UpdateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    UPDATE ODRProcesing.tasks 
	SET title=$1, description=$2, completed=$3, completed_at=$4, version = version + 1
    WHERE id = $5 AND version = $6
    RETURNING id, version, title, description,
	completed, author_user_id, created_at, completed_at;
    `

	row := r.pool.QueryRow(
		ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version,
	)

	var taskModel TaskModel

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"%v: task with id=`%d` concurently accessed: %w",
				task.AuthorUserId,
				err,
				core_errors.ErrConflict,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)

	return taskDomain, nil
}

