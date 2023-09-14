package version

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/krt/pkg/parse"
)

// Create creates a Version on the DB based on the content of a KRT file.
func (h *Handler) Create(
	ctx context.Context,
	user *entity.User,
	productID string,
	krtFile io.Reader,
) (*entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActCreateVersion); err != nil {
		return nil, err
	}

	_, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error product repo GetById: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "version")
	if err != nil {
		return nil, fmt.Errorf("error creating temp dir for version: %w", err)
	}

	h.logger.Info("Created temp dir to extract the KRT files at " + tmpDir)

	tmpKrtFile, err := h.copyStreamToTempFile(krtFile)
	if err != nil {
		return nil, fmt.Errorf("error creating temp krt file for version: %w", err)
	}

	krtYml, err := parse.ParseFile(tmpKrtFile.Name())
	if err != nil {
		return nil, ParsingKRTFileError(err)
	}

	err = krtYml.Validate()
	if err != nil {
		return nil, NewErrInvalidKRT(
			"invalid KRT file",
			err,
		)
	}

	_, err = h.versionRepo.GetByTag(ctx, productID, krtYml.Version)
	if err != nil && !errors.Is(err, ErrVersionNotFound) {
		return nil, fmt.Errorf("error version repo GetByTag: %w", err)
	} else if err == nil {
		return nil, ErrVersionDuplicated
	}

	versionCreated, err := h.versionRepo.Create(
		user.Email,
		productID,
		h.mapKrtToVersion(krtYml),
	)
	if err != nil {
		return nil, err
	}

	err = h.userActivityInteractor.RegisterCreateAction(user.Email, productID, versionCreated)
	if err != nil {
		return nil, fmt.Errorf("registering create version action: %w", err)
	}

	h.logger.Info("Version created")

	return versionCreated, nil
}

func (h *Handler) setStatusError(
	ctx context.Context,
	productID string,
	vers *entity.Version,
	errs []error,
	notifyCh chan *entity.Version,
) {
	errorMessages := make([]string, len(errs))

	var jointErrors error

	for idx, err := range errs {
		errorMessages[idx] = err.Error()
		jointErrors = errors.Join(jointErrors, err)
	}

	h.logger.Error(jointErrors, "Errors found in version", "version tag", vers.Tag)

	versionWithError, err := h.versionRepo.SetErrors(ctx, productID, vers, errorMessages)
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
