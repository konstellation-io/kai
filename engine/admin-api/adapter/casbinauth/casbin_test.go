//go:build unit

package casbinauth_test

import (
	"path"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	casbinModel  = path.Join("..", "..", "casbin_rbac_model.conf")
	casbinPolicy = path.Join("..", "..", "casbin_rbac_policy.csv")
)

func TestNewCasbinAccessControl_ErrorInitEnforcer(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	_, err := casbinauth.NewCasbinAccessControl(logger, "this is a invalid model", casbinPolicy)
	require.Error(t, err)
}

func TestGetUserProducts(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	expectedProducts := []string{"product-1", "product-2"}
	user := &entity.User{
		Roles: []string{"USER"},
		ProductGrants: map[string][]string{
			expectedProducts[0]:                {auth.ActViewProduct.String()},
			expectedProducts[1]:                {auth.ActViewProduct.String()},
			"product-without-view-permissions": {},
		},
	}

	userProducts := authorizer.GetUserProductsWithViewAccess(user)
	for _, p := range expectedProducts {
		assert.Contains(t, userProducts, p)
	}
}

func TestGetUserProducts_AdminUser(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	user := &entity.User{
		Roles: []string{auth.DefaultAdminRole},
		ProductGrants: map[string][]string{
			"product-1":                        {auth.ActViewProduct.String()},
			"product-2":                        {auth.ActViewProduct.String()},
			"product-without-view-permissions": {},
		},
	}

	userProducts := authorizer.GetUserProductsWithViewAccess(user)
	assert.Nil(t, userProducts)
}

func TestIsadmin(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		user    *entity.User
		isAdmin bool
	}{
		{
			name:    "user with cfg admin role in its roles is admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{auth.DefaultAdminRole}).Build(),
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
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	customAdminRole := "admin"

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy, casbinauth.WithAdminRole(customAdminRole))
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
			user:    testhelpers.NewUserBuilder().WithRoles([]string{"user", auth.DefaultAdminRole}).Build(),
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
