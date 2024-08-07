package errors

import "errors"

var (
	ErrURLNotFound = errors.New("URL not found")
	ErrRedisError  = errors.New("an error occurred while interacting with Redis")
)
