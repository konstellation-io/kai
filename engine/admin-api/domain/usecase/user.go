package usecase

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

const (
	updateUserProductGrantsWrapper = "update user product grants"
	revokeUserProductGrantsWrapper = "revoke user product grants"
	updateUserProductGrantsLog     = "Updated user %q grants for product %q: %v"
	revokeUserProductGrantsLog     = "Revoked user %q grants for product %q"
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
	user *entity.User,
	targetUserID,
	product string,
	grants []string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(updateUserProductGrantsWrapper + ": %w")
	if err := ui.accessControl.CheckRoleGrants(user, auth.ActUpdateUserGrants); err != nil {
		return wrapErr(err)
	}

	err := ui.userRegistry.UpdateUserProductGrants(targetUserID, product, grants)
	if err != nil {
		return wrapErr(err)
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
		return wrapErr(err)
	}

	ui.logger.Infof(updateUserProductGrantsLog, targetUserID, product, grants)

	return nil
}

func (ui *UserInteractor) RevokeUserProductGrants(
	user *entity.User,
	targetUserID,
	product string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(revokeUserProductGrantsWrapper + ": %w")

	if err := ui.accessControl.CheckRoleGrants(user, auth.ActUpdateUserGrants); err != nil {
		return wrapErr(err)
	}

	err := ui.userRegistry.UpdateUserProductGrants(targetUserID, product, []string{})
	if err != nil {
		return wrapErr(err)
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
		return wrapErr(err)
	}

	ui.logger.Infof(revokeUserProductGrantsLog, targetUserID, product)

	return nil
}
