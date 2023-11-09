package auth

import "errors"

var (
	ErrWeakPassword          = errors.New("auth.error.weak_password")
	ErrUsernameAlreadyExists = errors.New("auth.error.username_already_exists")
)
