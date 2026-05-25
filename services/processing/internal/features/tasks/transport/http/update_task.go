package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	core_http_request "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/request"
	core_http_response "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/response"
	core_http_types "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/types"
)

type PatchTaskResponse TaskDTOResponse

type UpdateTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *UpdateTaskRequest) Validate() error {

	if !r.Title.Set && !r.Description.Set && !r.Completed.Set {
		return fmt.Errorf("empty patch")
	}

	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("TitleLen must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("DescriptionLen must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed can't be NULL")
		}
	}

	return nil
}

// @Summary Update Task
// @Description Update an existing task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body UpdateTaskRequest true "Task update data"
// @Success 200 {object} TaskDTOResponse
// @Failure 404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Bad Request"
// @Router /task/{id} [patch]
func (h *TasksHTTPHandler) UpdateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Update Tasks handler")

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "can't get taskId from path")
		return
	}

	var request UpdateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.UpdateTask(ctx, userId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to update user")
		return
	}
	response := PatchTaskResponse(TaskDtoFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func taskPatchFromRequest(request UpdateTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}

