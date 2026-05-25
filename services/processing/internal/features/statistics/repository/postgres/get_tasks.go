package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
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
	`)

	/*
		/dynamic SQL query building
	*/

	args := []any{}
	conditions := []string{}

	if userId != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userId)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		// query += fmt.Sprintf(" WHERE " + strings.Join(conditions, " AND "))
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	// query += " ORDER BY id ASC"
	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var tasksModal []TaskModel

	for rows.Next() {
		var taskModal TaskModel
		err := rows.Scan(
			&taskModal.ID,
			&taskModal.Version,
			&taskModal.Title,
			&taskModal.Description,
			&taskModal.Completed,
			&taskModal.AuthorUserId,
			&taskModal.CreatedAt,
			&taskModal.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		tasksModal = append(tasksModal, taskModal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows tasks: %w", err)
	}

	tasksDomains := TasksDomainsFromModels(tasksModal)

	return tasksDomains, nil
}

