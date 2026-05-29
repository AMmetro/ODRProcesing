# Kafka Integration Guide

## Overview

This project uses Kafka for asynchronous inter-service communication. The **processing** service produces task events, and the **reservation** service consumes them.

## Architecture

```
┌────────────────┐ ┌────────────────┐
│ Processing │ │ Reservation │
│ Service │ ──────── tasks-events ──────────────────────────► │ Service │
│ │ │ │
│ │ ◄────── tasks-responses (temporary consumer) ──── │ │
└────────────────┘ └────────────────┘
```

## Kafka Topics

### `tasks-events`

**Purpose**: Task lifecycle events  
**Producer**: Processing Service  
**Consumers**: Reservation Service (group: `reservation-group`)  

**Message Format (JSON)**:
```json
{
  "task_id": 123,
  "author_id": 45,
  "title": "Example task",
  "description": "Task description",
  "completed": false,
  "created_at": "2024-01-01T10:00:00Z",
  "event_type": "task.created"
}
```

## Implementation Details

### Producer (Processing Service)

**File**: `services/processing/internal/features/tasks/service/create_task.go`

The processing service sends a message to Kafka whenever a new task is created:

```go
go func() {
    taskMessage := map[string]interface{}{
        "task_id":     newTask.ID,
        "author_id":   newTask.AuthorUserId,
        "title":       newTask.Title,
        "description": newTask.Description,
        "completed":   newTask.Completed,
        "created_at":  newTask.CreatedAt,
        "event_type":  "task.created",
    }

    if err := s.kafkaProducer.SendMessage(ctx, "tasks-events", strconv.Itoa(newTask.ID), taskMessage); err != nil {
        log.Printf("error sending task creation message to kafka: %v", err)
    }
}()
```

**Key Features**:
- Message is sent asynchronously (doesn't block task creation)
- Task ID used as Kafka message key (ensures order of related messages)
- Errors are logged but don't affect the main operation

### Consumer (Reservation Service)

**File**: `services/reservation/cmd/reservation/main.go`

The reservation service initializes a Kafka consumer on startup:

```go
kafkaConsumer, err := messaging.NewKafkaConsumer(
    []string{"kafka:9092"},
    "reservation-group",
    "tasks-events",
)

kafkaConsumer.SetMessageHandler(func(message []byte) error {
    return service.ProcessTaskMessage(context.Background(), message)
})

if err := kafkaConsumer.Start(ctx); err != nil {
    log.Fatalf("failed to start kafka consumer: %v", err)
}
```

**Processing Logic**: `services/reservation/internal/features/reservation/service/service.go`

```go
func (s *ReservationService) ProcessTaskMessage(ctx context.Context, message []byte) error {
    var taskEvent map[string]interface{}
    if err := messaging.UnmarshalMessage(message, &taskEvent); err != nil {
        return err
    }
    
    // Process the event
    // - Create reservations
    // - Notify users
    // - Log events
    
    return nil
}
```

## Kafka Libraries

The project uses **segmentio/kafka-go** for Kafka client implementation.

**Shared Layer**: `shared/pkg/core/messaging/`
- `kafka_producer.go` - Producer implementation
- `kafka_consumer.go` - Consumer implementation

## Configuration

### Kafka Brokers

Default broker: `kafka:9092`

To change, update in:
- **Processing**: `services/processing/cmd/processing/main.go`
  ```go
  kafkaProducer, err := core_messaging.NewKafkaProducer([]string{"kafka:9092"})
  ```

- **Reservation**: `services/reservation/cmd/reservation/main.go`
  ```go
  kafkaConsumer, err := messaging.NewKafkaConsumer(
      []string{"kafka:9092"},
      "reservation-group",
      "tasks-events",
  )
  ```

### Consumer Group

- **Group ID**: `reservation-group`
- **Topic**: `tasks-events`

Multiple consumer instances can use the same group ID for load balancing.

## Running with Docker Compose

The docker-compose.yaml already includes Kafka and Zookeeper:

```bash
docker compose up -d
```

This starts:
- Kafka broker (port 9092)
- Zookeeper (port 2181)
- Both microservices

## Error Handling

### Producer Errors
- If Kafka is unavailable, task creation still succeeds (message is lost)
- Errors are logged but don't affect the main operation
- For critical scenarios, implement retry logic or dead letter queue

### Consumer Errors
- Messages with processing errors are not committed
- Consumer will retry the message on next restart
- Invalid JSON messages are skipped after logging

## Adding New Events

To add a new event type:

1. **Producer**: Add new message struct in processing service
2. **Topic**: Use same or create new topic
3. **Consumer**: Add handler in reservation service
4. **Message Format**: Document the JSON schema

## Monitoring

Monitor Kafka with:

```bash
# List topics
docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Describe topic
docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic tasks-events

# Console consumer (for debugging)
docker compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic tasks-events --from-beginning
```

## Performance Tuning

Adjust in producer/consumer config:
- **Batch size**: `kafka.Writer.WriteBackoffMin/Max`
- **Compression**: Enable in writer config
- **Replication factor**: Consider for production deployments
