package user

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) RevokeUserProductGrants(ctx context.Context, userEmail string, product string, grants []auth.Action) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	users, err := ur.client.GetUsers(ctx, ur.token.AccessToken, viper.GetString(config.KeycloakRealmKey), gocloak.GetUsersParams{Email: gocloak.StringP(userEmail)})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return ErrUserNotFound
	}

	user := users[0]

	if user.Attributes == nil {
		user.Attributes = &map[string][]string{}
	}

	rolesAttribute, ok := (*user.Attributes)["product_roles"]
	if !ok {
		rolesAttribute = make([]string, 1)
	}

	userProductGrants := rolesAttribute[0]
	userGrantsByProduct := make(map[string][]auth.Action)

	if userProductGrants == "" {
		userProductGrants = "{}"
	}

	if err = json.Unmarshal([]byte(userProductGrants), &userGrantsByProduct); err != nil {
		return fmt.Errorf("unmarshalling user's product grants: %w", err)
	}

	userGrantsByProduct[product] = slices.DeleteFunc(userGrantsByProduct[product], func(e auth.Action) bool {
		return slices.Contains(grants, e)
	})

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
