package auth

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/auth_${GOFILE} -package=mocks

type AccessControlAction string

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

const ActView AccessControlAction = "view"
const ActEdit AccessControlAction = "edit"

func (e AccessControlAction) IsValid() bool {
	switch e {
	case ActView, ActEdit, ActCreateProduct, ActStartVersion, ActStopVersion,
		ActPublishVersion, ActUnpublishVersion, ActEditVersion, ActViewMetrics,
		ActViewUserActivities, ActViewProduct, ActCreateVersion, ActViewVersion:
		return true
	}

	return false
}

func (e AccessControlAction) String() string {
	return string(e)
}

//nolint:godox // Remove this nolint statement after the TODO is done.
type AccessControl interface { // TODO: move to middleware.
	CheckPermission(user *entity.User, product string, action AccessControlAction) error
}
