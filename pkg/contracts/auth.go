package contracts

import "errors"

type (
	Validate interface {
		Validate() error
	}
)

type (
	SignUpRequest struct {
		Email    string
		Password string
		Role     int
	}

	SignUpResponse struct {
		ProfileID string
	}
)

func (sr *SignUpRequest) Validate() error {
	return errors.New("")
}
