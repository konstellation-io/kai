package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

type GocloakService interface {
	CreateUser(userData entity.UserGocloakData) error
	GetUserByID(userID string) (entity.UserGocloakData, error)
	UpdateUserProductPermissions(userID string, product string, roles []string) error
}
