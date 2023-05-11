package usecase

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

// UserInteractor contains app logic to handle User entities.
type UserInteractor struct {
	logger                 logging.Logger
	userActivityInteractor UserActivityInteracter
	accessControl          auth.AccessControl
}

// NewUserInteractor creates a new UserInteractor.
func NewUserInteractor(
	logger logging.Logger,
	userActivityInteractor UserActivityInteracter,
	accessControl auth.AccessControl,
) *UserInteractor {
	return &UserInteractor{
		logger,
		userActivityInteractor,
		accessControl,
	}
}
