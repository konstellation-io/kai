package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
	"github.com/spf13/viper"
)

const (
	CommentUserNotAuthorized          = "user not authorized"
	CommentVersionNotFound            = "version not found"
	CommentInvalidVersionStatus       = "invalid version status before starting"
	CommentErrorCreatingNATSResources = "error creating NATS resources"
	CommentErrorStartingVersion       = "error starting version"
)

var (
	ErrUpdatingVersionStatus   = fmt.Errorf("error updating version status")
	ErrUpdatingVersionErrors   = fmt.Errorf("error updating version errors")
	ErrRegisteringUserActivity = fmt.Errorf("error registering user activity")
)

// Start a previously created Version.
func (h *Handler) Start(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionTag string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	h.logger.Info("Starting version", "userID", user.ID, "versionTag", versionTag, "productID", productID)

	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerStartActionFailed(user.ID, productID, v, CommentUserNotAuthorized)
		return nil, nil, err
	}

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerStartActionFailed(user.ID, productID, v, CommentVersionNotFound)
		return nil, nil, err
	}

	if !vers.CanBeStarted() {
		h.registerStartActionFailed(user.ID, productID, vers, CommentInvalidVersionStatus)
		return nil, nil, internalerrors.ErrInvalidVersionStatusBeforeStarting
	}

	versionCfg, err := h.getVersionConfig(ctx, productID, vers)
	if err != nil {
		h.registerStartActionFailed(user.ID, productID, vers, CommentErrorCreatingNATSResources)
		return nil, nil, err
	}

	vers.Status = entity.VersionStatusStarting

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarting,
		)
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	go h.startAndNotify(user.ID, productID, comment, vers, versionCfg, notifyStatusCh)

	return vers, notifyStatusCh, nil
}

func (h *Handler) registerStartActionFailed(userID, productID string, vers *entity.Version, comment string) {
	err := h.userActivityInteractor.RegisterStartAction(userID, productID, vers, comment)
	if err != nil {
		h.logger.Error(ErrRegisteringUserActivity, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}
}

func (h *Handler) getVersionConfig(ctx context.Context, productID string, vers *entity.Version) (*entity.VersionConfig, error) {
	versionStreamCfg, err := h.natsManagerService.CreateStreams(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating streams for version %q: %w", vers.Tag, err)
	}

	objectStoreCfg, err := h.natsManagerService.CreateObjectStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating objects stores for version %q: %w", vers.Tag, err)
	}

	kvStoreCfg, err := h.natsManagerService.CreateKeyValueStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating key-value stores for version %q: %w", vers.Tag, err)
	}

	versionCfg := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)

	return versionCfg, nil
}

func (h *Handler) startAndNotify(
	userID string,
	productID string,
	comment string,
	vers *entity.Version,
	versionConfig *entity.VersionConfig,
	notifyStatusCh chan *entity.Version,
) {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
	}()

	err := h.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		h.registerStartActionFailed(userID, productID, vers, CommentErrorStartingVersion)
		h.handleVersionServiceStartError(ctx, productID, vers, notifyStatusCh, err)
		return
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarted)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarted,
		)
	}

	err = h.userActivityInteractor.RegisterStartAction(userID, productID, vers, comment)
	if err != nil {
		h.logger.Error(ErrRegisteringUserActivity, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
}

func (h *Handler) handleVersionServiceStartError(
	ctx context.Context, productID string, vers *entity.Version,
	notifyStatusCh chan *entity.Version, startErr error,
) {
	err := h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusError)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusError,
		)
	}

	_, err = h.versionRepo.SetError(ctx, productID, vers, startErr.Error())
	if err != nil {
		h.logger.Error(ErrUpdatingVersionErrors, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusError,
		)
	}

	vers.Status = entity.VersionStatusError
	vers.Error = startErr.Error()
	notifyStatusCh <- vers
}
