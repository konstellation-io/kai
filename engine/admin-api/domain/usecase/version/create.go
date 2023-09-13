package version

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/krt"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

// Create creates a Version on the DB based on the content of a KRT file.
func (h *Handler) Create(
	ctx context.Context,
	user *entity.User,
	productID string,
	krtFile io.Reader,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActCreateVersion); err != nil {
		return nil, nil, err
	}

	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("error product repo GetById: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "version")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating temp dir for version: %w", err)
	}

	h.logger.Info("Created temp dir to extract the KRT files at " + tmpDir)

	tmpKrtFile, err := h.copyStreamToTempFile(krtFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating temp krt file for version: %w", err)
	}

	krtYml, err := krt.ParseFile(tmpKrtFile.Name())
	if err != nil {
		return nil, nil, internalerrors.ParsingKRTFileError(err)
	}

	err = krtYml.Validate()
	if err != nil {
		return nil, nil, internalerrors.NewErrInvalidKRT(
			"invalid KRT file",
			err,
		)
	}

	_, err = h.versionRepo.GetByTag(ctx, productID, krtYml.Version)
	if err != nil && !errors.Is(err, internalerrors.ErrVersionNotFound) {
		return nil, nil, fmt.Errorf("error version repo GetByTag: %w", err)
	} else if err == nil {
		return nil, nil, internalerrors.ErrVersionDuplicated
	}

	versionCreated, err := h.versionRepo.Create(
		user.ID,
		productID,
		krt.MapKrtYamlToVersion(krtYml),
	)

	if err != nil {
		return nil, nil, err
	}

	h.logger.Info("Version created")

	notifyStatusCh := make(chan *entity.Version, 1)

	go h.completeVersionCreation(
		user.ID, tmpKrtFile, tmpDir, product, versionCreated, notifyStatusCh,
	)

	return versionCreated, notifyStatusCh, nil
}

func (h *Handler) completeVersionCreation(
	loggedUserID string,
	tmpKrtFile *os.File,
	tmpDir string,
	product *entity.Product,
	versionCreated *entity.Version,
	notifyStatusCh chan *entity.Version,
) {
	ctx := context.Background()

	defer close(notifyStatusCh)

	defer func() {
		err := tmpKrtFile.Close()
		if err != nil {
			h.logger.Error(err, "Error closing file")
			return
		}

		err = os.Remove(tmpKrtFile.Name())
		if err != nil {
			h.logger.Error(err, "Error removing file")
		}
	}()

	var contentErrors []error

	dashboardsFolder := path.Join(tmpDir, "metrics/dashboards")
	contentErrors = h.saveKRTDashboards(ctx, dashboardsFolder, product, versionCreated, contentErrors)

	err := h.versionRepo.UploadKRTYamlFile(product.ID, versionCreated, tmpKrtFile.Name())
	if err != nil {
		contentErrors = append(contentErrors, internalerrors.ErrStoringKRTFile)
	}

	if len(contentErrors) > 0 {
		h.setStatusError(ctx, product.ID, versionCreated, contentErrors[0], notifyStatusCh)
		return
	}

	err = h.versionRepo.SetStatus(ctx, product.ID, versionCreated.ID, entity.VersionStatusCreated)
	if err != nil {
		versionCreated.Status = entity.VersionStatusError

		h.logger.Error(err, "Error setting version status")
		notifyStatusCh <- versionCreated

		return
	}

	// Notify state
	versionCreated.Status = entity.VersionStatusCreated
	notifyStatusCh <- versionCreated

	err = h.userActivityInteractor.RegisterCreateAction(loggedUserID, product.ID, versionCreated)
	if err != nil {
		h.logger.Error(err, "Error registering activity")
	}
}

// TODO discuss what will happen with this.
//
//nolint:godox // To be done.
func (h *Handler) saveKRTDashboards(
	ctx context.Context,
	dashboardsFolder string,
	product *entity.Product,
	versionCreated *entity.Version,
	contentErrors []error,
) []error {
	if _, err := os.Stat(path.Join(dashboardsFolder)); err == nil {
		err := h.storeDashboards(ctx, dashboardsFolder, product.ID, versionCreated.Tag)
		if err != nil {
			contentErrors = append(contentErrors, internalerrors.ErrCreatingDashboard)
		}
	}

	return contentErrors
}

func (h *Handler) setStatusError(
	ctx context.Context,
	productID string,
	vers *entity.Version,
	err error,
	notifyCh chan *entity.Version,
) {
	h.logger.Error(err, "Error found in version", "version tag", vers.Tag)

	versionWithError, err := h.versionRepo.SetError(ctx, productID, vers, err.Error())
	if err != nil {
		h.logger.Error(err, "Error saving version error state")
	}

	notifyCh <- versionWithError
}

func (h *Handler) copyStreamToTempFile(krtFile io.Reader) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "version")

	if err != nil {
		return nil, fmt.Errorf("error creating temp file for version: %w", err)
	}

	_, err = io.Copy(tmpFile, krtFile)
	if err != nil {
		return nil, fmt.Errorf("error copying temp file for version: %w", err)
	}

	h.logger.Info("Created temp file", "file name", tmpFile.Name())

	return tmpFile, nil
}
