package users_service

import (
	"context"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUserAgent(
		ctx context.Context,
		user domain.UserAgent,
	) (domain.UserAgent, error)
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
