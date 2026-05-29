package reservation_repository

import (
	"context"
	"fmt"
	"time"

	reservation_domain "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/domain"
)

type ReservationRepository struct {
	items []reservation_domain.ReservationItem
}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{
		items: []reservation_domain.ReservationItem{
			{ID: 1, Name: "order-status", Description: "Order processing state from reservation service"},
			{ID: 2, Name: "health", Description: "Health and readiness information"},
			{ID: 3, Name: "metrics", Description: "Sample metrics data for the reservation service"},
		},
	}
}

func (r *ReservationRepository) GetItems(ctx context.Context) ([]reservation_domain.ReservationItem, error) {
	result := make([]reservation_domain.ReservationItem, len(r.items))
	copy(result, r.items)
	return result, nil
}

func (r *ReservationRepository) GetReservations(ctx context.Context, taskID int) (int, error) {

	fmt.Println("create new product reservation - validate permission")

	var approvalNumber int

	select {
	case <-time.After(3 * time.Second):
		for _, item := range r.items {
			if item.ID == taskID {
				approvalNumber = item.ID
				break
			}
		}
		fmt.Println("product reservation created")
		return approvalNumber, nil

	case <-ctx.Done():
		return 0, ctx.Err()
	}

}
