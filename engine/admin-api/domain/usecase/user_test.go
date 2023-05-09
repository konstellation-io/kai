package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kre/engine/admin-api/mocks"
)

type ContextUserManagerSuite struct {
	suite.Suite
	mockGokeycloak         *mocks.MockGocloakService
	logger                 *mocks.MockLogger
	userActivityInteractor *mocks.MockUserActivityInteracter
	userManager            *UserInteractor
}

func TestContextMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(ContextUserManagerSuite))
}

func (suite *ContextUserManagerSuite) SetupSuite() {
	mockController := gomock.NewController(suite.T())
	suite.mockGokeycloak = mocks.NewMockGocloakService(mockController)
	suite.logger = mocks.NewMockLogger(mockController)
	suite.userActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	suite.userManager = NewUserInteractor(suite.logger, suite.userActivityInteractor, suite.mockGokeycloak)
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

func (suite *ContextUserManagerSuite) TestCreateUser() {
	testUserData := suite.GetTestUserData()
	suite.mockGokeycloak.EXPECT().CreateUser(testUserData).Times(1).Return(nil)
	err := suite.userManager.CreateUser(testUserData)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestCreateUserError() {
	testUserData := suite.GetTestUserData()
	suite.mockGokeycloak.EXPECT().CreateUser(testUserData).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.CreateUser(testUserData)
	suite.Error(err)
	suite.ErrorContains(err, "create user")
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
	suite.ErrorContains(err, "get user by id")
}

func (suite *ContextUserManagerSuite) TestUpdateUserRoles() {
	testUserData := suite.GetTestUserData()
	testProduct := "test-product"
	roles := []string{"role1", "role2"}
	suite.mockGokeycloak.EXPECT().UpdateUserRoles(testUserData.ID, testProduct, roles).Times(1).Return(nil)

	err := suite.userManager.UpdateUserRoles(testUserData.ID, testProduct, roles)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateUserRolesError() {
	testUserData := suite.GetTestUserData()
	testProduct := "test-product"
	roles := []string{"role1", "role2"}
	suite.mockGokeycloak.EXPECT().UpdateUserRoles(testUserData.ID, testProduct, roles).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserRoles(testUserData.ID, testProduct, roles)
	suite.Error(err)
	suite.ErrorContains(err, "update user roles")
}

func (suite *ContextUserManagerSuite) TestRevokeProductRoles() {
	testUserData := suite.GetTestUserData()
	testProduct := "test-product"
	suite.mockGokeycloak.EXPECT().UpdateUserRoles(testUserData.ID, testProduct, []string{}).Times(1).Return(nil)

	err := suite.userManager.RevokeProductRoles(testUserData.ID, testProduct)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductRolesError() {
	testUserData := suite.GetTestUserData()
	testProduct := "test-product"
	suite.mockGokeycloak.EXPECT().UpdateUserRoles(testUserData.ID, testProduct, []string{}).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeProductRoles(testUserData.ID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, "revoke user roles")
}
