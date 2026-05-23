package web_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

// type GetTaskResponse TaskDTOResponse

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
func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Get MainPage handler")

	res, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to index.html for main page")
		return
	}

	response := res
	responseHandler.JSONResponse(response, http.StatusOK)

}
