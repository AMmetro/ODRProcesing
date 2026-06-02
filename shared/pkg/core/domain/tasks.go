package domain

import (
	"fmt"
	"time"

	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
)

type TaskStatus string

const (
	TaskStatusUnconfirmed TaskStatus = "unconfirmed"
	TaskStatusPending     TaskStatus = "pending"
	TaskStatusConfirmed   TaskStatus = "confirmed"
	TaskStatusRejected    TaskStatus = "rejected"
	TaskStatusFailed      TaskStatus = "failed"
	TaskStatusTimeout     TaskStatus = "timeout"
)

type Task struct {
	ID      int
	Version int

	Title        string
	Description  *string
	Completed    bool
	Status       TaskStatus
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserId int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	status TaskStatus,
	authorUserId int,
	createdAt time.Time,
	completedAt *time.Time,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		Status:       status,
		AuthorUserId: authorUserId,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
	}
}

func NewTaskInitialized(
	title string,
	description *string,
	completed bool,
	authorUserId int,
) Task {
	return Task{
		ID:           UninitializedID,
		Version:      UninitializedVersion,
		Title:        title,
		Description:  description,
		Completed:    false,
		Status:       "unconfirmed",
		CreatedAt:    time.Now(),
		CompletedAt:  nil,
		AuthorUserId: authorUserId,
	}
}

func (t Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `title` len: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}
	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {

	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		taskIsCompleted := *patch.Completed.Value

		if taskIsCompleted {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}

		tmp.Completed = taskIsCompleted
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(title Nullable[string], description Nullable[string], completed Nullable[bool]) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (t TaskPatch) Validate() error {
	if t.Title.Set && t.Title.Value == nil {
		return fmt.Errorf("Title can`t be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	if t.Completed.Set && t.Completed.Value == nil {
		return fmt.Errorf("Completed can`t be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}
	if t.CompletedAt == nil {
		return nil
	}

	duration := t.CompletedAt.Sub(t.CreatedAt)
	return &duration
}
