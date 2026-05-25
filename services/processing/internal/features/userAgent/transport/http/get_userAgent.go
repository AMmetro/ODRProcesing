package users_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

type GetUsersAgentResponse UserDTOResponse

// @Summary Get User
// @Description Retrieve a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserDTOResponse
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /user/{id} [get]
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

