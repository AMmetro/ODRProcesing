package tasks_service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	// Создаем задачу в БД
	newTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	fmt.Println("===========================11111111111111===================")

	// Сообщение для Kafka
	taskMessage := map[string]interface{}{
		"task_id":     newTask.ID,
		"author_id":   newTask.AuthorUserId,
		"title":       newTask.Title,
		"description": newTask.Description,
		"completed":   newTask.Completed,
		"created_at":  newTask.CreatedAt,
		"event_type":  "task.created",
	}

	log.Println("SENDING MESSAGE TO KAFKA")

	// Отправка сообщения в Kafka
	err = s.kafkaProducer.SendMessage(
		context.Background(),
		"tasks-events",
		strconv.Itoa(newTask.ID),
		taskMessage,
	)

	if err != nil {
		log.Printf("error sending task creation message to kafka: %v", err)
	} else {
		log.Println("MESSAGE SENT TO KAFKA")
	}

	return newTask, nil
}
