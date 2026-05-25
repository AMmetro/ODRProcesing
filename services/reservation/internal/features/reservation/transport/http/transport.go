package reservation_transport_http

import (
	"context"
	"encoding/json"
	"net/http"

	reservation_service "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/service"
	core_http_server "github.com/AMmetro/ODRProcesing/shared/pkg/core/transport/http/server"
)

type ReservationHTTPHandler struct {
	reservationService ReservationService
}

type ReservationService interface {
	GetItems(ctx context.Context) ([]reservation_service.ReservationItem, error)
}

func NewReservationHTTPHandler(
	reservationService ReservationService,
) *ReservationHTTPHandler {
	return &ReservationHTTPHandler{
		reservationService: reservationService,
	}
}

func (h *ReservationHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/reservation",
			Handler: h.GetItems,
		},
	}
}

func (h *ReservationHTTPHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.reservationService.GetItems(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
