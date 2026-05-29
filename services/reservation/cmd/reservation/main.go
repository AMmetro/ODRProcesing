package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	reservation_repository "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/repository/inmemory"
	reservation_service "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/service"
	reservation_transport_http "github.com/AMmetro/ODRProcesing/services/reservation/internal/features/reservation/transport/http"
	"github.com/AMmetro/ODRProcesing/shared/pkg/core/messaging"
)

func main() {

	// =========================================================
	// CONFIG
	// =========================================================

	httpAddr := getEnv("HTTP_ADDR", ":8091")

	kafkaBrokers := strings.Split(
		getEnv("KAFKA_BROKERS", "localhost:9092"),
		",",
	)

	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "reservation-group")

	taskEventsTopic := getEnv("KAFKA_TASK_EVENTS_TOPIC", "tasks-events")

	// =========================================================
	// ROOT CONTEXT
	// =========================================================

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// =========================================================
	// REPOSITORY
	// =========================================================

	repository := reservation_repository.NewReservationRepository()

	// =========================================================
	// KAFKA PRODUCER
	// =========================================================

	kafkaProducer, err := messaging.NewKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatalf("failed to init kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// =========================================================
	// SERVICE
	// =========================================================

	service := reservation_service.NewReservationService(
		repository,
		kafkaProducer,
	)

	// =========================================================
	// HTTP HANDLER
	// =========================================================

	handler := reservation_transport_http.NewReservationHTTPHandler(service)

	// =========================================================
	// KAFKA CONSUMER
	// =========================================================

	kafkaConsumer, err := messaging.NewKafkaConsumer(
		kafkaBrokers,
		kafkaGroupID,
		taskEventsTopic,
	)
	if err != nil {
		log.Fatalf("failed to init kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	kafkaConsumer.SetMessageHandler(func(message []byte) error {
		return service.ProcessTaskMessage(appCtx, message)
	})

	if err := kafkaConsumer.Start(appCtx); err != nil {
		log.Fatalf("failed to start kafka consumer: %v", err)
	}

	// =========================================================
	// HTTP SERVER
	// =========================================================

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/reservation", handler.GetItems)

	srv := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("reservation service started on %s", httpAddr)

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf("reservation server failed: %v", err)
		}
	}()

	// =========================================================
	// GRACEFUL SHUTDOWN
	// =========================================================

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("reservation service shutting down...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("reservation service stopped")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
