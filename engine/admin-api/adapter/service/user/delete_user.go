package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (ur *KeycloakUserRegistry) DeleteUser(ctx context.Context, name string) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	users, err := ur.client.GetUsers(
		ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{Username: gocloak.StringP(name)})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return ErrUserNotFound
	}

	userID := *users[0].ID

	err = ur.client.DeleteUser(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), userID)
	if err != nil {
		return fmt.Errorf("deleting Keycloak user %q: %w", userID, err)
	}

	return nil
}
