package entity

import (
	"time"
)

type UserActivityType string

const (
	UserActivityTypeCreateProduct       UserActivityType = "CREATE_PRODUCT"
	UserActivityTypeCreateVersion       UserActivityType = "CREATE_VERSION"
	UserActivityTypePublishVersion      UserActivityType = "PUBLISH_VERSION"
	UserActivityTypeUnpublishVersion    UserActivityType = "UNPUBLISH_VERSION"
	UserActivityTypeStartVersion        UserActivityType = "START_VERSION"
	UserActivityTypeStopVersion         UserActivityType = "STOP_VERSION"
	UserActivityTypeUpdateProductGrants UserActivityType = "UPDATE_PRODUCT_GRANTS"
)

func (e UserActivityType) IsValid() bool {
	switch e {
	case UserActivityTypeCreateProduct,
		UserActivityTypeCreateVersion,
		UserActivityTypePublishVersion,
		UserActivityTypeUnpublishVersion,
		UserActivityTypeStartVersion,
		UserActivityTypeStopVersion,
		UserActivityTypeUpdateProductGrants:
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
