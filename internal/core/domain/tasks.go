package domain

import (
	"fmt"
	"time"

	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title        string
	Description  *string
	Completed    bool
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
		if descriptionLen < 10 || descriptionLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				`"CompletedAt" can't be nil if "Completed" == true: %w`,
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				`"CompletedAt" can't be before "CreatedAt": %w`,
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				`"CompletedAt" must be nil if "Completed" == false: %w`,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
