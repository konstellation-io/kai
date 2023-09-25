package errors

import (
	"errors"
)

var (
	// ErrInvalidVersionStatusBeforeStopping error.
	ErrInvalidVersionStatusBeforeStopping = errors.New("the version must be started before stopping")

	// ErrInvalidVersionStatusBeforeUnpublishing error.
	ErrInvalidVersionStatusBeforeUnpublishing = errors.New("the version must be published before unpublishing")
)
