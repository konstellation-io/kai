//go:build unit

package usecase

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const (
	targetUserID = "trigger-user-id"
	testProduct  = "test-product"
)

type ContextUserManagerSuite struct {
	testifySuite.Suite
	mockUserRegistry           *mocks.MockUserRegistry
	mockAccessControl          *mocks.MockAccessControl
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
	suite.mockAccessControl = mocks.NewMockAccessControl(mockController)
	suite.mockLogger = mocks.NewMockLogger(mockController)
	suite.mockUserActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	suite.userManager = NewUserInteractor(suite.mockLogger, suite.mockAccessControl, suite.mockUserActivityInteractor, suite.mockUserRegistry)
}

func (suite *ContextUserManagerSuite) GetTestUser() *entity.User {
	return testhelpers.NewUserBuilder().Build()
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductGrants() {
	loggedUser := suite.GetTestUser()
	grants := []string{"grant1", "grant2"}

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, grants).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		"",
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.UpdateUserProductGrants(loggedUser, targetUserID, testProduct, grants)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductGrantsGivenComment() {
	loggedUser := suite.GetTestUser()
	grants := []string{"grant1", "grant2"}
	testComment := "test comment"

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, grants).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		testComment,
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.UpdateUserProductGrants(loggedUser, targetUserID, testProduct, grants, testComment)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestUpdateProductGrantsErrorUnauthorized() {
	loggedUser := suite.GetTestUser()
	grants := []string{"grant1", "grant2"}
	exppectedError := errors.New("unauthorized")

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(exppectedError)

	err := suite.userManager.UpdateUserProductGrants(loggedUser, targetUserID, testProduct, grants)
	suite.Error(err)
	suite.ErrorIs(err, exppectedError)
}

func (suite *ContextUserManagerSuite) TestUpdateUserProductGrantsErrorInUserRegistry() {
	loggedUser := suite.GetTestUser()
	grants := []string{"grant1", "grant2"}

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, grants).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserProductGrants(loggedUser, targetUserID, testProduct, grants)
	suite.Error(err)
	suite.ErrorContains(err, updateUserProductGrantsWrapper)
}

func (suite *ContextUserManagerSuite) TestUpdateUserGrantsErrorInUserActivity() {
	loggedUser := suite.GetTestUser()
	grants := []string{"grant1", "grant2"}

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, grants).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		grants,
		"",
	).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.UpdateUserProductGrants(loggedUser, targetUserID, testProduct, grants)
	suite.Error(err)
	suite.ErrorContains(err, updateUserProductGrantsWrapper)
}

func (suite *ContextUserManagerSuite) TestRevokeProductGrants() {
	loggedUser := suite.GetTestUser()
	testComment := "test comment"

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		testComment,
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.RevokeUserProductGrants(loggedUser, targetUserID, testProduct, testComment)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductGrantsGivenComment() {
	loggedUser := suite.GetTestUser()

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(nil)
	mocks.AddLoggerExpects(suite.mockLogger)

	err := suite.userManager.RevokeUserProductGrants(loggedUser, targetUserID, testProduct)
	suite.NoError(err)
}

func (suite *ContextUserManagerSuite) TestRevokeProductGrantsErrorUnauthorized() {
	loggedUser := suite.GetTestUser()
	exppectedError := errors.New("unauthorized")

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(exppectedError)

	err := suite.userManager.RevokeUserProductGrants(loggedUser, targetUserID, testProduct)
	suite.Error(err)
	suite.ErrorIs(err, exppectedError)
}

func (suite *ContextUserManagerSuite) TestRevokeProductGrantsErrorInUserRegistry() {
	loggedUser := suite.GetTestUser()

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, []string{}).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeUserProductGrants(loggedUser, targetUserID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, revokeUserProductGrantsWrapper)
}

func (suite *ContextUserManagerSuite) TestRevokeUserGrantsErrorInUserActivity() {
	loggedUser := suite.GetTestUser()

	suite.mockAccessControl.EXPECT().CheckGrants(loggedUser, auth.ActUpdateUserGrants).Return(nil)
	suite.mockUserRegistry.EXPECT().UpdateUserProductGrants(targetUserID, testProduct, []string{}).Times(1).Return(nil)
	suite.mockUserActivityInteractor.EXPECT().RegisterUpdateProductGrants(
		loggedUser.ID,
		targetUserID,
		testProduct,
		[]string{},
		"",
	).Times(1).Return(fmt.Errorf("error"))

	err := suite.userManager.RevokeUserProductGrants(loggedUser, targetUserID, testProduct)
	suite.Error(err)
	suite.ErrorContains(err, revokeUserProductGrantsWrapper)
}
