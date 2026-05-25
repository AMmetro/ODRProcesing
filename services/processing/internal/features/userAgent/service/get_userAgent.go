package users_service

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (s *UsersService) GetUserAgent(
	ctx context.Context,
	id int,
) (domain.UserAgent, error) {
	users, err := s.usersRepository.GetUserAgent(ctx, id)
	if err != nil {
		return domain.UserAgent{}, fmt.Errorf("get userAgent from repository: %w", err)
	}
	return users, nil
}

