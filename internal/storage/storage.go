package storage

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInputTooLong = errors.New("input value(s) is(are) too long")
)
