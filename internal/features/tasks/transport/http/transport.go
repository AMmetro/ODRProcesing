package users_transport_http

import (
	"context"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTasks(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/tasks",
		// 	Handler: h.CreateUserAgent,
		// },
	}
}
