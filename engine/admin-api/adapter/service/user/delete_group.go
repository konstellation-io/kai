package user

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) DeleteGroup(ctx context.Context, name string) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	group, err := ur.client.GetGroupByPath(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), fmt.Sprintf("/%s", name))
	if err != nil {
		return err
	}

	err = ur.client.DeleteGroup(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), *group.ID)
	if err != nil {
		return fmt.Errorf("deleting Keycloak group %q: %w", name, err)
	}

	return nil
}
