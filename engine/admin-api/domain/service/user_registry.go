package service

import "context"

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type UserRegistry interface {
	UpdateUserProductGrants(ctx context.Context, userID string, product string, grants []string) error
	CreateGroupWithPolicy(ctx context.Context, name, policy string) error
	DeleteGroup(ctx context.Context, name string) error
	CreateUserWithinGroup(ctx context.Context, name, password, group string) error
	DeleteUser(ctx context.Context, name string) error
}
