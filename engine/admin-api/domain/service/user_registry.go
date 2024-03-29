package service

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type UserRegistry interface {
	AddProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error
	RevokeProductGrants(ctx context.Context, userEmail, product string, grants []auth.Action) error
	CreateGroupWithPolicy(ctx context.Context, name, policy string) error
	DeleteGroup(ctx context.Context, name string) error
	CreateUserWithinGroup(ctx context.Context, name, password, group string) error
	DeleteUser(ctx context.Context, name string) error
}
