package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated              int      `json:"task_created"`
	TasksCompleted            int      `json:"task_completed"`
	TaskCompletedRate         *float64 `json:"task_completed_rate"`
	TaskAverageComplitionTime *string  `json:"task_average_complition_time" example:"1m30s`
}

// @Summary Get Statistics
// @Description Retrieve statistics for tasks within a date range and optionally filtered by user ID
// @Tags statistics
// @Accept json
// @Produce json
// @Param userId query int false "User ID to filter statistics"
// @Param from query string false "Start date in YYYY-MM-DD format"
// @Param to query string false "End date in YYYY-MM-DD format"
// @Success 200 {object} GetStatisticsResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /statistics [get]
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

