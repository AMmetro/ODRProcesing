package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AMmetro/ODRProcesing/shared/pkg/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.UserAgent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `SELECT id, full_name, phone_number, version 
	FROM ODRProcesing.users ORDER BY id ASC
	LIMIT $1 OFFSET $2;`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var usersModels []UserAgentModel
	for rows.Next() {
		var usersModel UserAgentModel
		err := rows.Scan(
			&usersModel.ID,
			&usersModel.FullName,
			&usersModel.PhoneNumber,
			&usersModel.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		usersModels = append(usersModels, usersModel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainFromModels(usersModels)

	return userDomains, nil

}

