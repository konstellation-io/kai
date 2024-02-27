package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type Action string

const DefaultAdminRole = "ADMIN"

const (
	ActViewProduct   Action = "view_product"
	ActCreateProduct Action = "create_product"

	ActManageVersion Action = "manage_version"

	ActRegisterProcess         Action = "register_process"
	ActDeleteRegisteredProcess Action = "delete_registered_process"

	ActRegisterPublicProcess Action = "register_public_process"
	ActDeletePublicProcess   Action = "delete_public_process"

	ActManageCriticalVersion    Action = "manage_critical_version"
	ActManageProductUsers       Action = "manage_product_user"
	ActManageProductMaintainers Action = "manage_product_maintainers"

	ActViewServerInfo     Action = "view_server_info"     // To be deprecated
	ActViewUserActivities Action = "view_user_activities" // To be deprecated
)

func (e Action) IsValid() bool {
	switch e {
	case ActViewProduct, ActCreateProduct, ActManageVersion,
		ActRegisterProcess, ActDeleteRegisteredProcess, ActRegisterPublicProcess,
		ActDeletePublicProcess, ActManageCriticalVersion, ActViewUserActivities,
		ActViewServerInfo, ActManageProductUsers:
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
	GetUserProductsWithViewAccess(user *entity.User) []string
}

func GetDefaultUserGrants() []Action {
	return []Action{
		ActViewProduct,
		ActManageVersion,
		ActRegisterProcess,
	}
}

func GetDefaultMaintainerGrants() []Action {
	return append(
		GetDefaultUserGrants(),
		ActDeleteRegisteredProcess,
		ActManageCriticalVersion,
		ActManageProductUsers,
	)
}
