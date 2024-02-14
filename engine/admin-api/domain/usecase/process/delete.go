package process

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

type DeleteProcessOpts struct {
	Product  string
	Version  string
	Process  string
	IsPublic bool
}

func (o DeleteProcessOpts) Validate() error {
	if o.Product == "" && !o.IsPublic {
		return ErrMissingProductInParams
	}

	if o.Product != "" && o.IsPublic {
		return ErrIsPublicAndHasProduct
	}

	if o.Version == "" {
		return ErrMissingVersionInParams
	}

	if o.Process == "" {
		return ErrMissingProcessInParams
	}

	return nil
}

func (ps *Handler) DeleteProcess(
	ctx context.Context,
	user *entity.User,
	opts DeleteProcessOpts,
) (string, error) {
	ps.logger.Info("Deleting process", "Product", opts.Product, "Version", opts.Version, "Process", opts.Process, "IsPublic", opts.IsPublic)

	if err := opts.Validate(); err != nil {
		return "", err
	}

	if err := ps.checkDeleteGrants(user, opts); err != nil {
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

func (ps *Handler) checkDeleteGrants(user *entity.User, opts DeleteProcessOpts) error {
	if opts.IsPublic {
		return ps.accessControl.CheckRoleGrants(user, auth.ActDeletePublicProcess)
	}

	return ps.accessControl.CheckProductGrants(user, opts.Product, auth.ActDeleteProcess)
}
