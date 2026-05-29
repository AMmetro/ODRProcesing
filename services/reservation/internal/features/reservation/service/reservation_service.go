package reservation_service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	reservation_domain "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/domain"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

type ReservationItem = reservation_domain.ReservationItem

type ReservationService struct {
	repository ReservationRepository
	producer   messaging.KafkaProducer
}

type ReservationRepository interface {
	GetItems(ctx context.Context) ([]ReservationItem, error)
	GetReservations(ctx context.Context, taskID int) (int, error)
}

func NewReservationService(
	repository ReservationRepository,
	producer messaging.KafkaProducer,
) *ReservationService {
	return &ReservationService{
		repository: repository,
		producer:   producer,
	}
}

// func (s *ReservationService) createReservation(ctx context.Context, taskID int) error {
// 	s.repository.GetReservations(ctx) // taskID
// 	return nil
// }

func (s *ReservationService) sendReply(
	ctx context.Context,
	taskID int,
	status string,
	errMsg string,
) error {

	replyTopic := "tasks-responses"

	replyMessage := map[string]interface{}{
		"task_id":            taskID,
		"reservation_status": status,
		"event_type":         "reservation.completed",
		"completed_at":       time.Now(),
		"error":              errMsg,
	}

	return s.producer.SendMessage(
		ctx,
		replyTopic,
		strconv.Itoa(taskID),
		replyMessage,
	)
}

type TaskCreatedEvent struct {
	TaskID int `json:"task_id"`
}

// Kafka consumer message handler
func (s *ReservationService) ProcessTaskMessage(ctx context.Context, message []byte) error {

	var taskEvent TaskCreatedEvent
	if err := messaging.UnmarshalMessage(message, &taskEvent); err != nil {
		_ = s.sendReply(ctx, 0, "failed", "invalid_message")
		return fmt.Errorf("unmarshal task event: %w", err)
	}

	taskID := taskEvent.TaskID

	aprovalId, err := s.repository.GetReservations(ctx, taskID)

	if err != nil {

		// todo rolback
		_ = s.sendReply(ctx, taskID, "failed", err.Error())
		return fmt.Errorf("create reservation: %w", err)
	}

	_ = s.sendReply(ctx, aprovalId, "ok", "")

	return nil
}

func (s *ReservationService) GetItems(ctx context.Context) ([]ReservationItem, error) {
	return s.repository.GetItems(ctx)
}
