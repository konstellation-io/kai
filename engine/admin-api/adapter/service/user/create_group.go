package user

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) CreateGroupWithPolicy(ctx context.Context, name, policy string) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	_, err = ur.client.CreateGroup(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), gocloak.Group{
		Name: gocloak.StringP(name),
		Attributes: &map[string][]string{
			"policy": {policy},
		},
	})
	if err != nil {
		return fmt.Errorf("creating Keycloak group: %w", err)
	}

	return nil
}
