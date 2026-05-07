package statistics_transport_http

import (
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TotalUsers     int `json:"total_users"`
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	PendingTasks   int `json:"pending_tasks"`
	// TasksPerUser   float64 `json:"tasks_per_user"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetStatistics handler")

	summary, err := h.statisticsService.GetStatistics(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := GetStatisticsResponseFromDomain(summary)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func GetStatisticsResponseFromDomain(summary domain.StatisticsSummary) GetStatisticsResponse {
	return GetStatisticsResponse{
		TotalUsers:     summary.TotalUsers,
		TotalTasks:     summary.TotalTasks,
		CompletedTasks: summary.CompletedTasks,
		PendingTasks:   summary.PendingTasks,
		// TasksPerUser:   summary.TasksPerUser,
	}
}
