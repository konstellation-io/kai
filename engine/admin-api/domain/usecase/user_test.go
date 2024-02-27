//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const (
	_targetUserEmail = "test-user-email"
	_testProduct     = "test-product"
)

type userHandlerSuite struct {
	suite.Suite
	mockUserRegistry           *mocks.MockUserRegistry
	mockAccessControl          *mocks.MockAccessControl
	logger                     logr.Logger
	mockUserActivityInteractor *mocks.MockUserActivityInteracter
	userHandler                *usecase.UserHandler
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(userHandlerSuite))
}

func (s *userHandlerSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	s.mockUserRegistry = mocks.NewMockUserRegistry(mockController)
	s.mockAccessControl = mocks.NewMockAccessControl(mockController)
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.mockUserActivityInteractor = mocks.NewMockUserActivityInteracter(mockController)
	s.userHandler = usecase.NewUserHandler(s.logger, s.mockAccessControl, s.mockUserActivityInteractor, s.mockUserRegistry)
}

func (s *userHandlerSuite) getTestUser() *entity.User {
	return testhelpers.NewUserBuilder().Build()
}

func (s *userHandlerSuite) TestAddUserToProduct() {
	var (
		ctx        = context.Background()
		loggedUser = s.getTestUser()
		grants     = auth.GetDefaultUserGrants()
	)

	s.T().Run("OK", func(t *testing.T) {
		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(nil)
		s.mockUserRegistry.EXPECT().AddProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(nil)

		err := s.userHandler.AddUserToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.NoError(err)
	})

	s.T().Run("Unauthorized", func(t *testing.T) {
		expectedError := errors.New("unauthorized")

		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(expectedError)

		err := s.userHandler.AddUserToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})

	s.T().Run("UserRegistryFails", func(t *testing.T) {
		expectedError := errors.New("user registry error")

		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(nil)
		s.mockUserRegistry.EXPECT().AddProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(expectedError)

		err := s.userHandler.AddUserToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})
}

func (s *userHandlerSuite) TestRemoveUserFromProduct() {
	var (
		ctx        = context.Background()
		loggedUser = s.getTestUser()
		grants     = auth.GetDefaultUserGrants()
	)

	s.T().Run("OK", func(t *testing.T) {
		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(nil)
		s.mockUserRegistry.EXPECT().RevokeProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(nil)

		err := s.userHandler.RemoveUserFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.NoError(err)
	})

	s.T().Run("Unauthorized", func(t *testing.T) {
		expectedError := errors.New("unauthorized")

		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(expectedError)

		err := s.userHandler.RemoveUserFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})

	s.T().Run("UserRegistryFails", func(t *testing.T) {
		expectedError := errors.New("user registry error")

		s.mockAccessControl.EXPECT().CheckProductGrants(loggedUser, _testProduct, auth.ActManageProductUsers).Return(nil)
		s.mockUserRegistry.EXPECT().RevokeProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(expectedError)

		err := s.userHandler.RemoveUserFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})
}

func (s *userHandlerSuite) TestAddMaintainerToProduct() {
	var (
		ctx        = context.Background()
		loggedUser = s.getTestUser()
		grants     = auth.GetDefaultMaintainerGrants()
	)

	s.T().Run("OK", func(t *testing.T) {
		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(nil)
		s.mockUserRegistry.EXPECT().AddProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(nil)

		err := s.userHandler.AddMaintainerToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.NoError(err)
	})

	s.T().Run("Unauthorized", func(t *testing.T) {
		expectedError := errors.New("unauthorized")

		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(expectedError)

		err := s.userHandler.AddMaintainerToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})

	s.T().Run("UserRegistryFails", func(t *testing.T) {
		expectedError := errors.New("user registry error")

		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(nil)

		s.mockUserRegistry.EXPECT().AddProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(expectedError)

		err := s.userHandler.AddMaintainerToProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})
}

func (s *userHandlerSuite) TestRemoveMaintainerFromProduct() {
	var (
		ctx        = context.Background()
		loggedUser = s.getTestUser()
		grants     = auth.GetDefaultMaintainerGrants()
	)

	s.T().Run("OK", func(t *testing.T) {
		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(nil)
		s.mockUserRegistry.EXPECT().RevokeProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(nil)

		err := s.userHandler.RemoveMaintainerFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.NoError(err)
	})

	s.T().Run("Unauthorized", func(t *testing.T) {
		expectedError := errors.New("unauthorized")

		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(expectedError)

		err := s.userHandler.RemoveMaintainerFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})

	s.T().Run("UserRegistryFails", func(t *testing.T) {
		expectedError := errors.New("user registry error")

		s.mockAccessControl.EXPECT().CheckRoleGrants(loggedUser, auth.ActManageProductMaintainers).Return(nil)
		s.mockUserRegistry.EXPECT().RevokeProductGrants(ctx, _targetUserEmail, _testProduct, grants).Times(1).Return(expectedError)

		err := s.userHandler.RemoveMaintainerFromProduct(ctx, loggedUser, _targetUserEmail, _testProduct)
		s.ErrorIs(err, expectedError)
	})
}
