package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

var (
	ErrParsingKRTFile         = errors.New("error parsing KRT file")
	ErrVersionNotFound        = errors.New("error version not found")
	ErrVersionDuplicated      = errors.New("error version duplicated")
	ErrUserNotAuthorized      = errors.New("error user not authorized")
	ErrVersionCannotBeStarted = errors.New("error version cannot be started, status must be 'created', 'stopped' or 'failed'")
	ErrCreatingNATSResources  = errors.New("error creating NATS resources")
	ErrDeletingNATSResources  = errors.New("error deleting NATS resources")
	ErrStartingVersion        = errors.New("error starting version")
	ErrorStoppingVersion      = errors.New("error stopping version")
)

var (
	ErrUpdatingVersionStatus   = errors.New("error updating version status")
	ErrUpdatingVersionError    = errors.New("error updating version error")
	ErrRegisteringUserActivity = errors.New("error registering user activity")
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

func (h *Handler) registerActionFailed(userID, productID string, vers *entity.Version, comment, action string) {
	var err error
	if action == "start" {
		err = h.userActivityInteractor.RegisterStartAction(userID, productID, vers, comment)
	} else if action == "stop" {
		err = h.userActivityInteractor.RegisterStopAction(userID, productID, vers, comment)
	}

	if err != nil {
		h.logger.Error(ErrRegisteringUserActivity, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}
}

func (h *Handler) handleVersionServiceActionError(
	ctx context.Context, productID string, vers *entity.Version,
	notifyStatusCh chan *entity.Version, actionErr error,
) {
	_, err := h.versionRepo.SetError(ctx, productID, vers, actionErr.Error())
	if err != nil {
		h.logger.Error(ErrUpdatingVersionError, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"versionError", actionErr.Error(),
		)
	}

	vers.Status = entity.VersionStatusError
	vers.Error = actionErr.Error()
	notifyStatusCh <- vers
}
