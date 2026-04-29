package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Completed    bool    `json:"completed"` // todo exclude
	AuthorUserId int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserId int        `json:"author_user_id"`
}

func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Create Task handler")

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	taskDomain := domain.NewTaskInitialized(
		request.Title,
		request.Description,
		request.Completed,
		request.AuthorUserId)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}
	response := taskDTOfromDomain(taskDomain)
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func taskDTOfromDomain(taskDomain domain.Task) CreateTaskResponse {
	return CreateTaskResponse{
		ID:           taskDomain.ID,
		Version:      taskDomain.Version,
		Title:        taskDomain.Title,
		Description:  taskDomain.Description,
		Completed:    taskDomain.Completed,
		CreatedAt:    taskDomain.CreatedAt,
		CompletedAt:  taskDomain.CompletedAt,
		AuthorUserId: taskDomain.AuthorUserId,
	}
}
