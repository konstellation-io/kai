package process

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func (ps *Handler) DeleteProcess(
	ctx context.Context,
	user *entity.User,
	opts DeleteProcessOpts,
) (string, error) {
	ps.logger.Info("Deleting process", "Product", opts.Product, "Version", opts.Version, "Process", opts.Process, "IsPublic", opts.IsPublic)

	if err := opts.Validate(); err != nil {
		return "", err
	}

	if err := ps.checkDeleteGrants(user, opts.IsPublic, opts.Product); err != nil {
		return "", err
	}

	scope := ps.getProcessRegisterScope(opts.IsPublic, opts.Product)
	processID := ps.getProcessID(scope, opts.Process, opts.Version)
	imageName := ps.getImageName(scope, opts.Process)

	_, err := ps.processRepository.GetByID(ctx, scope, processID)
	if err != nil {
		return "", err
	}

	if err := ps.processRegistry.DeleteProcess(ctx, imageName, opts.Version); err != nil {
		return "", err
	}

	if err := ps.processRepository.Delete(ctx, scope, processID); err != nil {
		return "", err
	}

	return processID, nil
}
