package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type AccessControlAction string

const DefaultAdminRole = "ADMIN"

const ActViewProduct AccessControlAction = "view_product"
const ActCreateProduct AccessControlAction = "create_product"

const ActCreateVersion AccessControlAction = "create_version"
const ActStartVersion AccessControlAction = "start_version"
const ActStopVersion AccessControlAction = "stop_version"
const ActPublishVersion AccessControlAction = "publish_version"
const ActUnpublishVersion AccessControlAction = "unpublish_version"
const ActEditVersion AccessControlAction = "edit_version"
const ActViewVersion AccessControlAction = "view_version"

const ActViewMetrics AccessControlAction = "view_metrics"

const ActViewUserActivities AccessControlAction = "view_user_activities"
const ActUpdateUserGrants AccessControlAction = "update_user_grants"

func (e AccessControlAction) IsValid() bool {
	switch e {
	case ActCreateProduct, ActStartVersion, ActStopVersion, ActUpdateUserGrants,
		ActPublishVersion, ActUnpublishVersion, ActEditVersion, ActViewMetrics,
		ActViewUserActivities, ActViewProduct, ActCreateVersion, ActViewVersion:
		return true
	}

	return false
}

func (e AccessControlAction) String() string {
	return string(e)
}

type AccessControl interface {
	CheckProductGrants(user *entity.User, product string, action AccessControlAction) error
	CheckGrants(user *entity.User, action AccessControlAction) error
	IsAdmin(user *entity.User) bool
}
