package tasks_postgres_repository

import (
	"time"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

type TasksModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	AuthorUserId int
	CreatedAt    time.Time
	CompletedAt  *time.Time
}

func TasksDomainsFromModels(tasks []TasksModel) []domain.Task {

	tasksDomain := make([]domain.Task, 0, len(tasks))

	for _, task := range tasks {
		tasksDomain = append(tasksDomain, domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.AuthorUserId,
			task.CreatedAt,
			task.CompletedAt,
		))
	}

	return tasksDomain
}
