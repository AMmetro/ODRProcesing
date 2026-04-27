package users_transport_http

import (
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
	core_http_types "github.com/AMmetro/ODRProcesing/internal/core/transport/http/types"
	core_http_utils "github.com/AMmetro/ODRProcesing/internal/core/transport/http/utils"
)

type UpdateUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

type UpdateUserAgentResponse UserDTOResponse

func (h *UsersHTTPHandler) UpdateUserRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Update UserAgent handler")

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "can't get userId from path")
		return
	}

	var request UpdateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	userAgentPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.UpdateUserAgent(ctx, userId, userAgentPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to update user")
		return
	}
	response := UpdateUserAgentResponse(UserDtoFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func userPatchFromRequest(request UpdateUserRequest) domain.UserAgentPatch {
	return domain.UserAgentPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
