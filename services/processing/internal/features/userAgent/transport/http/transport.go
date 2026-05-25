package users_transport_http

import (
	"context"
	"net/http"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_http_server "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUserAgent(
		ctx context.Context,
		user domain.UserAgent,
	) (domain.UserAgent, error)

	UpdateUserAgent(
		ctx context.Context,
		id int,
		user domain.UserAgentPatch,
	) (domain.UserAgent, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.UserAgent, error)

	GetUserAgent(
		ctx context.Context,
		id int,
	) (domain.UserAgent, error)

	DeleteUserAgent(
		ctx context.Context,
		id int,
	) error
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
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			// Middleware: []core_http_middleware.Middleware{
			// 	core_http_middleware.Dummy("get users middleware"),
			// },
		},
		{
			Method:  http.MethodGet,
			Path:    "/user/{id}",
			Handler: h.GetUserAgent,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUserAgent,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/user/{id}",
			Handler: h.DeleteUserAgentRequest,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/user/{id}",
			Handler: h.UpdateUserRequest,
		},
	}
}

