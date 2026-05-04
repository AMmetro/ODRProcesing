package tasks_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

func (s *TasksService) UpdateTask(
	ctx context.Context,
	taskId int,
	taskPatch domain.TaskPatch,
) (domain.Task, error) {

	existingTask, err := s.tasksRepository.GetTask(ctx, taskId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task for update: %w", err)
	}

	if err := existingTask.ApplyPatch(taskPatch); err != nil {
		return domain.Task{}, fmt.Errorf("applay patch task: %w", err)
	}

	updatedTask, err := s.tasksRepository.UpdateTask(ctx, existingTask)

	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}

	return updatedTask, nil
}
