package tasks_service

import (
	"context"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		limit *int,
		ofset *int,
		userId *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		taskId int,
	) (domain.Task, error)

	UpdateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		taskId int,
	) error
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
