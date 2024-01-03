package processrepository

import "time"

type registeredProcessDTO struct {
	ID         string    `bson:"_id"`
	Name       string    `bson:"name"`
	Version    string    `bson:"version"`
	Type       string    `bson:"type"`
	Image      string    `bson:"image"`
	UploadDate time.Time `bson:"uploadDate"`
	Owner      string    `bson:"owner"`
	Status     string    `bson:"status"`
	Logs       string    `bson:"logs"`
	IsPublic   bool      `bson:"isPublic"`
}
