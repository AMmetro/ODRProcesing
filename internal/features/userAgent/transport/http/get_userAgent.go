package users_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type GetUsersAgentResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUserAgent(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetUserAgent handler")

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId")
		return
	}

	userDomain, err := h.usersService.GetUserAgent(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersAgentResponse(UserDtoFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
