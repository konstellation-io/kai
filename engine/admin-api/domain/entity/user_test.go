package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"gotest.tools/v3/assert"
)

func TestUserIsAdmin_DefaultAdminRole(t *testing.T) {
	testCases := []struct {
		name    string
		user    *entity.User
		isAdmin bool
	}{
		{
			name:    "User with default admin role in his roles is an admin",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{entity.DefaultAdminRole}).Build(),
			isAdmin: true,
		},
		{
			name:    "User without default admin role on his roles is not an admin",
			user:    testhelpers.NewUserBuilder().Build(),
			isAdmin: false,
		},
		{
			name:    "User with multiple roles is admin if default admin role is one of them",
			user:    testhelpers.NewUserBuilder().WithRoles([]string{"user", entity.DefaultAdminRole}).Build(),
			isAdmin: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := tc.user.IsAdmin()
			assert.Equal(t, tc.isAdmin, isAdmin)
		})
	}
}

func TestUserIsAdmin_OptionalAdminRole(t *testing.T) {
	testCases := []struct {
		name      string
		adminRole string
		user      *entity.User
		isAdmin   bool
	}{
		{
			name:      "User with optional admin role in his roles is an admin",
			adminRole: "admin",
			user:      testhelpers.NewUserBuilder().WithRoles([]string{"admin"}).Build(),
			isAdmin:   true,
		},
		{
			name:      "User without optional admin role in his roles is not an admin",
			adminRole: "admin",
			user:      testhelpers.NewUserBuilder().Build(),
			isAdmin:   false,
		},
		{
			name:      "User with multiple roles is admin if optional admin role is one of them",
			adminRole: "admin",
			user:      testhelpers.NewUserBuilder().WithRoles([]string{"user", "admin"}).Build(),
			isAdmin:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := tc.user.IsAdmin(tc.adminRole)
			assert.Equal(t, tc.isAdmin, isAdmin)
		})
	}
}
