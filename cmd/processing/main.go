package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/AMmetro/ODRProcesing/internal/core/transport/http/middleware"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
	users_postgres_repository "github.com/AMmetro/ODRProcesing/internal/features/userAgent/repository/postgress"
	users_service "github.com/AMmetro/ODRProcesing/internal/features/userAgent/service"
	users_transport_http "github.com/AMmetro/ODRProcesing/internal/features/userAgent/transport/http"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application pool:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("Faled to init postges poll", zap.Error(err))
	}
	defer pool.Close()

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	usersRoutes := usersTransportHTTP.Routes()
	apiVersionRouter.RegisterRoutes(usersRoutes...)
	httpServer.RegisterAPIRoute(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
	}
}
