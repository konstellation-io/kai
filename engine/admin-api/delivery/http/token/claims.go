package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type CustomClaims struct {
	ProductGrants entity.ProductGrants `json:"product_roles"`
	RealmAccess   RealmAccess          `json:"realm_access"`
	Name          string               `json:"preferred_username"`
	Email         string               `json:"email"`
	jwt.RegisteredClaims
}

type ProductRoles map[string][]string

type RealmAccess struct {
	Roles []string `json:"roles"`
}
