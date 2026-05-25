package tasks_transport_http

import (
	"time"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserId int        `json:"user_id"`
}

func TaskDtoFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Description:  task.Description,
		Title:        task.Title,
		Completed:    task.Completed,
		CompletedAt:  task.CompletedAt,
		CreatedAt:    task.CreatedAt,
		AuthorUserId: task.AuthorUserId,
	}
}

func TasksDtoFromDomains(tasks []domain.Task) []TaskDTOResponse {
	res := make([]TaskDTOResponse, 0, len(tasks))
	for _, task := range tasks {
		dto := TaskDtoFromDomain(task)
		res = append(res, dto)
	}
	return res
}

