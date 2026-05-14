package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
	core_http_types "github.com/AMmetro/ODRProcesing/internal/core/transport/http/types"
)

type UpdateUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("FullName must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("PhoneNumber must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must startswith '+' symbol")
			}
		}
	}
	return nil
}

type UpdateUserAgentResponse UserDTOResponse

// @Summary Update User
// @Description Update an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} UserDTOResponse
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /user/{id} [patch]
func (h *UsersHTTPHandler) UpdateUserRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Update UserAgent handler")

	userId, err := core_http_request.GetIntPathValue(r, "id")
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
	return domain.NewUserAgentPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
