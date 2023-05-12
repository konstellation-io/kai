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
	revokedPermissionsComment           = "revoked all permissions"
)

// UserInteractor contains app logic to handle User entities
type UserInteractor struct {
	logger                 logging.Logger
	userActivityInteractor UserActivityInteracter
	gocloakManager         service.GocloakService
}

// TODO: use user activity interactor, use logger
// NewUserInteractor creates a new UserInteractor
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

func (ui *UserInteractor) UpdateUserProductPermissions(triggerUserID string, targetUserID string, product string, permissions []string) error {
	wrapErr := errors.Wrapper(updateUserProductPermissionsWrapper + ": %w")

	err := ui.gocloakManager.UpdateUserProductPermissions(targetUserID, product, permissions)
	if err != nil {
		return wrapErr(err)
	}

	err = ui.userActivityInteractor.RegisterUpdateProductPermissions(triggerUserID, targetUserID, product, permissions, "")
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func (ui *UserInteractor) RevokeUserProductPermissions(triggerUserID string, targetUserID string, product string) error {
	wrapErr := errors.Wrapper(revokeUserProductPermissionsWrapper + ": %w")

	err := ui.gocloakManager.UpdateUserProductPermissions(targetUserID, product, []string{})
	if err != nil {
		return wrapErr(err)
	}

	err = ui.userActivityInteractor.RegisterUpdateProductPermissions(triggerUserID, targetUserID, product, []string{}, revokedPermissionsComment)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}
