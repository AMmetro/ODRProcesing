package statistics_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool"
)

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}

func (r *StatisticsRepository) GetStatistics(
	ctx context.Context,
) (domain.StatisticsSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT
	(SELECT COUNT(*) FROM odrprocesing.users) AS total_users,
	(SELECT COUNT(*) FROM odrprocesing.tasks) AS total_tasks,
	(SELECT COUNT(*) FROM odrprocesing.tasks WHERE completed = TRUE) AS completed_tasks,
	(SELECT COUNT(*) FROM odrprocesing.tasks WHERE completed = FALSE) AS pending_tasks;`

	row := r.pool.QueryRow(ctx, query)
	var statisticsModel StatisticsModel
	if err := row.Scan(
		&statisticsModel.TotalUsers,
		&statisticsModel.TotalTasks,
		&statisticsModel.CompletedTasks,
		&statisticsModel.PendingTasks,
	); err != nil {
		return domain.StatisticsSummary{}, fmt.Errorf("scan statistics: %w", err)
	}

	return statisticsModel.ToDomain(), nil
}
