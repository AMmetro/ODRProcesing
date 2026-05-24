package web_transport_http

import (
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			/* disabled for awoid conflict with other routes in root path */
			// Method:  http.MethodGet,
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
