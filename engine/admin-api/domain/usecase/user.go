package usecase

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

const (
	getUserByIDWrapper             = "get user by id"
	updateUserProductGrantsWrapper = "update user product grants"
	revokeUserProductGrantsWrapper = "revoke user product grants"
	updateUserProductGrantsLog     = "Updated user %q grants for product %q: %v"
	revokeUserProductGrantsLog     = "Revoked user %q grants for product %q"
)

type UserInteractor struct {
	logger                 logging.Logger
	userActivityInteractor UserActivityInteracter
	userRegistry           service.UserRegistry
}

// NewUserInteractor creates a new UserInteractor.
//
// UserInteractor is the usecase to manage users.
func NewUserInteractor(
	logger logging.Logger,
	userActivityInteractor UserActivityInteracter,
	gocloakManager service.UserRegistry,
) *UserInteractor {
	return &UserInteractor{
		logger,
		userActivityInteractor,
		gocloakManager,
	}
}

func (ui *UserInteractor) GetUserByID(userID string) (*entity.User, error) {
	wrapErr := errors.Wrapper(getUserByIDWrapper + ": %w")

	user, err := ui.userRegistry.GetUserByID(userID)
	if err != nil {
		return nil, wrapErr(err)
	}

	return user, nil
}

func (ui *UserInteractor) UpdateUserProductGrants(
	triggerUserID,
	targetUserID,
	product string,
	grants []string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(updateUserProductGrantsWrapper + ": %w")

	err := ui.userRegistry.UpdateUserProductGrants(targetUserID, product, grants)
	if err != nil {
		return wrapErr(err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductGrants(
		triggerUserID,
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
	triggerUserID,
	targetUserID,
	product string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(revokeUserProductGrantsWrapper + ": %w")

	err := ui.userRegistry.UpdateUserProductGrants(targetUserID, product, []string{})
	if err != nil {
		return wrapErr(err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductGrants(
		triggerUserID,
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
