package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type UserRegistry interface {
	UpdateUserProductGrants(userID string, product string, grants []string) error
}
