package tasks_service

import (
	"context"
	"sync"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

type ResponseResult struct {
	Status string
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

// HandleTaskResponse обрабатывает ответ от других микросервисов (например, reservation)
// и отправляет результат в соответствующий канал
func (s *TasksService) HandleTaskResponse(ctx context.Context, response map[string]interface{}) error {
	// Извлекаем task_id из ответа
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

	// Получаем статус из ответа
	status, ok := response["reservation_status"].(string)
	if !ok {
		status = "error"
	}

	// Отправляем результат в канал
	s.responseMutex.RLock()
	responseChan, exists := s.responses[id]
	s.responseMutex.RUnlock()

	if exists {
		select {
		case responseChan <- ResponseResult{
			Status: status,
			Error:  nil,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
