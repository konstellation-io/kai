//go:build integration

package user

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (s *KeycloakSuite) TestCreateGroupWithPolicy() {
	var (
		ctx       = context.Background()
		groupName = "test-group"
		policy    = "test-policy"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, policy)
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	group, err := s.keycloakClient.GetGroupByPath(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		fmt.Sprintf("/%s", groupName),
	)
	s.Require().NoError(err)

	s.Assert().Contains((*group.Attributes)["policy"], policy)
}

func (s *KeycloakSuite) TestCreateGroupWithPolicy_InvalidGroupName() {
	var (
		ctx       = context.Background()
		groupName = ""
		policy    = "test-policy"
	)

	err := s.keycloakUserRegistry.CreateGroupWithPolicy(ctx, groupName, policy)
	s.Require().Error(err)
}
