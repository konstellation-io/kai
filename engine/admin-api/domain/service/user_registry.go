package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type UserRegistry interface {
	GetUserByID(userID string) (*entity.User, error)
	UpdateUserProductGrants(userID string, product string, grants []string) error
}
