package services

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyInQueue = errors.New("user already in queue")
	ErrQueueEnd       = errors.New("entry is at the end of the queue")
)
