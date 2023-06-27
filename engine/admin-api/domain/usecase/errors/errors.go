package errors

import (
	"errors"
	"fmt"
)

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

var (
	// ErrVersionNotFound error.
	ErrVersionNotFound = errors.New("error version not found")

	// ErrVersionDuplicated error.
	ErrVersionDuplicated = errors.New("error version duplicated")

	// ErrVersionConfigIncomplete error.
	ErrVersionConfigIncomplete = errors.New("version config is incomplete")

	// ErrVersionConfigInvalidKey error.
	ErrVersionConfigInvalidKey = errors.New("version config contains an unknown key")

	// ErrUpdatingStartedVersionConfig error.
	ErrUpdatingStartedVersionConfig = errors.New("config can't be incomplete for started version")

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

type ErrInvalidKRT struct {
	msg  string
	errs error
}

func (e ErrInvalidKRT) Error() string {
	return fmt.Sprintf("%s:\n%s", e.msg, e.errs)
}

func (e ErrInvalidKRT) GetErrors() error {
	return e.errs
}

func NewErrInvalidKRT(msg string, errs error) ErrInvalidKRT {
	return ErrInvalidKRT{
		msg:  msg,
		errs: errs,
	}
}
