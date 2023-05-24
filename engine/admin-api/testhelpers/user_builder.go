package testhelpers

import (
	"github.com/bxcodec/faker/v3"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type UserBuilder struct {
	user *entity.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &entity.User{
			ID:            faker.UUIDHyphenated(),
			Roles:         []string{"USER"},
			ProductGrants: nil,
		},
	}
}

func (u *UserBuilder) WithID(id string) *UserBuilder {
	u.user.ID = id
	return u
}

func (u *UserBuilder) WithRoles(roles []string) *UserBuilder {
	u.user.Roles = roles
	return u
}

func (u *UserBuilder) WithProductGrants(grants entity.ProductGrants) *UserBuilder {
	u.user.ProductGrants = grants
	return u
}

func (u *UserBuilder) Build() *entity.User {
	return u.user
}
