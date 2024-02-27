package user

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (ur *KeycloakUserRegistry) AddProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error {
	user, err := ur.getUserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("getting user by email: %w", err)
	}

	userGrantsByProduct, err := ur.getUserProductGrants(user)
	if err != nil {
		return fmt.Errorf("getting user's product grants: %w", err)
	}

	userGrantsByProduct[product] = mergeGrants(userGrantsByProduct[product], grants)

	return ur.updatedUserWithNewGrants(ctx, user, userGrantsByProduct)
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
