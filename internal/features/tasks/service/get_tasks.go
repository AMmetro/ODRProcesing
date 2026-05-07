package tasks_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	domain_utils "github.com/AMmetro/ODRProcesing/internal/core/domain/utils"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userId *int,
) ([]domain.Task, error) {

	err := domain_utils.LimitOffsetValidation(limit, offset)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("validate limit and offset: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset, userId)

	if err != nil {
		return []domain.Task{}, fmt.Errorf("get tasks: %w", err)
	}

	return tasks, nil
}
