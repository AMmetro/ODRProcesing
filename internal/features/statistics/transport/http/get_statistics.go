package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated              int      `json:"task_created"`
	TasksCompleted            int      `json:"task_completed"`
	TaskCompletedRate         *float64 `json:"task_completed_rate"`
	TaskAverageComplitionTime *string  `json:"task_average_complition_time"` // "1m30sec"
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetStatistics handler")

	userId, err := core_http_request.GetIntQueryParam(r, "userId")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId params")
		return
	}
	from, err := core_http_request.GetDataQueryParam(r, "from")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get fromDate from params")
		return
	}
	to, err := core_http_request.GetDataQueryParam(r, "to")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get toDate from params")
		return
	}

	fmt.Println(userId, from, to)

	summary, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	// summary, err := h.statisticsService.GetStatistics()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := GetStatisticsResponseFromDomain(summary)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func GetStatisticsResponseFromDomain(summary domain.Statistics) GetStatisticsResponse {

	var avgTime *string
	if summary.TaskAverageComplitionTime != nil {
		duration := summary.TaskAverageComplitionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:              summary.TaskCreated,
		TasksCompleted:            summary.TasksCompleted,
		TaskCompletedRate:         summary.TaskCompletedRate,
		TaskAverageComplitionTime: avgTime,
	}
}
