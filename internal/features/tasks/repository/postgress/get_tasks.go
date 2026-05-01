package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userId *int,
) ([]domain.Task, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT 
	id,
	version,
	title,
	description,
	completed,
	author_user_id,
	created_at,
	completed_at
	FROM ODRProcesing.tasks
    `

	/*
		/dynamic SQL query building
	*/
	defaultLimit := 100
	defaultOffset := 0

	if limit == nil {
		limit = &defaultLimit
	}

	if offset == nil {
		offset = &defaultOffset
	}
	args := []any{}
	argN := 1

	if userId != nil {
		query += fmt.Sprintf(" WHERE author_user_id = $%d", argN)
		args = append(args, *userId)
		argN++
	}

	query += " ORDER BY created_at ASC"

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, *limit, *offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var tasksModel []TaskModel

	for rows.Next() {
		var task TaskModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan error: %w", err)
		}

		tasksModel = append(tasksModel, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	tasksDomain := TasksDomainsFromModels(tasksModel)

	return tasksDomain, nil
}
