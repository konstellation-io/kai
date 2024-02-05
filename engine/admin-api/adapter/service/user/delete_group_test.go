//go:build integration

package user

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (s *KeycloakSuite) TestDeleteGroup() {
	var (
		ctx       = context.Background()
		groupName = "test-group"
		policy    = "test-policy"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, policy)
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.DeleteGroup(ctx, groupName)
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	groups, err := s.keycloakClient.GetGroups(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetGroupsParams{},
	)
	s.Require().NoError(err)
	s.Assert().Empty(groups)
}
