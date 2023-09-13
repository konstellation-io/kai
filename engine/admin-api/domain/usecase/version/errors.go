package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const (
	CommentUserNotAuthorized                  = "user not authorized"
	CommentVersionNotFound                    = "version not found"
	CommentInvalidVersionStatusBeforeStarting = "invalid version status before starting"
	CommentInvalidVersionStatusBeforeStopping = "invalid version status before stopping"
	CommentErrorCreatingNATSResources         = "error creating NATS resources"
	CommentErrorDeletingNATSResources         = "error deleting NATS resources"
	CommentErrorStartingVersion               = "error starting version"
	CommentErrorStoppingVersion               = "error stopping version"
)

var (
	ErrUpdatingVersionStatus   = fmt.Errorf("error updating version status")
	ErrUpdatingVersionError    = fmt.Errorf("error updating version error")
	ErrRegisteringUserActivity = fmt.Errorf("error registering user activity")
)

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
