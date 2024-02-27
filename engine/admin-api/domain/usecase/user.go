package usecase

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"golang.org/x/net/context"
)

type UserHandler struct {
	logger                 logr.Logger
	accessControl          auth.AccessControl
	userActivityInteractor UserActivityInteracter
	userRegistry           service.UserRegistry
}

// NewUserHandler creates a new UserHandler.
//
// UserHandler is the usecase to manage users.
func NewUserHandler(
	logger logr.Logger,
	accessControl auth.AccessControl,
	userActivityInteractor UserActivityInteracter,
	userRegistry service.UserRegistry,
) *UserHandler {
	return &UserHandler{
		logger,
		accessControl,
		userActivityInteractor,
		userRegistry,
	}
}

func (ui *UserHandler) AddUserToProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.AddProductGrants(ctx, targetUserEmail, product, auth.GetDefaultUserGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserHandler) RemoveUserFromProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.RevokeProductGrants(ctx, targetUserEmail, product, auth.GetDefaultUserGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserHandler) AddMaintainerToProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckRoleGrants(user, auth.ActManageProductMaintainers); err != nil {
		return err
	}

	err := ui.userRegistry.AddProductGrants(ctx, targetUserEmail, product, auth.GetDefaultMaintainerGrants())
	if err != nil {
		return fmt.Errorf("adding product grants in user's registry: %w", err)
	}

	ui.logger.Info("Maintainer added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserHandler) RemoveMaintainerFromProduct(
	ctx context.Context,
	user *entity.User,
	targetUserEmail,
	product string,
) error {
	if err := ui.accessControl.CheckRoleGrants(user, auth.ActManageProductMaintainers); err != nil {
		return err
	}

	err := ui.userRegistry.RevokeProductGrants(ctx, targetUserEmail, product, auth.GetDefaultMaintainerGrants())
	if err != nil {
		return fmt.Errorf("revoking product grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}
