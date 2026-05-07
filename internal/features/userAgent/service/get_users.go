package users_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	domain_utils "github.com/AMmetro/ODRProcesing/internal/core/domain/utils"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.UserAgent, error) {
	var users []domain.UserAgent

	err := domain_utils.LimitOffsetValidation(limit, offset)
	if err != nil {
		return []domain.UserAgent{}, fmt.Errorf("validate limit and offset: %w", err)
	}

	users, err = s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return []domain.UserAgent{}, fmt.Errorf("get users from repository: %w", err)
	}
	return users, nil
}
