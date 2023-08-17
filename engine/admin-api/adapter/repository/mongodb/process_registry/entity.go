package process_registry

import "time"

type processRegistryDTO struct {
	ID         string    `bson:"_id"`
	Name       string    `bson:"name"`
	Version    string    `bson:"version"`
	Type       string    `bson:"type"`
	UploadDate time.Time `bson:"uploadDate"`
	Owner      string    `bson:"owner"`
}
