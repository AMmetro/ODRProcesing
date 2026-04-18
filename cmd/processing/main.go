package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_server "github.com/AMmetro/ODRProcesing/internal/core/transport/http/server"
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
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()
	logger.Debug("application started")

	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
	)
	httpServer.RegisterAPIRoute(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
	}
}
