package users_transport_http

import (
	"net/http"

	domain_utils "github.com/AMmetro/ODRProcesing/internal/core/domain/utils"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// @Summary Get Users
// @Description Retrieve a list of users with optional limit and offset
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to retrieve" minimum(1) maximum(100)
// @Param offset query int false "Number of users to skip" minimum(0)
// @Success 200 {array} UserDTOResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetUsers handler")

	limit, offset, _, err := domain_utils.GetUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}
	response := GetUsersResponse(UsersDtoFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusCreated)

}
