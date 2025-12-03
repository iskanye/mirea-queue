package repositories

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyInQueue = errors.New("user already in queue")
)
