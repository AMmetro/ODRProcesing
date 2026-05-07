package statistics_postgres_repository

import "github.com/AMmetro/ODRProcesing/internal/core/domain"

type StatisticsModel struct {
	TotalUsers     int
	TotalTasks     int
	CompletedTasks int
	PendingTasks   int
}

func (m StatisticsModel) ToDomain() domain.StatisticsSummary {
	return domain.NewStatisticsSummary(
		m.TotalUsers,
		m.TotalTasks,
		m.CompletedTasks,
		m.PendingTasks,
	)
}
