package users_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"

	"go.uber.org/zap"
)

type DeleteUserAgentResponse UserDTOResponse

func (h *UsersHTTPHandler) DeleteUserAgentRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "can't get userId from path")
		return
	}

	log.Debug("invoke DeleteUserAgent handler, userId", zap.Int("userId", userId))

	err = h.usersService.DeleteUserAgent(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete userAgent")
		return
	}

	responseHandler.NoContentResponse()
}
