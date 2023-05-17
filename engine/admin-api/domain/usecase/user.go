package usecase

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

const (
	getUserByIDWrapper                  = "get user by id"
	updateUserProductPermissionsWrapper = "update user product permissions"
	revokeUserProductPermissionsWrapper = "revoke user product permissions"
	updateUserProductPermissionsLog     = "Updated user %q permissions for product %q: %v"
	revokeUserProductPermissionsLog     = "Revoked user %q permissions for product %q"
)

type UserInteractor struct {
	logger                 logging.Logger
	userActivityInteractor UserActivityInteracter
	gocloakManager         service.GocloakService
}

// NewUserInteractor creates a new UserInteractor.
//
// UserInteractor is the usecase to manage users.
func NewUserInteractor(
	logger logging.Logger,
	userActivityInteractor UserActivityInteracter,
	gocloakManager service.GocloakService,
) *UserInteractor {
	return &UserInteractor{
		logger,
		userActivityInteractor,
		gocloakManager,
	}
}

func (ui *UserInteractor) GetUserByID(userID string) (entity.UserGocloakData, error) {
	wrapErr := errors.Wrapper(getUserByIDWrapper + ": %w")

	user, err := ui.gocloakManager.GetUserByID(userID)
	if err != nil {
		return entity.UserGocloakData{}, wrapErr(err)
	}

	return user, nil
}

func (ui *UserInteractor) UpdateUserProductPermissions(
	triggerUserID,
	targetUserID,
	product string,
	permissions []string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(updateUserProductPermissionsWrapper + ": %w")

	err := ui.gocloakManager.UpdateUserProductPermissions(targetUserID, product, permissions)
	if err != nil {
		return wrapErr(err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductPermissions(
		triggerUserID,
		targetUserID,
		product,
		permissions,
		givenComment,
	)
	if err != nil {
		return wrapErr(err)
	}

	ui.logger.Infof(updateUserProductPermissionsLog, targetUserID, product, permissions)

	return nil
}

func (ui *UserInteractor) RevokeUserProductPermissions(triggerUserID,
	targetUserID,
	product string,
	comment ...string,
) error {
	wrapErr := errors.Wrapper(revokeUserProductPermissionsWrapper + ": %w")

	err := ui.gocloakManager.UpdateUserProductPermissions(targetUserID, product, []string{})
	if err != nil {
		return wrapErr(err)
	}

	var givenComment string
	if len(comment) > 0 {
		givenComment = comment[0]
	}

	err = ui.userActivityInteractor.RegisterUpdateProductPermissions(
		triggerUserID,
		targetUserID,
		product,
		[]string{},
		givenComment,
	)
	if err != nil {
		return wrapErr(err)
	}

	ui.logger.Infof(revokeUserProductPermissionsLog, targetUserID, product)

	return nil
}
