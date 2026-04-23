package users_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.UserAgent, error) {
	var users []domain.UserAgent
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non negative %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non negative %w", core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return []domain.UserAgent{}, fmt.Errorf("get users from repository: %w", err)
	}
	return users, nil
}
