package lib

import "errors"

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrMissingToken      = errors.New("missing token")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrClaimsFound       = errors.New("no claims found in context")
)
