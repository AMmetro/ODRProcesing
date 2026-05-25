package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
)

func (r *UsersRepository) DeleteUserAgent(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
    DELETE FROM ODRProcesing.users WHERE id = $1 
    `
	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("error while deleted user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d':'%w' ", id, core_errors.ErrNotFound)
	}

	return nil

}

