package statistics_service

import (
	"context"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatistics(
		ctx context.Context,
	) (domain.StatisticsSummary, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
) (domain.StatisticsSummary, error) {
	summary, err := s.statisticsRepository.GetStatistics(ctx)
	if err != nil {
		return domain.StatisticsSummary{}, err
	}

	if summary.TotalUsers > 0 {
		// summary.TasksPerUser = float64(summary.TotalTasks) / float64(summary.TotalUsers)
	}

	return summary, nil
}
