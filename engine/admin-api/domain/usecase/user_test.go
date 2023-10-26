//go:build unit

package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const (
	targetUserID = "trigger-user-id"
	testProduct  = "test-product"
)

type ContextUserManagerSuite struct {
	suite.Suite
	mockUserRegistry           *mocks.MockUserRegistry
	mockAccessControl          *mocks.MockAccessControl
	mockLogger                 *mocks.MockLogger
	mockUserActivityInteractor *mocks.MockUserActivityInteracter
	userManager                *UserInteractor
}

func TestContextMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(ContextUserManagerSuite))
}

func (s *ContextUserManagerSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	s.mockUserRegistry = mocks.NewMockUserRegistry(mockController)
	s.mockAccessControl = mocks.NewMockAccessControl(mockController)
	s.mockLogger = mocks.NewMockLogger(mockController)
	s.mockUserActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	s.userManager = NewUserInteractor(s.mockLogger, s.mockAccessControl, s.mockUserActivityInteractor, s.mockUserRegistry)

	mocks.AddLoggerExpects(s.mockLogger)
}

func (s *ContextUserManagerSuite) GetTestUser() *entity.User {
	return testhelpers.NewUserBuilder().Build()
}

func (s *ContextUserManagerSuite) TestUpdateUserProductGrants() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	grants := []string{"grant1", "grant2"}

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, grants).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		"",
	).Times(1).Return(nil)

	err := s.userManager.UpdateUserProductGrants(ctx, loggedUser, targetUserID, testProduct, grants)
	s.NoError(err)
}

func (s *ContextUserManagerSuite) TestUpdateUserProductGrantsGivenComment() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	grants := []string{"grant1", "grant2"}
	testComment := "test comment"

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, grants).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		testComment,
	).Times(1).Return(nil)

	err := s.userManager.UpdateUserProductGrants(ctx, loggedUser, targetUserID, testProduct, grants, testComment)
	s.NoError(err)
}

func (s *ContextUserManagerSuite) TestUpdateProductGrantsErrorUnauthorized() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	grants := []string{"grant1", "grant2"}
	exppectedError := errors.New("unauthorized")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(exppectedError)

	err := s.userManager.UpdateUserProductGrants(ctx, loggedUser, targetUserID, testProduct, grants)
	s.Error(err)
	s.ErrorIs(err, exppectedError)
}

func (s *ContextUserManagerSuite) TestUpdateUserProductGrantsErrorInUserRegistry() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	grants := []string{"grant1", "grant2"}

	expectedError := errors.New("user registry error")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, grants).Times(1).
		Return(expectedError)

	err := s.userManager.UpdateUserProductGrants(ctx, loggedUser, targetUserID, testProduct, grants)
	s.ErrorIs(err, expectedError)
}

func (s *ContextUserManagerSuite) TestUpdateUserGrantsErrorInUserActivity() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	grants := []string{"grant1", "grant2"}

	expecetedError := errors.New("user activity error")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, grants).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		"",
	).Times(1).Return(expecetedError)

	err := s.userManager.UpdateUserProductGrants(ctx, loggedUser, targetUserID, testProduct, grants)
	s.ErrorIs(err, expecetedError)
}

func (s *ContextUserManagerSuite) TestRevokeProductGrants() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	testComment := "test comment"

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, []string{}).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		testComment,
	).Times(1).Return(nil)

	err := s.userManager.RevokeUserProductGrants(ctx, loggedUser, targetUserID, testProduct, testComment)
	s.NoError(err)
}

func (s *ContextUserManagerSuite) TestRevokeProductGrantsGivenComment() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, []string{}).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(nil)

	err := s.userManager.RevokeUserProductGrants(ctx, loggedUser, targetUserID, testProduct)
	s.NoError(err)
}

func (s *ContextUserManagerSuite) TestRevokeProductGrantsErrorUnauthorized() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	exppectedError := errors.New("unauthorized")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(exppectedError)

	err := s.userManager.RevokeUserProductGrants(ctx, loggedUser, targetUserID, testProduct)
	s.Error(err)
	s.ErrorIs(err, exppectedError)
}

func (s *ContextUserManagerSuite) TestRevokeProductGrantsErrorInUserRegistry() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	expectedError := errors.New("user registry error")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, []string{}).Times(1).
		Return(expectedError)

	err := s.userManager.RevokeUserProductGrants(ctx, loggedUser, targetUserID, testProduct)
	s.Error(err)
}

func (s *ContextUserManagerSuite) TestRevokeUserGrantsErrorInUserActivity() {
	ctx := context.Background()
	loggedUser := s.GetTestUser()
	expectedError := errors.New("user activity error")

	s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	s.mockUserRegistry.EXPECT().UpdateUserProductGrants(ctx, targetUserID, testProduct, []string{}).Times(1).Return(nil)
	s.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(expectedError)

	err := s.userManager.RevokeUserProductGrants(ctx, loggedUser, targetUserID, testProduct)
	s.ErrorIs(err, expectedError)
}
