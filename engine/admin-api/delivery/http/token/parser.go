package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type Config struct {
	SigninKey interface{}
}

type Parser struct {
	parser *jwt.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: jwt.NewParser(),
	}
}

func (p *Parser) GetUserRoles(accessToken string) (*entity.User, error) {
	claims := &CustomClaims{}

	_, _, err := p.parser.ParseUnverified(accessToken, claims)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	return &entity.User{
		ID:            claims.Subject,
		ProductGrants: claims.ProductGrants,
		Roles:         claims.RealmAccess.Roles,
	}, nil
}
