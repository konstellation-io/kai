package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type Action string

const DefaultAdminRole = "ADMIN"

const (
	ActViewProduct   Action = "view_product"
	ActCreateProduct Action = "create_product"

	ActCreateVersion        Action = "create_version"
	ActStartVersion         Action = "start_version"
	ActStartCriticalVersion Action = "start_critical_version"
	ActStopVersion          Action = "stop_version"
	ActPublishVersion       Action = "publish_version"
	ActUnpublishVersion     Action = "unpublish_version"
	ActEditVersion          Action = "edit_version"
	ActViewVersion          Action = "view_version"

	ActViewMetrics    Action = "view_metrics"
	ActViewServerInfo Action = "view_server_info"

	ActViewUserActivities Action = "view_user_activities"
	ActUpdateUserGrants   Action = "update_user_grants"

	ActRegisterProcess       Action = "register_process"
	ActRegisterPublicProcess Action = "register_public_process"
)

func (e Action) IsValid() bool {
	switch e {
	case ActCreateProduct, ActStartVersion, ActStopVersion, ActUpdateUserGrants,
		ActPublishVersion, ActUnpublishVersion, ActEditVersion, ActViewMetrics,
		ActViewUserActivities, ActViewProduct, ActCreateVersion, ActViewVersion,
		ActViewServerInfo, ActStartCriticalVersion:
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
