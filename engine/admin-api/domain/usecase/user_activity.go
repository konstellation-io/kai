package usecase

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/usecase_${GOFILE} -package=mocks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

type UserActivityInteracter interface {
	Get(ctx context.Context, loggedUserID string, userEmail *string, types []entity.UserActivityType,
		versionIDs []string, fromDate *string, toDate *string, lastID *string) ([]*entity.UserActivity, error)
	RegisterCreateRuntime(userID string, runtime *entity.Runtime) error
	RegisterCreateAction(userID, runtimeID string, version *entity.Version) error
	RegisterStartAction(userID, runtimeID string, version *entity.Version, comment string) error
	RegisterStopAction(userID, runtimeID string, version *entity.Version, comment string) error
	RegisterPublishAction(userID, runtimeID string, version *entity.Version, prev *entity.Version, comment string) error
	RegisterUnpublishAction(userID, runtimeID string, version *entity.Version, comment string) error
	RegisterUpdateProductGrants(userID string, targetUserID string, product string, productGrants []string, comment string) error
}

// UserActivityInteractor  contains app logic about UserActivity entities.
type UserActivityInteractor struct {
	logger           logging.Logger
	userActivityRepo repository.UserActivityRepo
	accessControl    auth.AccessControl
}

// NewUserActivityInteractor creates a new UserActivityInteractor.
func NewUserActivityInteractor(
	logger logging.Logger,
	userActivityRepo repository.UserActivityRepo,
	accessControl auth.AccessControl,
) UserActivityInteracter {
	return &UserActivityInteractor{
		logger,
		userActivityRepo,
		accessControl,
	}
}

// Get return a list of UserActivities.
func (i *UserActivityInteractor) Get(
	ctx context.Context,
	loggedUserID string,
	userEmail *string,
	types []entity.UserActivityType,
	versionIDs []string,
	fromDate *string,
	toDate *string,
	lastID *string,
) ([]*entity.UserActivity, error) {
	if err := i.accessControl.CheckGrant(loggedUserID, auth.ResAudit, auth.ActView); err != nil {
		return nil, err
	}

	return i.userActivityRepo.Get(ctx, userEmail, types, versionIDs, fromDate, toDate, lastID)
}

// Create add a new UserActivity to the given user.
func (i *UserActivityInteractor) create(
	userID string,
	userActivityType entity.UserActivityType,
	vars []*entity.UserActivityVar,
) error {
	userActivity := entity.UserActivity{
		ID:     primitive.NewObjectID().Hex(),
		UserID: userID,
		Type:   userActivityType,
		Date:   time.Now(),
		Vars:   vars,
	}

	return i.userActivityRepo.Create(userActivity)
}

func checkUserActivityError(logger logging.Logger, err error) error {
	if err != nil {
		userActivityErr := fmt.Errorf("error creating userActivity: %w", err)
		logger.Error(userActivityErr.Error())

		return userActivityErr
	}

	return nil
}

func (i *UserActivityInteractor) RegisterCreateRuntime(
	userID string,
	runtime *entity.Runtime,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeCreateProduct,
		[]*entity.UserActivityVar{
			{
				Key:   "PRODUCT_ID",
				Value: runtime.ID,
			},
			{
				Key:   "RUNTIME_NAME",
				Value: runtime.Name,
			},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterCreateAction(
	userID,
	runtimeID string,
	version *entity.Version,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeCreateVersion,
		[]*entity.UserActivityVar{
			{Key: "RUNTIME_ID", Value: runtimeID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_NAME", Value: version.Name},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterStartAction(
	userID,
	runtimeID string,
	version *entity.Version,
	comment string,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeStartVersion,
		[]*entity.UserActivityVar{
			{Key: "RUNTIME_ID", Value: runtimeID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_NAME", Value: version.Name},
			{Key: "COMMENT", Value: comment},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterStopAction(
	userID,
	runtimeID string,
	version *entity.Version,
	comment string,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeStopVersion,
		[]*entity.UserActivityVar{
			{Key: "RUNTIME_ID", Value: runtimeID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_NAME", Value: version.Name},
			{Key: "COMMENT", Value: comment},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterPublishAction(
	userID,
	runtimeID string,
	version *entity.Version,
	prev *entity.Version,
	comment string,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypePublishVersion,
		[]*entity.UserActivityVar{
			{Key: "RUNTIME_ID", Value: runtimeID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_NAME", Value: version.Name},
			{Key: "OLD_PUBLISHED_VERSION_ID", Value: prev.ID},
			{Key: "OLD_PUBLISHED_VERSION_NAME", Value: prev.Name},
			{Key: "COMMENT", Value: comment},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterUnpublishAction(
	userID,
	runtimeID string,
	version *entity.Version,
	comment string,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeUnpublishVersion,
		[]*entity.UserActivityVar{
			{Key: "RUNTIME_ID", Value: runtimeID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_NAME", Value: version.Name},
			{Key: "COMMENT", Value: comment},
		})

	return checkUserActivityError(i.logger, err)
}

func (i *UserActivityInteractor) RegisterUpdateProductGrants(
	userID string,
	targetUserID string,
	product string,
	productGrants []string,
	comment string,
) error {
	err := i.create(
		userID,
		entity.UserActivityTypeUpdateProductGrants,
		[]*entity.UserActivityVar{
			{Key: "TARGET_USER_ID", Value: targetUserID},
			{Key: "PRODUCT", Value: product},
			{Key: "NEW_PRODUCT_GRANTS", Value: strings.Join(productGrants, ",")},
			{Key: "COMMENT", Value: comment},
		})

	return checkUserActivityError(i.logger, err)
}
