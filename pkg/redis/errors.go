package redis

import "errors"

var (
	ErrNotFound = errors.New("redis.not_found")
)
