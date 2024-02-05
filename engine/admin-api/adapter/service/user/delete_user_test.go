//go:build integration

package user

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (s *KeycloakSuite) TestDeleteUser() {
	var (
		ctx          = context.Background()
		userName     = "user-to-delete"
		testPassword = "test-password"
		groupName    = "test-group"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, "test-policy")
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.CreateUserWithinGroup(ctx, userName, testPassword, groupName)
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.DeleteUser(ctx, userName)
	s.Require().NoError(err)

	users, err := s.keycloakClient.GetUsers(ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{Username: gocloak.StringP(userName)},
	)
	s.Require().NoError(err)
	s.Empty(users)
}
