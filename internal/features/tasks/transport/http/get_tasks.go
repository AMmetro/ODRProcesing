package tasks_transport_http

import (
	"net/http"

	domain_utils "github.com/AMmetro/ODRProcesing/internal/core/domain/utils"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

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
