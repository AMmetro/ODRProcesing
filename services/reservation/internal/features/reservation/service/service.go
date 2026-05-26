package reservation_service

import (
	"context"

	reservation_domain "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/domain"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

type ReservationItem = reservation_domain.ReservationItem

type ReservationService struct {
	repository ReservationRepository
}

type ReservationRepository interface {
	GetItems(ctx context.Context) ([]ReservationItem, error)
}

func NewReservationService(
	repository ReservationRepository,
) *ReservationService {
	return &ReservationService{
		repository: repository,
	}
}

// ProcessTaskMessage обрабатывает сообщение из Kafka о новой задаче
// Это функция предназначена для использования в качестве message handler для Kafka consumer
func (s *ReservationService) ProcessTaskMessage(ctx context.Context, message []byte) error {

	var taskEvent map[string]interface{}
	if err := messaging.UnmarshalMessage(message, &taskEvent); err != nil {
		return err
	}

	// Здесь должна быть бизнес-логика обработки события создания задачи
	// Например:
	// - Создание резервации на основе задачи
	// - Уведомление пользователя
	// - Логирование события
	// и т.д.

	// Для демонстрации просто логируем событие
	// logger.Debug("Received task event", zap.Any("event", taskEvent))

	return nil
}

func (s *ReservationService) GetItems(ctx context.Context) ([]ReservationItem, error) {
	return s.repository.GetItems(ctx)
}
