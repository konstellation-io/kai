//go:build integration

package user

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (s *KeycloakSuite) TestCreateUserWithinGroup() {
	var (
		ctx          = context.Background()
		userName     = "test-user"
		testPassword = "test-password"
		groupName    = "test-group"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, "test-policy")
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.CreateUserWithinGroup(ctx, userName, testPassword, groupName)
	s.Require().NoError(err)

	user, err := s.findUserByName(ctx, userName)
	s.Require().NoError(err)

	groups, err := s.keycloakClient.GetUserGroups(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.PString(user.ID),
		gocloak.GetGroupsParams{},
	)
	s.Require().NoError(err)

	s.Assert().Truef(s.groupContainsGroupWithName(groups, groupName), "expected group not found in user's groups")
}

func (s *KeycloakSuite) TestCreateUserWithinGroup_GroupDoesntExist() {
	var (
		ctx          = context.Background()
		userName     = "new-user"
		testPassword = "test-password"
		groupName    = "test-group"
	)

	err := s.keycloakUserRegistry.CreateUserWithinGroup(ctx, userName, testPassword, groupName)
	s.Require().Error(err)

}

func (s *KeycloakSuite) findUserByName(ctx context.Context, name string) (*gocloak.User, error) {
	err := s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)
	users, err := s.keycloakClient.GetUsers(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{},
	)
	s.Require().NoError(err)

	for _, user := range users {
		if *user.Username == name {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (s *KeycloakSuite) groupContainsGroupWithName(
	groups []*gocloak.Group,
	name string,
) bool {
	for _, group := range groups {
		if *group.Name == name {
			return true
		}
	}

	return false
}
