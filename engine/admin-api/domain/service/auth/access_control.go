package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type Action string

const DefaultAdminRole = "ADMIN"

const (
	ActViewProduct   Action = "view_product"
	ActCreateProduct Action = "create_product"

	ActCreateVersion Action = "create_version"
	ActManageVersion Action = "manage_version"

	ActRegisterProcess         Action = "register_process"
	ActDeleteRegisteredProcess Action = "delete_registered_process"

	ActRegisterPublicProcess Action = "register_public_process"
	ActDeletePublicProcess   Action = "delete_public_process"

	ActManageCriticalVersion Action = "manage_critical_version"
	ActUpdateUserGrants      Action = "update_user_grants"
	ActManageProductUser     Action = "manage_product_user"

	ActViewServerInfo     Action = "view_server_info"     // To be deprecated
	ActViewUserActivities Action = "view_user_activities" // To be deprecated
)

func (e Action) IsValid() bool {
	switch e {
	case ActViewProduct, ActCreateProduct, ActCreateVersion, ActManageVersion,
		ActRegisterProcess, ActDeleteRegisteredProcess, ActRegisterPublicProcess,
		ActDeletePublicProcess, ActManageCriticalVersion, ActUpdateUserGrants,
		ActViewUserActivities, ActViewServerInfo, ActManageProductUser:
		return true
	}

	return false
}

func (e Action) String() string {
	return string(e)
}

type AccessControl interface {
	CheckRoleGrants(user *entity.User, action Action) error
	CheckProductGrants(user *entity.User, product string, action Action) error
	IsAdmin(user *entity.User) bool
	GetUserProducts(user *entity.User) []string
}

func GetProductMantainerGrants() []Action {
	return []Action{
		ActViewProduct,
		ActCreateVersion,
		ActManageVersion,
		ActRegisterProcess,
		ActDeleteRegisteredProcess,
		ActManageProductUser,
	}
}

func GetProductUserGrants() []Action {
	return []Action{
		ActViewProduct,
		ActCreateVersion,
		ActManageVersion,
		ActRegisterProcess,
		ActDeleteRegisteredProcess, // TODO: should a regular user be able to do this?
	}
}
