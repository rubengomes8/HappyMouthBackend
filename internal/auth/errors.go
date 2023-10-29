package auth

import "errors"

var (
	ErrWeakPassword = errors.New("auth.error.weak_password")
)
