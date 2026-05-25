package web_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

// type GetTaskResponse TaskDTOResponse

// @Summary Get WEB page
// @Description Retrieve a HTML about app
// @Web web
// @Accept json
// @Produce json
// @Success 200 {object} TaskDTOResponse
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /task/ [get]
func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Get MainPage handler")

	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to index.html for main page")
		return
	}

	responseHandler.HTMLResponse(html)

}

