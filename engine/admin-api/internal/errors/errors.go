package errors

import (
	"errors"
)

var (
	// ErrInvalidVersionStatusBeforeStarting error.
	ErrInvalidVersionStatusBeforeStarting = errors.New("the version must be stopped before starting")

	// ErrInvalidVersionStatusBeforeStopping error.
	ErrInvalidVersionStatusBeforeStopping = errors.New("the version must be started before stopping")

	// ErrInvalidVersionStatusBeforePublishing error.
	ErrInvalidVersionStatusBeforePublishing = errors.New("the version must be started before publishing")

	// ErrInvalidVersionStatusBeforeUnpublishing error.
	ErrInvalidVersionStatusBeforeUnpublishing = errors.New("the version must be published before unpublishing")

	// ErrCreatingDashboard error.
	ErrCreatingDashboard = errors.New("error creating dashboard")

	// ErrStoringKRTFile error.
	ErrStoringKRTFile = errors.New("error storing KRT file")
)
