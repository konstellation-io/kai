package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kre/engine/admin-api/mocks"
)

const (
	triggerUserID = "trigger-user-id"
	testProduct   = "test-product"
)

type ContextUserManagerSuite struct {
	suite.Suite
	mockGokeycloak             *mocks.MockGocloakService
	mockLogger                 *mocks.MockLogger
	mockUserActivityInteractor *mocks.MockUserActivityInteracter
	userManager                *UserInteractor
}

func TestContextMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(ContextUserManagerSuite))
}

func (suite *ContextUserManagerSuite) SetupSuite() {
	mockController := gomock.NewController(suite.T())
	suite.mockGokeycloak = mocks.NewMockGocloakService(mockController)
	suite.mockLogger = mocks.NewMockLogger(mockController)
	suite.mockUserActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	suite.userManager = NewUserInteractor(suite.mockLogger, suite.mockUserActivityInteractor, suite.mockGokeycloak)
}

func (suite *ContextUserManagerSuite) GetTestUserData() entity.UserGocloakData {
	return entity.UserGocloakData{
		ID:        "test-id",
		Username:  "test",
		Email:     "test@email.com",
		FirstName: "test name",
		LastName:  "test last name",
		Enabled:   true,
	}
}

func (suite *ContextUserManagerSuite) TestGetUserByID() {
	testUserData := suite.GetTestUserData()

	suite.mockGokeycloak.EXPECT().GetUserByID(testUserData.ID).Times(1).Return(testUserData, nil)

	userData, err := suite.userManager.GetUserByID(testUserData.ID)
	suite.NoError(err)
	suite.Equal(testUserData, userData)
}

func (suite *ContextUserManagerSuite) TestGetUserByIDError() {
	testUserData := suite.GetTestUserData()

	suite.mockGokeycloak.EXPECT().GetUserByID(testUserData.ID).Times(1).Return(entity.UserGocloakData{}, fmt.Errorf("error"))

	_, err := suite.userManager.GetUserByID(testUserData.ID)
	suite.Error(err)
	suite.ErrorContains(err, getUserByIDWrapper)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductPermissions() {
	testUserData := suite.GetTestUserData()
	permissions := []string{"permission1", "permission2"}

	suite.mockGokeycloak.EXPECT().UpdateUserProductPermissions(testUserData.ID, testProduct, permissions).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(triggerUserID, testUserData.ID, testProduct, permissions, "").Times(1).Return(nil)

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUserData.ID, testProduct, permissions)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductPermissionsError() {
	testUserData := suite.GetTestUserData()
	permissions := []string{"permission1", "permission2"}

	suite.mockGokeycloak.EXPECT().UpdateUserProductPermissions(testUserData.ID, testProduct, permissions).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserProductPermissions(triggerUserID, testUserData.ID, testProduct, permissions)
	suite.Error(err)
	suite.ErrorContains(err, updateUserProductPermissionsWrapper)
}

func (suite *ContextUserManagerSuite) TestRevokeProductPermissions() {
	testUserData := suite.GetTestUserData()

	suite.mockGokeycloak.EXPECT().UpdateUserProductPermissions(testUserData.ID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductPermissions(triggerUserID, testUserData.ID, testProduct, []string{}, revokedPermissionsComment).Times(1).Return(nil)

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUserData.ID, testProduct)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductPermissionsError() {
	testUserData := suite.GetTestUserData()

	suite.mockGokeycloak.EXPECT().UpdateUserProductPermissions(testUserData.ID, testProduct, []string{}).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeUserProductPermissions(triggerUserID, testUserData.ID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, revokeUserProductPermissionsWrapper)
}
