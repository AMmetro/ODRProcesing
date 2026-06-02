package tasks_service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
	core_logger "github.com/AMmetro/ODRProcesing/shared/pkg/core/logger"
)

const createConstantTimeout = 5 * time.Second

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {

	log := core_logger.FromContext(ctx)

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	task.Status = domain.TaskStatusUnconfirmed // add default status unconfirmed
	newTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	log.Debug("task created in DB with unconfirmed status", zap.Int("task_id", newTask.ID))

	// Регистрируем канал для ожидания ответа
	// task_id -> same partition -> same consumer instance IN PROD
	responseChan := make(chan ResponseResult, 1)
	s.responseMutex.Lock()
	s.responses[newTask.ID] = responseChan // responses[102] = responseChan
	s.responseMutex.Unlock()

	defer func() {
		s.responseMutex.Lock()
		delete(s.responses, newTask.ID)
		s.responseMutex.Unlock()
	}()

	// Send event to Kafka for reservation services with default status unconfirmed
	taskMessage := map[string]interface{}{
		"task_id":     newTask.ID,
		"author_id":   newTask.AuthorUserId,
		"title":       newTask.Title,
		"description": newTask.Description,
		"created_at":  newTask.CreatedAt,
		"event_type":  "task.created",
	}

	if err := s.kafkaProducer.SendMessage(
		ctx,
		"tasks-events",
		strconv.Itoa(newTask.ID),
		taskMessage,
	); err != nil {
		log.Error("failed to send task.created event to kafka",
			zap.Int("task_id", newTask.ID),
			zap.Error(err),
		)
		// Mark task as failed if we can't send event to Kafka
		newTask.Status = domain.TaskStatusFailed
		if _, updateErr := s.tasksRepository.UpdateTask(ctx, newTask); updateErr != nil {
			log.Error("failed to mark task as failed",
				zap.Int("task_id", newTask.ID),
				zap.Error(updateErr),
			)
		}
		return newTask, fmt.Errorf("send task.created event: %w", err)
	}

	newTask.Status = domain.TaskStatusPending
	if _, updateStatus := s.tasksRepository.UpdateTask(ctx, newTask); updateStatus != nil {
		log.Error("failed to mark task as pending",
			zap.Int("task_id", newTask.ID),
			zap.Error(updateStatus),
		)
	}

	log.Debug("task.created event sent to kafka, waiting for confirmation", zap.Int("task_id", newTask.ID))

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {

	case response := <-responseChan:
		if response.Error != nil {
			log.Error("reservation service returned error",
				zap.Int("task_id", newTask.ID),
				zap.Error(response.Error),
			)
			newTask.Status = domain.TaskStatusRejected
			if _, err := s.tasksRepository.UpdateTask(ctx, newTask); err != nil {
				log.Error("failed to update rejected task status",
					zap.Int("task_id", newTask.ID),
					zap.Error(err),
				)
			}
			return newTask, fmt.Errorf(
				"reservation service error: %w",
				response.Error,
			)
		}
		if response.Status == domain.TaskStatusConfirmed {
			newTask.Status = domain.TaskStatusConfirmed
			confirmedTask, err := s.tasksRepository.UpdateTask(ctx, newTask)
			if err != nil {
				log.Error("failed to confirm task",
					zap.Int("task_id", newTask.ID),
					zap.Error(err),
				)
				return domain.Task{}, fmt.Errorf("confirm task: %w", err)
			}

			log.Debug("task confirmed",
				zap.Int("task_id", newTask.ID),
			)
			return confirmedTask, nil
		}

		return newTask, fmt.Errorf(
			"invalid response from reservation service: status=%q error=nil",
			response.Status,
		)

	case <-timeoutCtx.Done():
		newTask.Status = domain.TaskStatusTimeout
		if _, err := s.tasksRepository.UpdateTask(ctx, newTask); err != nil {
			log.Error("failed to update timeout task status",
				zap.Int("task_id", newTask.ID),
				zap.Error(err),
			)
		}

		log.Error("timeout waiting reservation confirmation",
			zap.Int("task_id", newTask.ID),
		)

		return newTask, fmt.Errorf(
			"task confirmation timeout: %w",
			core_errors.ErrTimeout,
		)
	}
}
