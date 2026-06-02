package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {

	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be after before: %w", core_errors.ErrInvalidArgument)
		}
	}

	summary, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("the 'to' date must be after the 'from' date %w", err)
	}

	statistics := calcStatistics(summary)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	var totalCompletedDuration time.Duration
	tasksCompleted := 0
	for _, task := range tasks {

		completionDuration := task.CompletionDuration()

		if completionDuration != nil && *completionDuration > 0 {
			// if *task.CompleteionDuration() > 0 {
			tasksCompleted++
			completionDuration := task.CompletionDuration()
			if completionDuration != nil {
				totalCompletedDuration += *completionDuration
			}
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}
	return domain.NewStatistics(tasksCreated, tasksCompleted, &tasksCompletedRate, tasksAverageCompletionTime)
}
