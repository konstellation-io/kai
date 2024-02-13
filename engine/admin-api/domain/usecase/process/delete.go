package process

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// TODO
// Current config doesnt work, admin api is not able to get up
// Talk about usage of basic user pwd auth rather than token use
// Ensure config for both registry and admin-api are set up properly
// Add a test for this function
// Try it out in a local environment

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

	if err := ps.processRegistry.DeleteProcess(imageName, opts.Version); err != nil {
		return "", err
	}

	if err := ps.processRepository.Delete(ctx, scope, processID); err != nil {
		return "", err
	}

	return processID, nil
}
