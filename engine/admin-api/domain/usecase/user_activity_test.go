//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

type userActivitySuite struct {
	suite.Suite
	userActivity     *usecase.UserActivityInteractor
	logger           logr.Logger
	userActivityRepo *mocks.MockUserActivityRepo
	accessControl    *mocks.MockAccessControl
}

func TestUserActivitySuite(t *testing.T) {
	suite.Run(t, new(userActivitySuite))
}

func (s *userActivitySuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.userActivityRepo = mocks.NewMockUserActivityRepo(ctrl)
	s.accessControl = mocks.NewMockAccessControl(ctrl)

	s.userActivity = usecase.NewUserActivityInteractor(
		s.logger,
		s.userActivityRepo,
		s.accessControl,
	)
}

func (s *userActivitySuite) TestGet() {
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		userEmail  = "user@test.com"
		types      = []entity.UserActivityType{entity.UserActivityTypeStartVersion}
		versionIDs = []string{"test-version"}
		fromDate   = time.Now().String()
		toDate     = time.Now().String()
		lastID     = "last-id"

		expectedUserActivities = []*entity.UserActivity{
			{
				ID:     "test-id",
				Date:   time.Now(),
				UserID: user.ID,
				Type:   entity.UserActivityTypeStartVersion,
				Vars:   nil,
			},
		}
	)

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActViewUserActivities).Return(nil)
	s.userActivityRepo.EXPECT().Get(ctx, &userEmail, types, versionIDs, &fromDate, &toDate, &lastID).Return(expectedUserActivities, nil)

	actual, err := s.userActivity.Get(ctx, user, &userEmail, types, versionIDs, &fromDate, &toDate, &lastID)
	s.Assert().NoError(err)
	s.Assert().ElementsMatch(expectedUserActivities, actual)
}

func (s *userActivitySuite) TestGet_UnauthorizedUser() {
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		userEmail  = "user@test.com"
		types      = []entity.UserActivityType{entity.UserActivityTypeStartVersion}
		versionIDs = []string{"test-version"}
		fromDate   = time.Now().String()
		toDate     = time.Now().String()
		lastID     = "last-id"
	)

	expectedError := errors.New("access control error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActViewUserActivities).Return(expectedError)

	_, err := s.userActivity.Get(ctx, user, &userEmail, types, versionIDs, &fromDate, &toDate, &lastID)
	s.Assert().ErrorIs(err, expectedError)
}

func (s *userActivitySuite) TestRegisterCreateProduct() {
	product := &entity.Product{
		ID:          "test-product",
		Name:        "test-product",
		Description: "Test description",
	}

	user := testhelpers.NewUserBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: user.ID,
		Type:   entity.UserActivityTypeCreateProduct,
		Vars: []*entity.UserActivityVar{
			{
				Key:   "PRODUCT_ID",
				Value: product.ID,
			},
			{
				Key:   "PRODUCT_NAME",
				Value: product.Name,
			},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterCreateProduct(user.ID, product)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterCreateProduct_RepositoryError() {
	product := &entity.Product{
		ID:          "test-product",
		Name:        "test-product",
		Description: "Test description",
	}

	user := testhelpers.NewUserBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: user.ID,
		Type:   entity.UserActivityTypeCreateProduct,
		Vars: []*entity.UserActivityVar{
			{
				Key:   "PRODUCT_ID",
				Value: product.ID,
			},
			{
				Key:   "PRODUCT_NAME",
				Value: product.Name,
			},
		},
	}

	expectedError := errors.New("repository error")

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(expectedError)

	err := s.userActivity.RegisterCreateProduct(user.ID, product)
	s.Assert().ErrorIs(err, expectedError)
}

func (s *userActivitySuite) TestRegisterCreateAction() {
	const (
		userID    = "test-user"
		productID = "test-product"
	)

	version := testhelpers.NewVersionBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypeCreateVersion,
		Vars: []*entity.UserActivityVar{
			{Key: "PRODUCT_ID", Value: productID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_TAG", Value: version.Tag},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterCreateAction(userID, productID, &version)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterStartAction() {
	const (
		userID    = "test-user"
		productID = "test-product"
		comment   = "This is a test comment"
	)

	version := testhelpers.NewVersionBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypeStartVersion,
		Vars: []*entity.UserActivityVar{
			{Key: "PRODUCT_ID", Value: productID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_TAG", Value: version.Tag},
			{Key: "COMMENT", Value: comment},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterStartAction(userID, productID, &version, comment)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterStopAction() {
	const (
		userID    = "test-user"
		productID = "test-product"
		comment   = "This is a test comment"
	)

	version := testhelpers.NewVersionBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypeStopVersion,
		Vars: []*entity.UserActivityVar{
			{Key: "PRODUCT_ID", Value: productID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_TAG", Value: version.Tag},
			{Key: "COMMENT", Value: comment},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterStopAction(userID, productID, &version, comment)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterPublishAction() {
	const (
		userID    = "test-user"
		productID = "test-product"
		comment   = "This is a test comment"
	)

	newVersion := testhelpers.NewVersionBuilder().Build()
	previousVersion := testhelpers.NewVersionBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypePublishVersion,
		Vars: []*entity.UserActivityVar{
			{Key: "PRODUCT_ID", Value: productID},
			{Key: "VERSION_ID", Value: newVersion.ID},
			{Key: "VERSION_TAG", Value: newVersion.Tag},
			{Key: "OLD_PUBLISHED_VERSION_ID", Value: previousVersion.ID},
			{Key: "OLD_PUBLISHED_VERSION_TAG", Value: previousVersion.Tag},
			{Key: "COMMENT", Value: comment},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterPublishAction(userID, productID, &newVersion, &previousVersion, comment)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterUnpublishAction() {
	const (
		userID    = "test-user"
		productID = "test-product"
		comment   = "This is a test comment"
	)

	version := testhelpers.NewVersionBuilder().Build()

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypeUnpublishVersion,
		Vars: []*entity.UserActivityVar{
			{Key: "PRODUCT_ID", Value: productID},
			{Key: "VERSION_ID", Value: version.ID},
			{Key: "VERSION_TAG", Value: version.Tag},
			{Key: "COMMENT", Value: comment},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterUnpublishAction(userID, productID, &version, comment)
	s.Assert().NoError(err)
}

func (s *userActivitySuite) TestRegisterUpdateProductGrants() {
	const (
		userID       = "test-user"
		targetUserID = "target-user"
		productID    = "test-product"
		comment      = "This is a test comment"
	)

	productGrants := []string{"test", "view"}

	expectedUserActivity := entity.UserActivity{
		UserID: userID,
		Type:   entity.UserActivityTypeUpdateProductGrants,
		Vars: []*entity.UserActivityVar{
			{Key: "USER_ID", Value: userID},
			{Key: "TARGET_USER_ID", Value: targetUserID},
			{Key: "PRODUCT", Value: productID},
			{Key: "NEW_PRODUCT_GRANTS", Value: strings.Join(productGrants, ",")},
			{Key: "COMMENT", Value: comment},
		},
	}

	customMatcher := newUserActivityMatcher(expectedUserActivity)

	s.userActivityRepo.EXPECT().Create(customMatcher).Return(nil)

	err := s.userActivity.RegisterUpdateProductGrants(userID, targetUserID, productID, productGrants, comment)
	s.Assert().NoError(err)
}

type userActivityMatcher struct {
	expectedUserActivity entity.UserActivity
}

func newUserActivityMatcher(expectedUserActivity entity.UserActivity) *userActivityMatcher {
	return &userActivityMatcher{
		expectedUserActivity: expectedUserActivity,
	}
}

func (m userActivityMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedUserActivity)
}

func (m userActivityMatcher) Matches(actual interface{}) bool {
	actualUserActivity, ok := actual.(entity.UserActivity)
	if !ok {
		return false
	}

	return m.expectedUserActivity.UserID == actualUserActivity.UserID &&
		m.expectedUserActivity.Type == actualUserActivity.Type &&
		reflect.DeepEqual(m.expectedUserActivity.Vars, actualUserActivity.Vars)
}
