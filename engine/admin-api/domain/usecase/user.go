package usecase

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"golang.org/x/net/context"
)

type UserInteractor struct {
	logger                 logging.Logger
	accessControl          auth.AccessControl
	userActivityInteractor UserActivityInteracter
	userRegistry           service.UserRegistry
}

// NewUserInteractor creates a new UserInteractor.
//
// UserInteractor is the usecase to manage users.
func NewUserInteractor(
	logger logging.Logger,
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
	grants []string,
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

	ui.logger.Infof("Updated user %q grants for product %q: %v", targetUserID, product, grants)

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

	err := ui.userRegistry.UpdateUserProductGrants(ctx, targetUserID, product, []string{})
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
		[]string{},
		givenComment,
	)
	if err != nil {
		return fmt.Errorf("registering user activity: %w", err)
	}

	ui.logger.Infof("Revoked user %q grants for product %q", targetUserID, product)

	return nil
}
