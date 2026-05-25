package tasks_transport_http

import (
	"net/http"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Completed    bool    `json:"completed"` // todo exclude
	AuthorUserId int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TaskDTOResponse

/* todo for future options */
// type CreateTaskResponse struct {
// 	Task TaskDTOResponse
// 	NewOptions []any
// }

// @Summary Create Task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task data"
// @Success 200 {object} TaskDTOResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /tasks [post]
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
	response := CreateTaskResponse(TaskDtoFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

