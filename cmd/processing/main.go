package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AMmetro/ODRProcesing/internal/core/config"
	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_pgx_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/AMmetro/ODRProcesing/internal/core/transport/http/middleware"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/AMmetro/ODRProcesing/internal/features/statistics/repository/postgres"
	statistics_service "github.com/AMmetro/ODRProcesing/internal/features/statistics/service"
	statistics_transport_http "github.com/AMmetro/ODRProcesing/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/AMmetro/ODRProcesing/internal/features/tasks/repository/postgress"
	tasks_service "github.com/AMmetro/ODRProcesing/internal/features/tasks/service"
	tasks_transport_http "github.com/AMmetro/ODRProcesing/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/AMmetro/ODRProcesing/internal/features/userAgent/repository/postgress"
	users_service "github.com/AMmetro/ODRProcesing/internal/features/userAgent/service"
	users_transport_http "github.com/AMmetro/ODRProcesing/internal/features/userAgent/transport/http"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "github.com/AMmetro/ODRProcesing/docs"
)

func main() {

	// ============================================================================
	// 1. Load CONFIG ENV if exists
	// ============================================================================
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
	}

	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	// ============================================================================
	// 2. НАСТРОЙКА GRACEFUL SHUTDOWN <-- (Ctrl+C) or  Docker/k8s
	// ============================================================================
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	// ============================================================================
	// 3. LOGGER INITIALIZE
	// ============================================================================
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application pool:", err)
		os.Exit(1)
	}
	defer logger.Close() // Обязательно закрываем файл лога при выходе

	logger.Debug("application time zone", zap.Any("timezone", time.Local))
	logger.Debug("initializing feature", zap.String("feature", "users"))
	logger.Debug("initializing feature", zap.String("feature", "tasks"))

	// ============================================================================
	// 4. CREATE CONNECTION POOL
	// ============================================================================

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("Faled to init postges poll", zap.Error(err))
	}
	defer pool.Close()

	// ============================================================================
	// 5. SETTING LAYERS (DEPENDENCY INJECTION):  DB → Service → Transport (HTTP handlers)
	// ============================================================================
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	// ============================================================================
	// 6. SETUP HTTP SERVER
	// ============================================================================
	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	// ============================================================================
	// 7. CREATE ROUTER API VERSION /api/v1/
	// ============================================================================
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)

	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	// apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2,
	// 	core_http_middleware.Dummy("api v2 middleware"),
	// )
	// apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	// setup router - to HTTP server
	httpServer.RegisterAPIRoute(apiVersionRouterV1)
	httpServer.RegisterSwagger()

	// httpServer.RegisterAPIRoute(apiVersionRouterV1, apiVersionRouterV2)

	// ============================================================================
	// 8. START SERVER AND WAIT SIGNAL GRACEFUL SHUTDOWN
	// ============================================================================

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
	}
}
