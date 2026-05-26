package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	reservation_repository "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/repository/inmemory"
	reservation_service "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/service"
	reservation_transport_http "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/transport/http"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

func main() {
	addr := getEnv("HTTP_ADDR", ":8091")

	repository := reservation_repository.NewReservationRepository()
	service := reservation_service.NewReservationService(repository)
	handler := reservation_transport_http.NewReservationHTTPHandler(service)

	// Initialize Kafka Consumer
	kafkaConsumer, err := messaging.NewKafkaConsumer(
		[]string{"localhost:9092"},
		"reservation-group",
		"tasks-events",
	)
	if err != nil {
		log.Fatalf("failed to init kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// Setup message handler
	kafkaConsumer.SetMessageHandler(func(message []byte) error {
		return service.ProcessTaskMessage(context.Background(), message)
	})

	// Start consumer
	ctx := context.Background()
	if err := kafkaConsumer.Start(ctx); err != nil {
		log.Fatalf("failed to start kafka consumer: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/reservation", handler.GetItems)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("reservation service starting on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("reservation server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("reservation service shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("reservation server graceful shutdown failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
