package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AMmetro/ODRProcesing/internal/core/domain"
	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUserAgent(
	ctx context.Context,
	id int,
) (domain.UserAgent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    SELECT id, version, full_name, phone_number FROM ODRProcesing.users WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserAgentModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserAgent{}, fmt.Errorf("user not found:: %v: %w", err, core_errors.ErrNotFound)
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
