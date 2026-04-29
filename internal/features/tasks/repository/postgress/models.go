package tasks_postgres_repository

import "time"

type TasksModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	AuthorUserId int
	CreatedAt    time.Time
	CompletedAt  *time.Time
}
