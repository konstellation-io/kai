package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

type KeycloakUserRegistry struct {
	client                *gocloak.GoCloak
	token                 *gocloak.JWT
	tokenExpiresAt        time.Time
	refreshTokenExpiresAt time.Time
}

func WithClient(endpoint string) *gocloak.GoCloak {
	client := gocloak.NewClient(endpoint)

	return client
}

func NewKeycloakUserRegistry(client *gocloak.GoCloak) (*KeycloakUserRegistry, error) {
	ctx := context.Background()
	now := time.Now()

	token, err := client.LoginAdmin(
		ctx,
		viper.GetString(config.KeycloakAdminUserKey),
		viper.GetString(config.KeycloakAdminPasswordKey),
		viper.GetString(config.KeycloakMasterRealmKey),
	)
	if err != nil {
		return nil, fmt.Errorf("login with admin user: %w", err)
	}

	return &KeycloakUserRegistry{
		client:                client,
		token:                 token,
		tokenExpiresAt:        now.Add(time.Duration(token.ExpiresIn) * time.Second),
		refreshTokenExpiresAt: now.Add(time.Duration(token.RefreshExpiresIn) * time.Second),
	}, nil
}

func (ur *KeycloakUserRegistry) refreshToken(ctx context.Context) error {
	now := time.Now()

	if now.Before(ur.tokenExpiresAt) {
		return nil
	}

	var (
		token *gocloak.JWT
		err   error
	)

	if now.Before(ur.refreshTokenExpiresAt) {
		token, err = ur.client.RefreshToken(
			ctx,
			ur.token.RefreshToken,
			viper.GetString(config.KeycloakAdminClientIDKey),
			"",
			viper.GetString(config.KeycloakMasterRealmKey),
		)

		if err != nil {
			return fmt.Errorf("refreshing token: %w", err)
		}
	} else {
		token, err = ur.client.LoginAdmin(
			ctx,
			viper.GetString(config.KeycloakAdminUserKey),
			viper.GetString(config.KeycloakAdminPasswordKey),
			viper.GetString(config.KeycloakMasterRealmKey),
		)

		if err != nil {
			return fmt.Errorf("login with admin user: %w", err)
		}
	}

	ur.token = token
	ur.tokenExpiresAt = now.Add(time.Duration(token.ExpiresIn) * time.Second)
	ur.refreshTokenExpiresAt = now.Add(time.Duration(token.RefreshExpiresIn) * time.Second)

	return nil
}
