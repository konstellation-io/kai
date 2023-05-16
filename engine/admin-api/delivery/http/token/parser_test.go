package token_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/token"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tokenKey = []byte("")

func newTokenWithProductRoles(userRoles *entity.User) (string, error) {
	if userRoles == nil {
		return jwt.NewWithClaims(jwt.SigningMethodHS256, nil).SignedString(tokenKey)
	}

	claims := token.CustomClaims{
		ProductGrants: userRoles.ProductGrants,
		RealmAccess: token.RealmAccess{
			Roles: userRoles.Roles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userRoles.ID,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tokenKey)
}

func Test_CustomClaims(t *testing.T) {
	expectedUser := &entity.User{
		ID: "test-user",
		ProductGrants: entity.ProductGrants{
			"test": {
				"view_product",
			},
		},
		Roles: []string{"USER"},
	}

	accessToken, err := newTokenWithProductRoles(expectedUser)
	require.NoError(t, err)

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
	accessToken, err := newTokenWithProductRoles(nil)
	require.NoError(t, err)

	tokenParser := token.NewParser()

	_, err = tokenParser.GetUser(accessToken)
	assert.NoError(t, err)
}
