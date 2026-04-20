package users_transport_http

import (
	"context"
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUserAgent(
		ctx context.Context,
		user domain.UserAgent,
	) (domain.UserAgent, error)
}

func NewUserHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUserAgent,
		},
	}
}
