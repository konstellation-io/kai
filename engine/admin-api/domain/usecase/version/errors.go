package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const (
	StartAction = iota
	StopAction
)

var (
	ErrParsingKRTFile             = errors.New("error parsing KRT file")
	ErrVersionNotFound            = errors.New("error version not found")
	ErrVersionDuplicated          = errors.New("error version duplicated")
	ErrUserNotAuthorized          = errors.New("error user not authorized")
	ErrVersionCannotBeStarted     = errors.New("error version cannot be started, status must be 'created', 'stopped' or 'failed'")
	ErrVersionCannotBeStopped     = errors.New("error version cannot be stopped, status must be 'started'")
	ErrVersionCannotBePublished   = errors.New("error publishing version, status must be 'started'")
	ErrVersionCannotBeUnpublished = errors.New("error unpublishing version, status must be 'published'")
	ErrCreatingNATSResources      = errors.New("error creating NATS resources")
	ErrDeletingNATSResources      = errors.New("error deleting NATS resources")
	ErrStartingVersion            = errors.New("error starting version")
	ErrStoppingVersion            = errors.New("error stopping version")
	ErrUnpublishingVersion        = errors.New("error unpublishing version")
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

func (h *Handler) registerStopActionFailed(userEmail, productID string, vers *entity.Version, incomingErr error) {
	err := h.userActivityInteractor.RegisterStopAction(userEmail, productID, vers, incomingErr.Error())
	if err != nil {
		h.logger.Error(err, "Error registering user activity",
			"productID", productID,
			"versionTag", vers.Tag,
			"error", incomingErr.Error(),
		)
	}
}

func (h *Handler) handleVersionServiceActionError(
	ctx context.Context, productID string, vers *entity.Version,
	notifyStatusCh chan *entity.Version, actionErr error,
) {
	err := h.versionRepo.SetError(ctx, productID, vers, actionErr.Error())
	if err != nil {
		h.logger.Error(err, "Error updating version error",
			"productID", productID,
			"versionTag", vers.Tag,
			"versionError", actionErr.Error(),
		)
	}

	vers.Status = entity.VersionStatusError
	vers.Error = actionErr.Error()
	notifyStatusCh <- vers
}
