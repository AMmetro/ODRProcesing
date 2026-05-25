package reservation_service

import (
	"context"

	reservation_domain "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/domain"
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

func (s *ReservationService) GetItems(ctx context.Context) ([]ReservationItem, error) {
	return s.repository.GetItems(ctx)
}
