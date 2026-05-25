package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/AMmetro/ODRProcesing/shared/pkg/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type valitable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	// parse JSON from request to dest
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error
	// check if dest implements interface valitable
	v, ok := dest.(valitable)
	if ok {
		// if implements, validate with implemented methods
		err = v.Validate()
	} else {
		// validate structure dest {json:"full_name" validate:"required,min=3,max=100"}
		// according to tags in structure call methods Validate
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}

