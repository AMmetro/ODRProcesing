package tasks_postgres_repository

import (
	"time"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	Status       domain.TaskStatus
	AuthorUserId int
	CreatedAt    time.Time
	CompletedAt  *time.Time
}

func TasksDomainsFromModels(tasks []TaskModel) []domain.Task {
	tasksDomain := make([]domain.Task, 0, len(tasks))
	for _, task := range tasks {
		tasksDomain = append(tasksDomain, domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.Status,
			task.AuthorUserId,
			task.CreatedAt,
			task.CompletedAt,
		))
	}
	return tasksDomain
}

func TaskDomainFromModel(task TaskModel) domain.Task {
	taskDomain := domain.NewTask(
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.Status,
		task.AuthorUserId,
		task.CreatedAt,
		task.CompletedAt,
	)
	return taskDomain
}
