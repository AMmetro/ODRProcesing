package tasks_service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {

	log := core_logger.FromContext(ctx)

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	newTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	taskMessage := map[string]interface{}{
		"task_id":     newTask.ID,
		"author_id":   newTask.AuthorUserId,
		"title":       newTask.Title,
		"description": newTask.Description,
		"completed":   newTask.Completed,
		"created_at":  newTask.CreatedAt,
		"event_type":  "task.created",
	}

	err = s.kafkaProducer.SendMessage(
		context.Background(),
		"tasks-events",
		strconv.Itoa(newTask.ID),
		taskMessage,
	)

	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	} else {
		log.Debug("Send message to Kafka")
	}

	replyTopic := "tasks-responses"
	groupID := "processing-replies-" + strconv.Itoa(newTask.ID) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	kafkaConsumer, err := messaging.NewKafkaConsumer([]string{"localhost:9092"}, groupID, replyTopic)
	if err != nil {
		log.Debug("failed to init temporary kafka consumer: ")
		return newTask, nil
	}
	defer kafkaConsumer.Close()

	replyCh := make(chan map[string]interface{}, 1)

	kafkaConsumer.SetMessageHandler(func(message []byte) error {
		var msg map[string]interface{}
		if err := messaging.UnmarshalMessage(message, &msg); err != nil {
			// log.Printf("failed to unmarshal reply message: %v", err)
			return err
		}

		// Простая корреляция по task_id
		if idRaw, ok := msg["task_id"]; ok {
			var id int
			switch v := idRaw.(type) {
			case float64:
				id = int(v)
			case int:
				id = v
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					id = parsed
				}
			}

			if id == newTask.ID {
				select {
				case replyCh <- msg:
				default:
				}
			}
		}

		return nil
	})

	// Таймаут ожидания ответа
	waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := kafkaConsumer.Start(waitCtx); err != nil {
		// log.Printf("failed to start temporary kafka consumer: %v", err)
		return newTask, nil
	}

	select {
	case reply := <-replyCh:
		// log.Printf("received reservation reply for task %d: %v", newTask.ID, reply)

		// Если в ответе есть статус резервации, можно обновить задачу
		if status, ok := reply["reservation_status"].(string); ok && status == "ok" {
			newTask.Completed = true
			if updated, err := s.tasksRepository.UpdateTask(context.Background(), newTask); err == nil {
				newTask = updated
			} else {
				// log.Debug("failed to update task after reservation reply: %v", err)
			}
		}
	case <-waitCtx.Done():
		// log.Printf("timed out waiting for reservation reply for task %d", newTask.ID)
	}

	return newTask, nil
}
