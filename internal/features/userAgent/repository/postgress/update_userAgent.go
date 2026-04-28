package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
	core_postgres_pool "github.com/AMmetro/ODRProcesing/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) UpdateUserAgent(
	ctx context.Context,
	user domain.UserAgent,
) (domain.UserAgent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    UPDATE ODRProcesing.users
	SET full_name=$1, phone_number=$2,
	version = version + 1
	WHERE id = $3 AND version = $4
    RETURNING id, version, full_name, phone_number
    `
	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, user.ID, user.Version)

	var userModel UserAgentModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserAgent{}, fmt.Errorf("user with id='%d' concurently access: %w",
				user.ID, core_errors.ErrConflict)
		}
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
