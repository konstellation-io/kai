package usecase

import (
	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kre/engine/admin-api/domain/service"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kre/engine/admin-api/internal/errors"
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
	wrapErr := errors.Wrapper("get user by id: %w")

	user, err := ui.gocloakManager.GetUserByID(userID)
	if err != nil {
		return entity.UserGocloakData{}, wrapErr(err)
	}

	return user, nil
}

func (ui *UserInteractor) UpdateUserRoles(userID string, product string, roles []string) error {
	wrapErr := errors.Wrapper("update user roles: %w")

	err := ui.gocloakManager.UpdateUserRoles(userID, product, roles)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func (ui *UserInteractor) RevokeProductRoles(userID string, product string) error {
	wrapErr := errors.Wrapper("revoke user roles: %w")

	err := ui.gocloakManager.UpdateUserRoles(userID, product, []string{})
	if err != nil {
		return wrapErr(err)
	}

	return nil
}
