package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUserAgent(
	ctx context.Context,
	id int,
) error {

	err := s.usersRepository.DeleteUserAgent(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

