package entity

import (
	"time"
)

type UserActivityType string

const (
	UserActivityTypeCreateRuntime            UserActivityType = "CREATE_RUNTIME"
	UserActivityTypeCreateVersion            UserActivityType = "CREATE_VERSION"
	UserActivityTypePublishVersion           UserActivityType = "PUBLISH_VERSION"
	UserActivityTypeUnpublishVersion         UserActivityType = "UNPUBLISH_VERSION"
	UserActivityTypeStartVersion             UserActivityType = "START_VERSION"
	UserActivityTypeStopVersion              UserActivityType = "STOP_VERSION"
	UserActivityTypeUpdateProductPermissions UserActivityType = "UPDATE_ACCESS_LEVELS"
)

func (e UserActivityType) IsValid() bool {
	switch e {
	case UserActivityTypeCreateRuntime,
		UserActivityTypeCreateVersion,
		UserActivityTypePublishVersion,
		UserActivityTypeUnpublishVersion,
		UserActivityTypeStartVersion,
		UserActivityTypeStopVersion,
		UserActivityTypeUpdateProductPermissions:
		return true
	}
	return false
}

func (e UserActivityType) String() string {
	return string(e)
}

type UserActivityVar struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type UserActivity struct {
	ID     string             `bson:"_id"`
	Date   time.Time          `bson:"date"`
	UserID string             `bson:"userId"`
	Type   UserActivityType   `bson:"type"`
	Vars   []*UserActivityVar `bson:"vars"`
}
