//go:build integration

package user

import (
	"context"
)

func (s *KeycloakSuite) TestGroupExists_OK() {
	var (
		ctx       = context.Background()
		groupName = "test-group"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, "test-policy")
	s.Require().NoError(err)

	exists, err := s.keycloakUserRegistry.GroupExists(ctx, groupName)
	s.Require().NoError(err)

	s.True(exists)
}

func (s *KeycloakSuite) TestGroupExists_GroupNotFound() {
	var (
		ctx       = context.Background()
		groupName = "test-group"
	)

	exists, err := s.keycloakUserRegistry.GroupExists(ctx, groupName)
	s.Require().NoError(err)

	s.False(exists)
}
