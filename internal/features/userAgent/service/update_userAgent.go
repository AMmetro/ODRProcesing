package users_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

func (s *UsersService) UpdateUserAgent(
	ctx context.Context,
	userId int,
	patch domain.UserAgentPatch,
) (domain.UserAgent, error) {
	existingUser, err := s.usersRepository.GetUserAgent(ctx, userId)
	if err != nil {
		return domain.UserAgent{}, fmt.Errorf("get user for update: %w", err)
	}

	if err := existingUser.ApplyPatch(patch); err != nil {
		return domain.UserAgent{}, fmt.Errorf("applay patch user: %w", err)
	}

	updatedUser, err := s.usersRepository.UpdateUserAgent(ctx, existingUser)
	if err != nil {
		return domain.UserAgent{}, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}
