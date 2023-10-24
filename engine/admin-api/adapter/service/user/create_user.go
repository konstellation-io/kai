package user

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) CreateUserWithinGroup(ctx context.Context, name, password, group string) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	_, err = ur.client.CreateUser(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), gocloak.User{
		Username:      gocloak.StringP(name),
		EmailVerified: gocloak.BoolP(true),
		Groups:        &[]string{group},
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Value: gocloak.StringP(password),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("creating Keycloak user: %w", err)
	}

	return nil
}
