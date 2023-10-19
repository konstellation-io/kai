package service

import "context"

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type UserRegistry interface {
	UpdateUserProductGrants(ctx context.Context, userID string, product string, grants []string) error
}
