//go:build unit

package casbinauth_test

import (
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAdminGrants(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		user          *entity.User
		act           auth.Action
		expectedError error
	}{
		{
			name:          "user with ADMIN role can update user roles",
			user:          testhelpers.NewUserBuilder().WithRoles([]string{auth.DefaultAdminRole}).Build(),
			act:           auth.ActManageVersion,
			expectedError: nil,
		},
		{
			name: "user without ADMIN role can't update user roles",
			user: testhelpers.NewUserBuilder().WithRoles([]string{}).Build(),
			act:  auth.ActManageProductUsers,
			expectedError: auth.UnauthorizedError{
				Action: auth.ActManageProductUsers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := authorizer.CheckRoleGrants(tc.user, tc.act)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
