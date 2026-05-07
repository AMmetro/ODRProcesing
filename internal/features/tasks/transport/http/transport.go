package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		limit *int,
		ofset *int,
		userId *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		taskId int,
	) (domain.Task, error)

	UpdateTask(
		ctx context.Context,
		taskId int,
		taskPatch domain.TaskPatch,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		taskId int,
	) error
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
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/task/{id}",
			Handler: h.UpdateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/task/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/task/{id}",
			Handler: h.DeleteTask,
		},
	}
}
