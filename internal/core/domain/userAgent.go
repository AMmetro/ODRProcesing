package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
)

type UserAgent struct {
	ID          int
	FullName    string
	Version     int
	PhoneNumber *string
}

func NewUserAgent(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) UserAgent {
	return UserAgent{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserAgentInitialized(
	fullName string,
	phoneNumber *string,
) UserAgent {
	return UserAgent{
		ID:          UninitializedID,
		Version:     UninitializedVersion,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (u *UserAgent) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}
	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}
		re := regexp.MustCompile(`^\+?[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil
}

func (u *UserAgent) ApplyPatch(patch UserAgentPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

type UserAgentPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserAgentPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("Full name can`t be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}
