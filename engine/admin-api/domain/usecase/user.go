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

func (ui *UserInteractor) UpdateUserProductGrants(
	ctx context.Context,
	user *entity.User,
	targetUserID,
	product string,
	grants []auth.Action,
	comment ...string,
) error {
	if err := ui.accessControl.CheckRoleGrants(user, auth.ActUpdateUserGrants); err != nil {
		return fmt.Errorf("checking role grants: %w", err)
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserID, product, grants)
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductGrants(
		user.ID,
		targetUserID,
		product,
		grants,
		givenComment,
	)
	if err != nil {
		return fmt.Errorf("registering user activity: %w", err)
	}

	ui.logger.Info("Updated user grants for product", "user", targetUserID, "product", product, "grants", grants)

	return nil
}

func (ui *UserInteractor) RevokeUserProductGrants(
	ctx context.Context,
	user *entity.User,
	targetUserID,
	product string,
	comment ...string,
) error {
	if err := ui.accessControl.CheckRoleGrants(user, auth.ActUpdateUserGrants); err != nil {
		return fmt.Errorf("checking role grants: %w", err)
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserID, product, []auth.Action{})
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductGrants(
		user.ID,
		targetUserID,
		product,
		[]auth.Action{},
		givenComment,
	)
	if err != nil {
		return fmt.Errorf("registering user activity: %w", err)
	}

	ui.logger.Info("Revoked user grants for product", "user", targetUserID, "product", product)

	return nil
}

func (ui *UserInteractor) AddUserToProduct(
	ctx context.Context,
	user *entity.User,
	product,
	targetUserEmail string,
) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUser); err != nil {
		return err
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserEmail, product, auth.GetProductUserGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) RemoveUserFromProduct(
	ctx context.Context,
	user *entity.User,
	product,
	targetUserEmail string,
) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUser); err != nil {
		return err
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserEmail, product, []auth.Action{})
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) AddMaintainerToProduct(
	ctx context.Context,
	user *entity.User,
	product,
	targetUserEmail string,
) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUser); err != nil {
		return err
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserEmail, product, auth.GetProductMantainerGrants())
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("Maintainer added to product", "user", targetUserEmail, "product", product)

	return nil
}

func (ui *UserInteractor) RemoveMaintainerFromProduct(
	ctx context.Context,
	user *entity.User,
	product,
	targetUserEmail string,
) error {
	if err := ui.accessControl.CheckProductGrants(user, product, auth.ActManageProductUser); err != nil {
		return err
	}

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserEmail, product, []auth.Action{})
	if err != nil {
		return fmt.Errorf("updating grants in user's registry: %w", err)
	}

	ui.logger.Info("User deleted from product", "user", targetUserEmail, "product", product)

	return nil
}
