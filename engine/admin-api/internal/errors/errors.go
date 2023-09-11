package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrVersionNotFound error.
	ErrVersionNotFound = errors.New("error version not found")

	// ErrVersionDuplicated error.
	ErrVersionDuplicated = errors.New("error version duplicated")

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

	// ErrParsingKRTFile error.
	ErrParsingKRTFile = errors.New("error parsing KRT file")

	// ErrStoringKRTFile error.
	ErrStoringKRTFile = errors.New("error storing KRT file")
)

func ParsingKRTFileError(err error) error {
	return fmt.Errorf("%w: %w", ErrParsingKRTFile, err)
}

type KRTValidationError struct {
	msg  string
	errs error
}

func (e KRTValidationError) Error() string {
	return fmt.Sprintf("%s:\n%s", e.msg, e.errs)
}

func (e KRTValidationError) GetErrors() error {
	return e.errs
}

func NewErrInvalidKRT(msg string, errs error) KRTValidationError {
	return KRTValidationError{
		msg:  msg,
		errs: errs,
	}
}
