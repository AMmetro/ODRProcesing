package users_transport_http

import (
	"net/http"

	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"

	"go.uber.org/zap"
)

type DeleteUserAgentResponse UserDTOResponse

// @Summary Delete User
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /user/{id} [delete]
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

