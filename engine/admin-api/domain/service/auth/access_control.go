package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type Action string

const DefaultAdminRole = "ADMIN"

const ActViewProduct Action = "view_product"
const ActCreateProduct Action = "create_product"

const ActCreateVersion Action = "create_version"
const ActStartVersion Action = "start_version"
const ActStartCriticalVersion Action = "start_critical_version"
const ActStopVersion Action = "stop_version"
const ActPublishVersion Action = "publish_version"
const ActUnpublishVersion Action = "unpublish_version"
const ActEditVersion Action = "edit_version"
const ActViewVersion Action = "view_version"

const ActViewMetrics Action = "view_metrics"
const ActViewServerInfo Action = "view_server_info"

const ActViewUserActivities Action = "view_user_activities"
const ActUpdateUserGrants Action = "update_user_grants"

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
