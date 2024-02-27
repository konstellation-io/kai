package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
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

func (ur *KeycloakUserRegistry) getUserProductGrants(user *gocloak.User) (map[string][]auth.Action, error) {
	if user.Attributes == nil {
		user.Attributes = &map[string][]string{}
	}

	rolesAttribute, ok := (*user.Attributes)["product_roles"]
	if !ok {
		rolesAttribute = make([]string, 1)
	}

	userProductGrants := rolesAttribute[0]
	userGrantsByProduct := make(map[string][]auth.Action)

	if userProductGrants != "" {
		if err := json.Unmarshal([]byte(userProductGrants), &userGrantsByProduct); err != nil {
			return nil, fmt.Errorf("unmarshalling user's product grants: %w", err)
		}
	}

	return userGrantsByProduct, nil
}

func (ur *KeycloakUserRegistry) getUserByEmail(ctx context.Context, userEmail string) (*gocloak.User, error) {
	err := ur.refreshToken(ctx)
	if err != nil {
		return nil, err
	}

	users, err := ur.client.GetUsers(
		ctx, ur.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{Email: gocloak.StringP(userEmail)},
	)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	return users[0], nil
}

func (ur *KeycloakUserRegistry) updatedUserWithNewGrants(
	ctx context.Context,
	user *gocloak.User,
	newGrants map[string][]auth.Action,
) error {
	marshalledRoles, err := json.Marshal(newGrants)
	if err != nil {
		return fmt.Errorf("marshaling user's product grants: %w", err)
	}

	(*user.Attributes)["product_roles"] = []string{string(marshalledRoles)}

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
