//go:build integration

package user

import (
	"context"
	"encoding/json"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (s *KeycloakSuite) TestAddUserProductGrants_NoPreviousExisting() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	// WHEN updating grants for a product for the first time
	err := s.keycloakUserRegistry.AddProductGrants(
		ctx,
		*user.Email,
		product,
		[]auth.Action{auth.ActViewProduct, auth.ActManageVersion},
	)
	s.Require().NoError(err)

	// THEN grants for the product are added
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product: []interface{}{auth.ActViewProduct.String(), auth.ActManageVersion.String()},
	}

	s.ElementsMatch(expectedResult[product], obtainedResult[product])
}

func (s *KeycloakSuite) TestAddUserProductGrants_MergeNewGrants_NoDups() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	err := s.keycloakUserRegistry.AddProductGrants(
		ctx,
		*user.Email,
		product,
		[]auth.Action{auth.ActViewProduct, auth.ActManageVersion},
	)
	s.Require().NoError(err)

	// WHEN
	err = s.keycloakUserRegistry.AddProductGrants(
		ctx,
		*user.Email,
		product,
		[]auth.Action{auth.ActManageVersion},
	)
	s.Require().NoError(err)

	// THEN grants for the product are added
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product: []interface{}{
			auth.ActViewProduct.String(),
			auth.ActManageVersion.String(),
		},
	}

	s.ElementsMatch(expectedResult[product], obtainedResult[product])
}

func (s *KeycloakSuite) TestAddUserProductGrants_MergeNewGrants_WithDups() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	err := s.keycloakUserRegistry.AddProductGrants(
		ctx,
		*user.Email,
		product,
		[]auth.Action{auth.ActViewProduct, auth.ActManageVersion},
	)
	s.Require().NoError(err)

	// WHEN
	err = s.keycloakUserRegistry.AddProductGrants(
		ctx,
		*user.Email,
		product,
		[]auth.Action{auth.ActViewProduct, auth.ActManageVersion, auth.ActManageVersion},
	)
	s.Require().NoError(err)

	// THEN grants for the product are added
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product: []interface{}{
			auth.ActViewProduct.String(),
			auth.ActManageVersion.String(),
		},
	}

	s.ElementsMatch(expectedResult[product], obtainedResult[product])
}
