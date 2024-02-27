package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

func (ur *KeycloakUserRegistry) AddProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error {
	user, err := ur.getUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	userGrantsByProduct, err := ur.getUserProductGrants(user)
	if err != nil {
		return err
	}

	userGrantsByProduct[product] = mergeGrants(userGrantsByProduct[product], grants)

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

func mergeGrants(actualGrants, newGrants []auth.Action) []auth.Action {
	grantsSet := map[auth.Action]bool{}

	for _, grant := range append(actualGrants, newGrants...) {
		grantsSet[grant] = true
	}

	mergedGrants := make([]auth.Action, 0, len(actualGrants)+len(newGrants))
	for key := range grantsSet {
		mergedGrants = append(mergedGrants, key)
	}

	return mergedGrants
}
