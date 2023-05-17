package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
)

const (
	triggerUserID = "trigger-user-id"
	testProduct   = "test-product"
)

type ContextUserManagerSuite struct {
	testifySuite.Suite
	mockUserRegistry           *mocks.MockUserRegistry
	mockLogger                 *mocks.MockLogger
	mockUserActivityInteractor *mocks.MockUserActivityInteracter
	userManager                *UserInteractor
}

func TestContextMeasurementTestSuite(t *testing.T) {
	testifySuite.Run(t, new(ContextUserManagerSuite))
}

func (suite *ContextUserManagerSuite) SetupSuite() {
	mockController := gomock.NewController(suite.T())
	suite.mockUserRegistry = mocks.NewMockUserRegistry(mockController)
	suite.mockLogger = mocks.NewMockLogger(mockController)
	suite.mockUserActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	suite.userManager = NewUserInteractor(suite.mockLogger, suite.mockUserActivityInteractor, suite.mockUserRegistry)
}

func (suite *ContextUserManagerSuite) GetTestUser() entity.User {
	return entity.User{
		ID:        "test-id",
		Username:  "test",
		Email:     "test@email.com",
		FirstName: "test name",
		LastName:  "test last name",
		Enabled:   true,
	}
}

func (suite *ContextUserManagerSuite) TestGetUserByID() {
	testUser := suite.GetTestUser()

	suite.mockUserRegistry.EXPECT().GetUserByID(testUser.ID).Times(1).Return(&testUser, nil)

	user, err := suite.userManager.GetUserByID(testUser.ID)
	suite.NoError(err)
	suite.Equal(testUser, *user)
}

func (suite *ContextUserManagerSuite) TestGetUserByIDErrorInUserRegistry() {
	testUser := suite.GetTestUser()

	suite.mockUserRegistry.EXPECT().GetUserByID(testUser.ID).Times(1).Return(nil, fmt.Errorf("error"))

	_, err := suite.userManager.GetUserByID(testUser.ID)
	suite.Error(err)
	suite.ErrorContains(err, getUserByIDWrapper)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductPermissions() {
	testUser := suite.GetTestUser()
	permissions := []string{"permission1", "permission2"}

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, permissions).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		permissions,
		"",
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUser.ID, testProduct, permissions)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductPermissionsGivenComment() {
	testUser := suite.GetTestUser()
	permissions := []string{"permission1", "permission2"}
	testComment := "test comment"

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, permissions).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		permissions,
		testComment,
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUser.ID, testProduct, permissions, testComment)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductPermissionsErrorInUserRegistry() {
	testUser := suite.GetTestUser()
	permissions := []string{"permission1", "permission2"}

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, permissions).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUser.ID, testProduct, permissions)
	suite.Error(err)
	suite.ErrorContains(err, updateUserProductPermissionsWrapper)
}

func (suite *ContextUserManagerSuite) TestUpdateUserPermissionsErrorInUserActivity() {
	testUser := suite.GetTestUser()
	permissions := []string{"permission1", "permission2"}

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, permissions).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		permissions,
		"",
	).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUser.ID, testProduct, permissions)
	suite.Error(err)
	suite.ErrorContains(err, updateUserProductPermissionsWrapper)
}

func (suite *ContextUserManagerSuite) TestRevokeProductPermissions() {
	testUser := suite.GetTestUser()
	testComment := "test comment"

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		[]string{},
		testComment,
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUser.ID, testProduct, testComment)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductPermissionsGivenComment() {
	testUser := suite.GetTestUser()

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUser.ID, testProduct)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductPermissionsErrorInUserRegistry() {
	testUser := suite.GetTestUser()

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, []string{}).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUser.ID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, revokeUserProductPermissionsWrapper)
}

func (suite *ContextUserManagerSuite) TestRevokeUserPermissionsErrorInUserActivity() {
	testUser := suite.GetTestUser()

	suite.mockUserRegistry.EXPECT().UpdateUserProductPermissions(testUser.ID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(
		triggerUserID,
		testUser.ID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUser.ID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, revokeUserProductPermissionsWrapper)
}
