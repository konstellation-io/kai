package version

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
	"github.com/spf13/viper"
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
		h.registerActionFailed(user.ID, productID, v, CommentUserNotAuthorized, "start")

		return nil, nil, err
	}

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.ID, productID, v, CommentVersionNotFound, "start")

		return nil, nil, err
	}

	if !vers.CanBeStarted() {
		h.registerActionFailed(user.ID, productID, vers, CommentInvalidVersionStatus, "start")
		return nil, nil, internalerrors.ErrInvalidVersionStatusBeforeStarting
	}

	versionCfg, err := h.getVersionConfig(ctx, productID, vers)
	if err != nil {
		h.registerActionFailed(user.ID, productID, vers, CommentErrorCreatingNATSResources, "start")
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

	vers.Tag = strings.ReplaceAll(vers.Tag, ".", "-")

	err := h.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		h.registerActionFailed(userID, productID, vers, CommentErrorStartingVersion, "start")
		h.handleVersionServiceActionError(ctx, productID, vers, notifyStatusCh, err)

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
