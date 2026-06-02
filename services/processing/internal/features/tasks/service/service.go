package tasks_service

import (
	"context"
	"errors"
	"sync"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

type ResponseResult struct {
	Status domain.TaskStatus
	Error  error
}

type TasksService struct {
	tasksRepository TasksRepository
	kafkaProducer   messaging.KafkaProducer

	// Для синхронизации ответов от других микросервисов
	responseMutex sync.RWMutex
	responses     map[int]chan ResponseResult
}

type TasksRepository interface {
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
		task domain.Task,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		taskId int,
	) error
}

func NewTasksService(
	tasksRepository TasksRepository,
	kafkaProducer messaging.KafkaProducer,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
		kafkaProducer:   kafkaProducer,
		responses:       make(map[int]chan ResponseResult),
	}
}

// HandleTaskResponse process response from other microservices (reservation)
// and extracts result to coresponding chanale
func (s *TasksService) HandleTaskResponse(ctx context.Context, response map[string]interface{}) error {
	// Extract task_id from response
	taskID, ok := response["task_id"]
	if !ok {
		return nil
	}

	// Конвертируем task_id в int
	var id int
	switch v := taskID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return nil
	}

	// Get status from response
	statusStr, ok := response["reservation_status"].(string)
	if !ok {
		statusStr = string(domain.TaskStatusRejected)
	}

	status := domain.TaskStatus(statusStr)

	// Get error message from respone
	var responseErr error
	if errMsg, ok := response["error"].(string); ok && errMsg != "" {
		responseErr = errors.New(errMsg)
	}

	// send result to chanale
	s.responseMutex.RLock()
	responseChan, exists := s.responses[id]
	s.responseMutex.RUnlock()

	if exists {
		select {
		case responseChan <- ResponseResult{
			// Status: domain.TaskStatus(status),
			Status: status,
			Error:  responseErr,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
