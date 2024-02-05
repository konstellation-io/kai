package user

import (
	"context"
	"encoding/json"
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

func (ur *KeycloakUserRegistry) UpdateUserProductGrants(ctx context.Context, userID, product string, grants []string) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	user, err := ur.client.GetUserByID(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), userID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if user.Attributes == nil {
		user.Attributes = &map[string][]string{}
	}

	rolesAttribute, ok := (*user.Attributes)["product_roles"]
	if !ok {
		rolesAttribute = make([]string, 1)
	}

	userProductGrants := rolesAttribute[0]
	userGrantsByProduct := make(map[string]interface{})

	if userProductGrants == "" {
		userProductGrants = "{}"
	}

	if err = json.Unmarshal([]byte(userProductGrants), &userGrantsByProduct); err != nil {
		return fmt.Errorf("unmarshalling user's product grants: %w", err)
	}

	if len(grants) == 0 {
		delete(userGrantsByProduct, product)
	} else {
		userGrantsByProduct[product] = grants
	}

	marshalledRoles, err := json.Marshal(userGrantsByProduct)
	if err != nil {
		return fmt.Errorf("marshaling user's product grants: %w", err)
	}

	rolesAttribute[0] = string(marshalledRoles)

	(*user.Attributes)["product_roles"] = rolesAttribute

	err = ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	err = ur.client.UpdateUser(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), *user)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	return nil
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
