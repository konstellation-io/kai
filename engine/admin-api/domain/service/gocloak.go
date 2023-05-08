package service

import (
	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

type GocloakService interface {
	CreateUser(userData entity.UserGocloakData) error
	GetUserByID(userID string) (entity.UserGocloakData, error)
	UpdateUserRoles(userID string, product string, roles []string) error
}
