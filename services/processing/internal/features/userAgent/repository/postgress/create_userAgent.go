package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (r *UsersRepository) CreateUserAgent(
	ctx context.Context,
	user domain.UserAgent,
) (domain.UserAgent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    INSERT INTO ODRProcesing.users (full_name, phone_number)
    VALUES ($1, $2)
    RETURNING id, version, full_name, phone_number;
    `
	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserAgentModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.UserAgent{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUserAgent(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil

}

