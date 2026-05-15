package users_transport_http

import (
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+1234567890"`
}

type CreateUserResponse UserDTOResponse

// @Summary Create User
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} UserDTOResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /users [post]
func (h *UsersHTTPHandler) CreateUserAgent(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUserAgent handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUserAgent(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}
	response := CreateUserResponse(UserDtoFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.UserAgent {
	return domain.NewUserAgentInitialized(dto.FullName, dto.PhoneNumber)
}
