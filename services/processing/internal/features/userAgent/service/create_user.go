package users_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (s *UsersService) CreateUserAgent(
	ctx context.Context,
	userAgent domain.UserAgent,
) (domain.UserAgent, error) {
	if err := userAgent.Validate(); err != nil {
		return domain.UserAgent{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreateUserAgent(ctx, userAgent)
	if err != nil {
		return domain.UserAgent{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

