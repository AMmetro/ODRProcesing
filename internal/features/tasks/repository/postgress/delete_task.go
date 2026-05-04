package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	taskId int,
) error {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    DELETE FROM ODRProcesing.tasks WHERE id = $1;
    `

	cmdTag, err := r.pool.Exec(ctx, query, taskId)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id='%d': %w", taskId, core_errors.ErrNotFound)
	}

	return nil
}
