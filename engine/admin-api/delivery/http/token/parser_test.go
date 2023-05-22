//go:build unit

package token_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/token"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CustomClaims(t *testing.T) {
	expectedUser := testhelpers.NewUserBuilder().
		WithRoles([]string{"USER"}).
		WithProductGrants(entity.ProductGrants{
			"test": {
				"view_product",
			},
		}).
		Build()

	accessToken := testhelpers.NewTokenBuilder().
		WithUser(expectedUser).
		Build()

	tokenParser := token.NewParser()

	userRoles, err := tokenParser.GetUser(accessToken)
	require.NoError(t, err)
	assert.Equal(t, expectedUser, userRoles)
}

func Test_CustomClaims_FailsIfTokenIsNotValid(t *testing.T) {
	accessToken := "this-is-an-invalid-token"

	tokenParser := token.NewParser()

	productRoles, err := tokenParser.GetUser(accessToken)
	assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	assert.Nil(t, productRoles)
}

func Test_CustomClaims_ReturnNilIfNoAccessClaims(t *testing.T) {
	accessToken := testhelpers.NewTokenBuilder().Build()

	tokenParser := token.NewParser()

	_, err := tokenParser.GetUser(accessToken)
	assert.NoError(t, err)
}
