package domain

import "time"

type Statistics struct {
	TaskCreated               int
	TasksCompleted            int
	TaskCompletedRate         *float64
	TaskAverageComplitionTime *time.Duration
}

func NewStatistics(
	taskCreated int,
	tasksCompleted int,
	taskCompletedRate *float64,
	taskAverageComplitionTime *time.Duration,
) Statistics {
	return Statistics{
		TaskCreated:               taskCreated,
		TasksCompleted:            tasksCompleted,
		TaskCompletedRate:         taskCompletedRate,
		TaskAverageComplitionTime: taskAverageComplitionTime,
	}
}

