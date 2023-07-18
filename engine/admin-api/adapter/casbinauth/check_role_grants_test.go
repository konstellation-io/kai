//go:build unit

package casbinauth_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAdminGrants(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

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
			act:           auth.ActUpdateUserGrants,
			expectedError: nil,
		},
		{
			name: "user without ADMIN role can't update user roles",
			user: testhelpers.NewUserBuilder().WithRoles([]string{}).Build(),
			act:  auth.ActUpdateUserGrants,
			expectedError: auth.UnauthorizedError{
				Action: auth.ActUpdateUserGrants,
			},
		},
		{
			name:          "user with ADMIN role can view server info",
			user:          testhelpers.NewUserBuilder().WithRoles([]string{"ADMIN"}).Build(),
			act:           auth.ActViewServerInfo,
			expectedError: nil,
		},
		{
			name:          "user with MLE role can view server info",
			user:          testhelpers.NewUserBuilder().WithRoles([]string{"MLE"}).Build(),
			act:           auth.ActViewServerInfo,
			expectedError: nil,
		},
		{
			name: "an user who is neither ADMIN or MLE can't view server info",
			user: testhelpers.NewUserBuilder().WithRoles([]string{"USER"}).Build(),
			act:  auth.ActViewServerInfo,
			expectedError: auth.UnauthorizedError{
				Action: auth.ActViewServerInfo,
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
