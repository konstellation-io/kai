package user

import (
	"context"
	"fmt"
	"slices"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (ur *KeycloakUserRegistry) RevokeProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error {
	user, err := ur.getUserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("getting user by email: %w", err)
	}

	userGrantsByProduct, err := ur.getUserProductGrants(user)
	if err != nil {
		return fmt.Errorf("getting user's product grants: %w", err)
	}

	userGrantsByProduct[product] = slices.DeleteFunc(userGrantsByProduct[product], func(e auth.Action) bool {
		return slices.Contains(grants, e)
	})

	return ur.updatedUserWithNewGrants(ctx, user, userGrantsByProduct)
}
