package user

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) GroupExists(ctx context.Context, name string) (bool, error) {
	err := ur.refreshToken(ctx)
	if err != nil {
		return false, err
	}

	groups, err := ur.client.GetGroups(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), gocloak.GetGroupsParams{})
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if group.Name != nil && *group.Name == name {
			return true, nil
		}
	}

	return false, nil
}
