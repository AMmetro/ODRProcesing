package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// @Summary Get Task
// @Description Retrieve a single task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} TaskDTOResponse
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /task/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Get Task handler")

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}
	response := GetTaskResponse(TaskDtoFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}
