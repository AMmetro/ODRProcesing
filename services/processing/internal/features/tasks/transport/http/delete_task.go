package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

// @Summary Delete Task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 404 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /task/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Dekete Tasks handler")

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId")
		return
	}

	err = h.tasksService.DeleteTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete users")
		return
	}
	responseHandler.JSONResponse(nil, http.StatusNoContent)

}

