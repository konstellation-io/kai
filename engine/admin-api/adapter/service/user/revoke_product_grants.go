package user

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) RevokeProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error {
	err := ur.refreshToken(ctx)
	if err != nil {
		return err
	}

	user, err := ur.getUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	userGrantsByProduct, err := ur.getUserProductGrants(user)
	if err != nil {
		return err
	}

	userGrantsByProduct[product] = slices.DeleteFunc(userGrantsByProduct[product], func(e auth.Action) bool {
		return slices.Contains(grants, e)
	})

	marshalledRoles, err := json.Marshal(userGrantsByProduct)
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
