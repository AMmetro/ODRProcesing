// package statistics_transport_http

// import (
// 	"context"
// 	"net/http"

// 	"github.com/AMmetro/ODRProcesing/internal/core/domain"
// 	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
// )

// type StatisticsHTTPHandler struct {
// 	statisticsService StatisticsService
// }

// type StatisticsService interface {
// 	GetStatistics(
// 		ctx context.Context,
// 	) (domain.StatisticsSummary, error)
// }

// func NewStatisticsHTTPHandler(
// 	statisticsService StatisticsService,
// ) *StatisticsHTTPHandler {
// 	return &StatisticsHTTPHandler{
// 		statisticsService: statisticsService,
// 	}
// }

// func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
// 	return []core_http_server.Route{
// 		{
// 			Method:  http.MethodGet,
// 			Path:    "/statistics",
// 			Handler: h.GetStatistics,
// 		},
// 	}
// }

package statistics_transport_http

import (
	"context"
	"net/http"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
	) (domain.StatisticsSummary, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
