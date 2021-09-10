package domain

import "errors"

var (
	ErrInvalidToken = errors.New("invalid jwt token")
)
