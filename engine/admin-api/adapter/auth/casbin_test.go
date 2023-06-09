//go:build unit

package auth_test

import (
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	auth2 "github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	casbinModel  = path.Join("..", "..", "casbin_rbac_model.conf")
	casbinPolicy = path.Join("..", "..", "casbin_rbac_policy.csv")
)

func TestCheckProductGrants(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := auth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	product01 := "product-01"
	product02 := "product-02"

	testCases := []struct {
		name     string
		product  string
		user     *entity.User
		act      auth2.AccessControlAction
		hasError bool
	}{
		{
			name:    "user with grants to view product-01 can view product-01",
			product: product01,
			user: testhelpers.NewUserBuilder().
				WithProductGrants(
					entity.ProductGrants{
						product01: []string{auth2.ActViewProduct.String()},
					},
				).
				Build(),
			act:      auth2.ActViewProduct,
			hasError: false,
		},
		{
			name:     "user without grants to view product-01 cannot view product-01",
			product:  product01,
			user:     testhelpers.NewUserBuilder().Build(),
			act:      auth2.ActViewProduct,
			hasError: true,
		},
		{
			name:    "user with grant to view product-01 but no product-02 cannot view product-02",
			product: product02,
			user: testhelpers.NewUserBuilder().WithProductGrants(
				entity.ProductGrants{
					product01: []string{auth2.ActViewProduct.String()},
				},
			).Build(),
			act:      auth2.ActViewProduct,
			hasError: true,
		},
		{
			name:    "user with grant to view product-01 cannot create product-01",
			product: product01,
			user: testhelpers.NewUserBuilder().WithProductGrants(
				entity.ProductGrants{
					product01: []string{auth2.ActViewProduct.String()},
				},
			).Build(),
			act:      auth2.AccessControlAction(auth2.ActCreateProduct.String()),
			hasError: true,
		},
		{
			name:     "admin user can do anything without specifying product grants",
			product:  product01,
			user:     testhelpers.NewUserBuilder().WithRoles([]string{auth2.DefaultAdminRole}).Build(),
			act:      auth2.AccessControlAction(auth2.ActCreateProduct.String()),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := authorizer.CheckProductGrants(tc.user, tc.product, tc.act)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckAdminGrants(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := auth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		user     *entity.User
		act      auth2.AccessControlAction
		hasError bool
	}{
		{
			name:     "user with admin grants can update user roles",
			user:     testhelpers.NewUserBuilder().WithRoles([]string{auth2.DefaultAdminRole}).Build(),
			act:      auth2.ActUpdateUserGrants,
			hasError: false,
		},
		{
			name:     "user without admin grants cannot update user roles",
			user:     testhelpers.NewUserBuilder().WithRoles([]string{}).Build(),
			act:      auth2.ActUpdateUserGrants,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := authorizer.CheckAdminGrants(tc.user, tc.act)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsadmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := auth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		user    *entity.User
		isAdmin bool
	}{
		{
			name:    "user with cfg admin role in its roles is admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{auth2.DefaultAdminRole}).Build(),
			isAdmin: true,
		},
		{
			name:    "user without cfg admin role in its roles is not admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{"user", "maintainer"}).Build(),
			isAdmin: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := authorizer.IsAdmin(tc.user)
			assert.Equal(t, tc.isAdmin, isAdmin)
		})
	}
}

func TestIsadmin_WithOptAdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	customAdminRole := "admin"

	authorizer, err := auth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy, auth.WithAdminRole(customAdminRole))
	require.NoError(t, err)

	testCases := []struct {
		name    string
		user    *entity.User
		isAdmin bool
	}{
		{
			name:    "user with cfg admin role in its roles is admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{customAdminRole}).Build(),
			isAdmin: true,
		},
		{
			name:    "user without cfg admin role in its roles is not admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{"user", auth2.DefaultAdminRole}).Build(),
			isAdmin: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := authorizer.IsAdmin(tc.user)
			assert.Equal(t, tc.isAdmin, isAdmin)
		})
	}
}

func TestCheckProductGrants_InvalidAct(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := auth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	user := testhelpers.NewUserBuilder().Build()
	product := "product-01"

	err = authorizer.CheckProductGrants(user, product, "invalid act")
	assert.ErrorIs(t, auth.ErrInvalidAccessControlAction, err)
}

func TestNewCasbinAccessControl_ErrorInitEnforcer(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	_, err := auth.NewCasbinAccessControl(logger, "this is a invalid model", casbinPolicy)
	require.Error(t, err)
}
