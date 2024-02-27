package usecase

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"golang.org/x/net/context"
)

type UserInteractor struct {
	logger                 logr.Logger
	accessControl          auth.AccessControl
	userActivityInteractor UserActivityInteracter
	userRegistry           service.UserRegistry
}

// NewUserInteractor creates a new UserInteractor.
//
// UserInteractor is the usecase to manage users.
func NewUserInteractor(
	logger logr.Logger,
	accessControl auth.AccessControl,
	userActivityInteractor UserActivityInteracter,
	userRegistry service.UserRegistry,
) *UserInteractor {
	return &UserInteractor{
		logger,
		accessControl,
		userActivityInteractor,
		userRegistry,
	}
}

func (ui *UserInteractor) AddUserToProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.AddProductGrants(ctx, targetUserEmail, product, auth.GetProductUserGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) RemoveUserFromProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.RevokeProductGrants(ctx, targetUserEmail, product, auth.GetProductUserGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) AddMaintainerToProduct(ctx context.Context, user *entity.User, targetUserEmail, product string) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.AddProductGrants(ctx, targetUserEmail, product, auth.GetProductMaintainerGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("Maintainer added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) RemoveMaintainerFromProduct(
	ctx context.Context,
	user *entity.User,
	targetUserEmail,
	product string,
) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUsers); err != nil {
		return err
	}

	err := ui.userRegistry.RevokeProductGrants(ctx, targetUserEmail, product, auth.GetProductMaintainerGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}
