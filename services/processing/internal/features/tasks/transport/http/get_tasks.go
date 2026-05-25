package tasks_transport_http

import (
	"net/http"

	domain_utils "github.com/AMmetro/ODRProcesing/shared/pkg/core/domain/utils"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

// @Summary Get Tasks
// @Description Retrieve a list of tasks with optional filtering by user ID, limit, and offset
// @Tags tasks
// @Accept json
// @Produce json
// @Param limit query int false "Number of tasks to retrieve" minimum(1) maximum(100)
// @Param offset query int false "Number of tasks to skip" minimum(0)
// @Param userId query int false "Filter tasks by user ID"
// @Success 200 {array} TaskDTOResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Get Tasks handler")

	limit, offset, userId, err := domain_utils.GetUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	tasksDomain, err := h.tasksService.GetTasks(ctx, limit, offset, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}
	response := TasksDtoFromDomains(tasksDomain)
	responseHandler.JSONResponse(response, http.StatusOK)

}

