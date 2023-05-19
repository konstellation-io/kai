package testhelpers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/token"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const defaultTokenKey = ""

type TokenBuilder struct {
	user *entity.User
}

func NewTokenBuilder() *TokenBuilder {
	return &TokenBuilder{
		user: &entity.User{},
	}
}

func (tb *TokenBuilder) WithUser(user *entity.User) *TokenBuilder {
	tb.user = user
	return tb
}

func (tb *TokenBuilder) Build() string {
	claims := token.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: tb.user.ID,
		},
		ProductGrants: tb.user.ProductGrants,
		RealmAccess: token.RealmAccess{
			Roles: tb.user.Roles,
		},
	}
	t, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(defaultTokenKey))

	return t
}
