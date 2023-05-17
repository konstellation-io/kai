package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type GocloakService interface {
	GetUserByID(userID string) (*entity.User, error)
	UpdateUserProductPermissions(userID string, product string, permissions []string) error
}
