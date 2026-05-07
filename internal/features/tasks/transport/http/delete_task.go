package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

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
