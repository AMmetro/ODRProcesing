package domain_utils

import (
	"fmt"
	"net/http"

	core_errors "github.com/AMmetro/ODRProcesing/internal/core/errors"
	core_http_request "github.com/AMmetro/ODRProcesing/internal/core/transport/http/request"
)

func LimitOffsetValidation(limit *int, offset *int) error {
	if limit != nil && *limit < 0 {
		return fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func GetUserIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
		idQueryParamKey     = "userId"
	)

	userId, err := core_http_request.GetIntQueryParam(r, idQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userId' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, userId, nil
}
