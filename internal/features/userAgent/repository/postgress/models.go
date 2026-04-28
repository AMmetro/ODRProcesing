package users_postgres_repository

import "github.com/AMmetro/ODRProcesing/internal/core/domain"

type UserAgentModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainFromModels(users []UserAgentModel) []domain.UserAgent {
	res := make([]domain.UserAgent, len(users))
	for i, user := range users {
		res[i] = domain.NewUserAgent(user.ID, user.Version, user.FullName, user.PhoneNumber)
	}
	return res
}
