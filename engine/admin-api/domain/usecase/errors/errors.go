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

	// ErrInvalidKRT error.
	ErrInvalidKRT = errors.New("invalid KRT")

	// ErrStoringKRTFile error.
	ErrStoringKRTFile = errors.New("error storing KRT file")
)

func ParsingKRTFileError(err error) error {
	return fmt.Errorf("%w: %w", ErrParsingKRTFile, err)
}

func InvalidKRTError(validationErrors error) error {
	return fmt.Errorf("%w:\n%w", ErrInvalidKRT, validationErrors)
}
